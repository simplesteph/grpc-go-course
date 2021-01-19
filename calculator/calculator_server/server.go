package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/leonardo-tozato/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	firstNumber := req.GetSum().GetFirstNumber()
	secondNumber := req.GetSum().GetSecondNumber()

	result := firstNumber + secondNumber

	res := &calculatorpb.SumResponse{
		Result: result,
	}

	return res, nil
}

func (*server) PrimeDecomposition(req *calculatorpb.PrimeDecompositionRequest, stream calculatorpb.CalculatorService_PrimeDecompositionServer) error {
	numToDecompose := req.GetNumToDecompose()
	k := int32(2)
	for numToDecompose > 1 {
		if numToDecompose%k == 0 {
			res := &calculatorpb.PrimeDecompositionResponse{
				PrimeFactor: k,
			}
			stream.Send(res)

			numToDecompose /= k
		} else {
			k++
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	fmt.Print("server listening on port 50051")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
