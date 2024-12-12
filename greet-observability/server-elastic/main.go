package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/url"

	greetpb "github.com/jadeidev/grpc-go-course/greet-observability/gen/greet/v1"

	apmgrpc "go.elastic.co/apm/module/apmgrpc/v2"
	"go.elastic.co/apm/v2"
	"go.elastic.co/apm/v2/transport"

	"google.golang.org/grpc"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v\n", req)
	// demo additional downstream api request when this rpc is called
	tx := apm.TransactionFromContext(ctx)
	if tx == nil {
		tx = apm.DefaultTracer().StartTransaction("HTTP GET", "request")
	}
	defer tx.End()

	span := tx.StartSpan("GET https://jsonplaceholder.typicode.com/posts/1", "external.http", nil)
	defer span.End()
	ctx = apm.ContextWithSpan(ctx, span)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", "https://jsonplaceholder.typicode.com/posts/1", nil)
	if err != nil {
		return nil, err
	}
	client := apmhttp.WrapClient(http.DefaultClient)
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Response from external API: %s\n", body)

	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func SetupAPM() (*apm.Tracer, error) {
	apmServerURL, err := url.Parse("https://localhost:8200")
	if err != nil {
		return nil, err
	}

	transport, err := transport.NewHTTPTransport(transport.HTTPTransportOptions{
		ServerURLs: []*url.URL{apmServerURL},
	})
	if err != nil {
		return nil, err
	}

	tracer, err := apm.NewTracerOptions(apm.TracerOptions{
		ServiceName: "imconsole-grpc-greeter-server-elastic",
		Transport:   transport,
	})
	if err != nil {
		return nil, err
	}

	return tracer, nil
}

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	//setup apm
	tracer, err := SetupAPM()
	if err != nil {
		log.Fatalf("Failed to setup APM: %v", err)
	}
	defer tracer.Flush(nil)
	apm.SetDefaultTracer(tracer)

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(apmgrpc.NewUnaryServerInterceptor(apmgrpc.WithTracer(tracer))),
	)
	greetpb.RegisterGreetServiceServer(s, &server{})

	fmt.Println("Server started")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
