package main

import (
	"fmt"

	"github.com/bufbuild/protovalidate-go"

	greetpb "github.com/jadeidev/grpc-go-course/greet-buf-validate/gen/go/greet/v1"
)

func main() {
	greeting := &greetpb.Greeting{}

	v, err := protovalidate.New()
	if err != nil {
		fmt.Println("failed to initialize validator:", err)
	}

	if err = v.Validate(greeting); err != nil {
		fmt.Println("validation failed:", err)
	} else {
		fmt.Println("validation succeeded")
	}
	greeting = &greetpb.Greeting{FirstName: "Stephane", LastName: "Maarek"}
	if err = v.Validate(greeting); err != nil {
		fmt.Println("validation failed:", err)
	} else {
		fmt.Println("validation succeeded")
	}
}
