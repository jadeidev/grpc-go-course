package main

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	grpchealth "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/test/bufconn"

	greetpb "github.com/jadeidev/grpc-go-course/greet-buf-validate/gen/go/greet/v1"
)

// this is a test function to demo how to use the protovalidate library
func TestValidate(t *testing.T) {
	// create a Greeting object (as if it was received in a gRPC request)
	greeting := &greetpb.Greeting{}

	v, err := protovalidate.New()
	if err != nil {
		fmt.Println("failed to initialize validator:", err)
	}

	// validation will fail because FirstName and LastName are required
	if err = v.Validate(greeting); err != nil {
		fmt.Println("validation failed:", err)
	} else {
		fmt.Println("validation succeeded")
	}
	greeting = &greetpb.Greeting{FirstName: "Stephane", LastName: "Maarek"}
	// validation will succeed because FirstName and LastName are provided
	if err = v.Validate(greeting); err != nil {
		fmt.Println("validation failed:", err)
	} else {
		fmt.Println("validation succeeded")
	}
}

// here we want to test the server implementation using golang test suite
// (rather than using the client we created greet_client/client.go)

// first we wanna simulate server
func setupServer(t *testing.T) (*grpc.Server, *bufconn.Listener) {
	// create a server
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})
	srv, err := NewServer()
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}
	t.Cleanup(func() {
		srv.Stop()
	})
	go func() {
		if err := srv.Serve(lis); err != nil {
			panic(fmt.Sprintf("Server exited with error: %v", err))
		}
	}()
	return srv, lis
}

// now we wanna simulate client
func setupClient(t *testing.T) (context.Context, *grpc.ClientConn) {
	_, lis := setupServer(t)
	ctx := context.Background()
	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
	cc, err := grpc.NewClient(
		"passthrough://bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	t.Cleanup(func() {
		cc.Close()
	})
	return ctx, cc
}

func TestGreetService(t *testing.T) {
	ctx, cc := setupClient(t)
	c := greetpb.NewGreetServiceClient(cc)

	// test the Greet RPC
	t.Run("Greet", func(t *testing.T) {
		t.Run("valid", func(t *testing.T) {
			req := &greetpb.GreetRequest{
				Greeting: &greetpb.Greeting{
					FirstName: "Stephane",
					LastName:  "Maarek",
				},
			}
			res, err := c.Greet(ctx, req)
			if err != nil {
				t.Fatalf("error while calling Greet RPC: %v", err)
			}
			if res.Result != "Hello Stephane" {
				t.Fatalf("unexpected result: %v", res.Result)
			}
		})
		t.Run("invalid", func(t *testing.T) {
			req := &greetpb.GreetRequest{
				Greeting: &greetpb.Greeting{
					FirstName: "St",
					LastName:  "Maarek",
				},
			}
			_, err := c.Greet(ctx, req)
			if err == nil {
				t.Fatalf("expected error but got nil")
			}
		})
	})
}

func TestHealthcheck(t *testing.T) {
	// we are going to simulate what health probe does in kubernetes
	_, cc := setupClient(t)
	c := grpchealth.NewHealthClient(cc)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := c.Check(ctx, &grpchealth.HealthCheckRequest{})
	if err != nil {
		t.Fatalf("error while calling Check RPC: %v", err)
	} else if resp.Status != grpchealth.HealthCheckResponse_SERVING {
		t.Fatalf("unexpected status: %v", resp.Status)
	}
	resp, err = c.Check(ctx, &grpchealth.HealthCheckRequest{Service: "greet"})
	if err != nil {
		t.Fatalf("error while calling Check RPC: %v", err)
	} else if resp.Status != grpchealth.HealthCheckResponse_SERVING {
		t.Fatalf("unexpected status: %v", resp.Status)
	}
}
