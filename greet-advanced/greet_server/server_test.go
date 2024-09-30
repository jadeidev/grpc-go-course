package main

import (
	"fmt"
	"testing"

	"github.com/bufbuild/protovalidate-go"

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
