package main

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"go.elastic.co/apm/module/apmgrpc/v2"
	"go.elastic.co/apm/v2"
	"go.elastic.co/apm/v2/transport"

	greetpb "github.com/jadeidev/grpc-go-course/greet-observability/gen/greet/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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
		ServiceName: "imconsole-grpc-greeter-client-elastic",
		Transport:   transport,
	})
	if err != nil {
		return nil, err
	}
	return tracer, nil

}

func main() {
	//setup apm
	tracer, err := SetupAPM()
	if err != nil {
		log.Fatalf("Failed to setup APM: %v", err)
	}
	defer tracer.Flush(nil)
	apm.SetDefaultTracer(tracer)

	cc, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(apmgrpc.NewUnaryClientInterceptor()),
	)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	doUnary(c)

}

func doUnary(c greetpb.GreetServiceClient) {
	tx := apm.DefaultTracer().StartTransaction("main", "request")
	defer tx.End()
	ctx := apm.ContextWithTransaction(context.Background(), tx)
	fmt.Println("Starting to do a Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Stephane",
			LastName:  "Maarek",
		},
	}
	res, err := c.Greet(ctx, req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}
