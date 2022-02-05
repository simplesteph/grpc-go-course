package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"../calculatorpb"

	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Calculator Client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	// fmt.Printf("Created client: %f", c)

	//doUnary(c)

	doStreamingService(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.SumRequest{
		FirstNumber:  1,
		SecondNumber: 20,
	}

	res, err := c.Sum(context.Background(), req)

	if err != nil {
		log.Fatalf("Calculator Request error:%v", err)
	}

	log.Printf("Calculator Sum Result: %v", res.SumResult)
}

func doStreamingService(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.PNDRequest{
		Number: 11247480,
	}

	stream, err := c.PND(context.Background(), req)

	if err != nil {
		log.Fatalf("PND Request error:%v", err)
	}

	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Streaming error:%v", err)
		}
		log.Printf("Received:%v", msg.GetPrime())
	}
}
