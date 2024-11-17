
cd into the greet-gateway folder
`buf dep update`
`buf generate`

to start grpc server:
`go run greet_server/server.go`

to start grpc gateway server:
`go run greet_gateway_server/gateway_server.go`

to start client:
`go run greet_client/client.go`

to make a post to gateway
```
curl -X POST http://localhost:8081/v1/greet \
     -H "Content-Type: application/json" \
     -d '{"greeting": {"firstName": "John", "lastName": "Doe"}}'
```

Use with swagger-ui:
`docker run --rm -p 80:8080 -e SWAGGER_JSON=/tmp/oapi.yaml --mount type=bind,source=/Users/<user>/source/grpc-go-course/greet-gateway/gen/openapiv31/greet/v1/greet.openapi.yaml,target=/tmp/oapi.yaml swaggerapi/swagger-ui`


We also demo generation of Open API 3.1 spec using `github.com/sudorandom/protoc-gen-connect-openapi`