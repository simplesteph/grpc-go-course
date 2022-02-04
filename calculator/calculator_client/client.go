package main

import (
	"context"
	"fmt"
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

	doUnary(c)

}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.SumRequest{
		FirstNumber:  1,
		SecondNumber: 20,
	}

	res, err := c.Sum(context.Background(), req)

	if err != nil {
		log.Fatalf("Calculator Request error:", err)
	}

	log.Printf("Calculator Sum Result: %v", res.SumResult)
}
