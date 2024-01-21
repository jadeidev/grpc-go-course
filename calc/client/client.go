package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/jadeidev/grpc-go-course/calc/calcpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client")
	options := grpc.WithTransportCredentials(insecure.NewCredentials())
	cc, err := grpc.Dial("localhost:50051", options)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()
	c := calcpb.NewCalcServiceClient(cc)
	// DoUnary
	// req := calcpb.CalcRequest{FirstNumber: 10, SecondNumber: 3}
	// res, err := c.Calc(context.Background(), &req)
	// if err != nil {
	// 	log.Fatalf("error while calling Calc RPC: %v", err)
	// }
	// log.Printf("Response from Calc: %v", res.Result)

	// DoServerStreaming
	// req := calcpb.PrimeNumberRequest{Number: 120}
	// resStream, err := c.CalcPrimeNumber(context.Background(), &req)
	// if err != nil {
	// 	log.Fatalf("error while calling CalcPrimeNumber RPC: %v", err)
	// }
	// for {
	// 	msg, err := resStream.Recv()
	// 	if err == io.EOF {
	// 		// we've reached the end of the stream
	// 		break
	// 	}
	// 	if err != nil {
	// 		log.Fatalf("error while reading stream: %v", err)
	// 	}
	// 	log.Printf("Response from CalcPrimeNumber: %v", msg.GetResult())
	// }

	// DoClientStreaming
	stream, err := c.CalcAverage(context.Background())
	if err != nil {
		log.Fatalf("error while calling CalcAverage: %v", err)
	}
	numbers := []int32{3, 5, 9, 54, 23}
	for _, number := range numbers {
		fmt.Printf("Sending number: %v\n", number)
		stream.Send(&calcpb.AverageRequest{
			Number: number,
		})
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from CalcAverage: %v", err)
	}
	fmt.Printf("CalcAverage Response: %v\n", res)

}
