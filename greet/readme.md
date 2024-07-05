cd into this folder where readme is:
generate with protoc
`protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative greetpb/greet.proto`

to start server:
`go run greet_server/server.go`

to start client:
`go run greet_client/client.go`