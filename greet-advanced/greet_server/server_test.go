package main

import (
	"fmt"
	"testing"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc"
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

// here we want to test the server implementation
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
