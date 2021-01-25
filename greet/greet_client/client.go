package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"time"

	"github.com/leonardo-tozato/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	certFile := "ssl/ca.crt"
	creds, err := credentials.NewClientTLSFromFile(certFile, "")

	if err != nil {
		log.Fatalf("error while loading CA trust certificates: %v", err)
	}

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("could not connect %v", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)
	doUnary(c)
	// doServerStreaming(c)
	//doClientStreaming(c)
	//doBidiStreaming(c)
	//doUnaryWithDeadLine(c)
}

func doUnaryWithDeadLine(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Bruce",
			LastName:  "Batma",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {

		statusErr, ok := status.FromError(err)
		if ok && statusErr.Code() == codes.DeadlineExceeded {
			fmt.Println("Timeout!")
			return
		} else {
			log.Fatalf("error while calling Greet RPC: %v", err)
		}
	}

	log.Printf("response from Greet: %v", res.Result)
}

func doBidiStreaming(c greetpb.GreetServiceClient) {
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while creating stream: %v", err)
		return
	}

	requests := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Batman",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Superman",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Flash",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Wonder Woman",
			},
		},
	}

	waitc := make(chan struct{})
	go func () {
		for _, req := range requests {
			log.Printf("sending message: %v", req.GetGreeting().GetFirstName())
			_ = stream.Send(req)
		}
		_ = stream.CloseSend()
	}()

	go func () {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while receiving response: %v", err)
			}
			log.Printf("received message: %v", res.GetResult())
		}
		close(waitc)
	}()

	<-waitc
}
func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Batman",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Superman",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Flash",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Wonder Woman",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet RPC: %v", err)
	}

	for _, req := range requests {
		stream.Send(req)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from LongGreet: %v", err)
	}

	log.Printf("LongGreet response: %v", res)

}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Zika",
			LastName:  "Memo",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}

		log.Printf("result from stream: %v", msg.GetResult())
	}

}

func doUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Bruce",
			LastName:  "Batma",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}

	log.Printf("response from Greet: %v", res.Result)
}
