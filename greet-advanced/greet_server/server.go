package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	grpchealth "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/bufbuild/protovalidate-go"
	greetpb "github.com/jadeidev/grpc-go-course/greet-buf-validate/gen/go/greet/v1"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/reflection"
)

type Server struct {
	greetpb.UnimplementedGreetServiceServer
	// lets add the standard grpc health check service, the proto file for this is all on the grpc side
	grpchealth.UnimplementedHealthServer
	healthMu  sync.RWMutex
	statusMap map[string]grpchealth.HealthCheckResponse_ServingStatus
}

func NewServer() (*grpc.Server, error) {
	// create a server
	opts := []grpc.ServerOption{}
	tls := false
	if tls {
		certFile := "../../ssl/server.crt"
		keyFile := "../../ssl/server.pem"
		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
		if sslErr != nil {
			log.Fatalf("Failed loading certificates: %v", sslErr)
			return nil, sslErr
		}
		opts = append(opts, grpc.Creds(creds))
	}

	s := grpc.NewServer(opts...)
	server := &Server{
		statusMap: map[string]grpchealth.HealthCheckResponse_ServingStatus{
			"":      grpchealth.HealthCheckResponse_SERVING,
			"greet": grpchealth.HealthCheckResponse_SERVING,
		},
	}
	greetpb.RegisterGreetServiceServer(s, server)
	// just like the greet service, we need to register the health service
	grpchealth.RegisterHealthServer(s, server)
	reflection.Register(s)
	return s, nil
}

// see how validate is used in the Greet function
func (*Server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v\n", req)
	// validate input
	v, err := protovalidate.New()
	if err != nil {
		fmt.Println("failed to initialize validator:", err)
	}
	if err = v.Validate(req); err != nil {
		fmt.Println("validation failed:", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	} else {
		fmt.Println("validation succeeded")
	}
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	// for proto messages marshal and unmarshal only use protojson
	marshaler := protojson.MarshalOptions{
		Multiline:         true,
		EmitUnpopulated:   true,
		EmitDefaultValues: true,
	}

	jsonReq, _ := marshaler.Marshal(req)
	fmt.Println("Request in Json format: \n", string(jsonReq))
	return res, nil
}

// see how validate is used in the GreetManyTimes function
func (*Server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function was invoked with %v\n", req)
	v, err := protovalidate.New()
	if err != nil {
		fmt.Println("failed to initialize validator:", err)
	}
	if err = v.Validate(req); err != nil {
		fmt.Println("validation failed:", err)
	} else {
		fmt.Println("validation succeeded")
	}
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*Server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet function was invoked with a streaming request\n")
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we have finished reading the client stream
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		result += "Hello " + firstName + "! "
	}
}

func (*Server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("GreetEveryone function was invoked with a streaming request\n")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}
		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName + "! "

		sendErr := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v", sendErr)
			return sendErr
		}
	}

}

func (*Server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	fmt.Printf("GreetWithDeadline function was invoked with %v\n", req)
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.DeadlineExceeded {
			// the client canceled the request
			fmt.Println("The client canceled the request!")
			return nil, status.Error(codes.Canceled, "the client canceled the request")
		}
		time.Sleep(1 * time.Second)
	}
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetWithDeadlineResponse{
		Result: result,
	}
	return res, nil
}

// implement the health check service
func (s *Server) Check(ctx context.Context, in *grpchealth.HealthCheckRequest) (*grpchealth.HealthCheckResponse, error) {
	if in.Service == "greet" {
		// perform basic greet to validate things are working
		_, err := s.Greet(ctx, &greetpb.GreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Health",
				LastName:  "Check",
			},
		})
		if err != nil {
			s.SetServiceStatus(in.Service, grpchealth.HealthCheckResponse_NOT_SERVING)
		} else {
			s.SetServiceStatus(in.Service, grpchealth.HealthCheckResponse_SERVING)
		}
	}
	s.healthMu.RLock()
	defer s.healthMu.RUnlock()
	if serviceStatus, ok := s.statusMap[in.Service]; ok {
		return &grpchealth.HealthCheckResponse{
			Status: serviceStatus,
		}, nil
	}
	return nil, status.Error(codes.NotFound, "unknown service")
}

func (s *Server) SetServiceStatus(service string, status grpchealth.HealthCheckResponse_ServingStatus) {
	s.healthMu.Lock()
	defer s.healthMu.Unlock()
	s.statusMap[service] = status
}

func main() {
	fmt.Println("Hello world")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s, err := NewServer()
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
