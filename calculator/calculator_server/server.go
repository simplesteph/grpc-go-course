package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"../calculatorpb"
	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {

	result := req.GetFirstNumber() + req.GetSecondNumber()

	res := &calculatorpb.SumResponse{
		SumResult: result,
	}
	return res, nil
}

func (*server) PND(req *calculatorpb.PNDRequest, stream calculatorpb.CalculatorService_PNDServer) error {

	k := req.GetNumber()

	n := int32(2)

	for {
		if k <= 1 {
			break
		}
		for {
			if k%n == 0 {
				k = k / n
				res := &calculatorpb.PNDResponse{
					Prime: n,
				}
				stream.Send(res)
				time.Sleep(time.Second)
			} else {
				n++
				break
			}
		}
	}
	return nil
}

func main() {
	fmt.Println("Calculator Server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	tls := false
	if tls {
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"
		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
		if sslErr != nil {
			log.Fatalf("Failed loading certificates: %v", sslErr)
			return
		}
		opts = append(opts, grpc.Creds(creds))
	}

	s := grpc.NewServer(opts...)
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
