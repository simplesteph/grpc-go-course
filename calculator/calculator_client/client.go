package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/leonardo-tozato/grpc-go-course/calculator/calculatorpb"

	"google.golang.org/grpc"
)

func main() {
	args := os.Args[1:]
	firstNumber, _ := strconv.Atoi(args[0])
	// secondNumber, _ := strconv.Atoi(args[1])

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect %v", err)
	}
	defer conn.Close()

	c := calculatorpb.NewCalculatorServiceClient(conn)
	doServerStreaming(c, int32(firstNumber))
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
