package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	greetpb "github.com/jadeidev/grpc-go-course/greet-gateway/gen/go/greet/v1"

	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Hello I'm a client")

	tls := false
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	if tls {
		certFile := "../../ssl/ca.crt" // Certificate Authority Trust certificate
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Fatalf("Error while loading CA trust certificate: %v", sslErr)
			return
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	cc, err := grpc.NewClient("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	// fmt.Printf("Created client: %f", c)

	doUnary(c)
	doUnaryWithRestGateway()

}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC direct call to the gRPC server...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Stephane",
			LastName:  "Maarek",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}

func doUnaryWithRestGateway() {
	fmt.Println("Starting to do a Unary RPC with REST API gateway...")
	url := "http://localhost:8081/v1/greet"
	payload := bytes.NewBuffer([]byte(`{"greeting": {"firstName": "John", "lastName": "Doe"}}`))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	fmt.Println("Response:", string(body))
}
