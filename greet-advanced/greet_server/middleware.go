package main

import (
	"context"

	"google.golang.org/grpc"
)

// example for interceptors
// we are going to inject some variable to the context
// and then we are going to read it in the server interceptor

type contextKey string

const myContextKey contextKey = "my_context_key"

func contextInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// inject some value to the context
		ctx = context.WithValue(ctx, myContextKey, "some_value")
		return handler(ctx, req)
	}
}
