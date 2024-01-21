package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/jadeidev/grpc-go-course/calc/calcpb"

	"google.golang.org/grpc"
)

type server struct {
	calcpb.UnimplementedCalcServiceServer
}

func (*server) Calc(ctx context.Context, req *calcpb.CalcRequest) (*calcpb.CalcResponse, error) {
	fmt.Printf("Calc function was invoked with %v\n", req)
	firstNumber := req.GetFirstNumber()
	secondNumber := req.GetSecondNumber()
	result := firstNumber + secondNumber
	res := &calcpb.CalcResponse{
		Result: result,
	}
	return res, nil
}

func (*server) CalcPrimeNumber(req *calcpb.PrimeNumberRequest, stream calcpb.CalcService_CalcPrimeNumberServer) error {
	fmt.Printf("CalcPrimeNumber function was invoked with %v\n", req)
	number := req.GetNumber()
	k := int32(2)
	for number > 1 {
		if number%k == 0 {
			stream.Send(&calcpb.PrimeNumberResponse{
				Result: k,
			})
			number = number / k
		} else {
			k = k + 1
		}
	}
	return nil

}

func (*server) CalcAverage(stream calcpb.CalcService_CalcAverageServer) error {
	fmt.Printf("CalcAverage function was invoked with a streaming request\n")
	var sum = int32(0)
	var count = int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we've finished reading the client stream
			average := float64(sum) / float64(count)
			return stream.SendAndClose(&calcpb.AverageResponse{
				Result: average,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		sum += req.GetNumber()
		count++
	}
}

func main() {
	fmt.Println("Calculator Server")

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	options := []grpc.ServerOption{}
	s := grpc.NewServer(options...)
	calcpb.RegisterCalcServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
