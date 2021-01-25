package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"strconv"

	"github.com/leonardo-tozato/grpc-go-course/calculator/calculatorpb"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect %v", err)
	}
	defer conn.Close()

	c := calculatorpb.NewCalculatorServiceClient(conn)
	doErrorUnary(c)
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	// correct call
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: int32(-10)})

	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			log.Println(respErr.Message())
			log.Println(respErr.Code())
			fmt.Println("Negative number sent to the server")
			return
		} else {
			log.Fatalf("big error calling sqrt: %v", err)
		}
	}

	log.Printf("square root of number: %v", res.GetNumberRoot())
}

func doMaxNumberStreaming(c calculatorpb.CalculatorServiceClient) {
	stream, err := c.StreamMaxNumber(context.Background())
	if err != nil {
		log.Fatalf("error while creating stream: %v", err)
		return
	}

	waitc := make(chan struct{})

	go func() {
		fmt.Print("Enter numbers: ")
		for {
			var num int32
			_, err := fmt.Scan(&num)
			if err != nil {
				close(waitc)
			}
			if err := stream.Send(&calculatorpb.StreamMaxNumberRequest{
				Num: num,
			}); err != nil {
				log.Fatalf("failed to send number to server: %v", err)
			}
		}
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
			log.Printf("received message: %v", res.GetMaxNum())
		}
		close(waitc)
	}()

	<-waitc

}

func doAverageStreaming(c calculatorpb.CalculatorServiceClient, nums []string) {
	fmt.Println("starting to do a Client Streaming RPC...")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error while calling ComputeAverage RPC: %v", err)
	}

	for _, numStr := range nums {
		num, _ := strconv.Atoi(numStr)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Num: int32(num),
		})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from ComputeAverage: %v", err)
	}

	log.Printf("Average: %v", res.GetAverage())
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient, firstNumber int32) {
	fmt.Println("starting to do a Server Streaming RPC...")

	req := &calculatorpb.PrimeDecompositionRequest{
		NumToDecompose: firstNumber,
	}

	resStream, err := c.PrimeDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling PrimeDecomposition RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}

		factor := msg.GetPrimeFactor()
		log.Printf("prime factor: %v", factor)
	}
}

func doUnary(c calculatorpb.CalculatorServiceClient, firstNumber int32, secondNumber int32) {
	req := &calculatorpb.SumRequest{
		Sum: &calculatorpb.Sum{
			FirstNumber:  firstNumber,
			SecondNumber: secondNumber,
		},
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}

	log.Printf("response from Greet: %v", res.Result)
}
