
cd into the greet-gateway folder
`buf dep update`
`buf generate`

to start grpc server:
`go run greet_server/server.go`

to start grpc gateway server:
`go run greet_gateway_server/server.go`

to start client:
`go run greet_client/client.go`
`curl -X POST http://localhost:8081/v1/greet \
     -H "Content-Type: application/json" \
     -d '{"greeting": {"firstName": "John", "lastName": "Doe"}}'