package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"math"
	"net"

	"github.com/leonardo-tozato/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Received a negative number: %v", number))
	}

	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func (*server) StreamMaxNumber(stream calculatorpb.CalculatorService_StreamMaxNumberServer) error {
	var numbers_received []int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error while reading client stream: %v", err)
			return err
		}
		numbers_received = append(numbers_received, req.GetNum())
		result := calcMax(numbers_received)
		if err := stream.Send(&calculatorpb.StreamMaxNumberResponse{
			MaxNum: result,
		}); err != nil {
			log.Fatalf("error while sendind response to the client: %v", err)
			return err
		}
	}
}

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

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	numbersReceived := []int32{}

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			average := calcAverage(numbersReceived)
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: average,
			})
		}

		if err != nil {
			log.Fatalf("error while reading client stream: %v", err)
		}

		numbersReceived = append(numbersReceived, req.GetNum())
	}
}

func calcAverage(nums []int32) float64 {
	var sum int32
	for _, v := range nums {
		sum += v
	}
	return float64(sum) / float64(len(nums))
}

func calcMax(nums []int32) int32 {
	max := nums[0]
	for _, val := range nums {
		if max < val {
			max = val
		}
	}

	return max
}
func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	reflection.Register(s)

	fmt.Print("server listening on port 50051")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
