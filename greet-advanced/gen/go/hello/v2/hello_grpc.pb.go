// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: hello/v2/hello.proto

package hellov2

import (
	context "context"
	v1 "github.com/jadeidev/grpc-go-course/greet-advanced/gen/go/greet/v1"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	HelloService_Greet_FullMethodName = "/hello.v2.HelloService/Greet"
)

// HelloServiceClient is the client API for HelloService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HelloServiceClient interface {
	// Unary
	Greet(ctx context.Context, in *v1.GreetRequest, opts ...grpc.CallOption) (*v1.GreetResponse, error)
}

type helloServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewHelloServiceClient(cc grpc.ClientConnInterface) HelloServiceClient {
	return &helloServiceClient{cc}
}

func (c *helloServiceClient) Greet(ctx context.Context, in *v1.GreetRequest, opts ...grpc.CallOption) (*v1.GreetResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(v1.GreetResponse)
	err := c.cc.Invoke(ctx, HelloService_Greet_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HelloServiceServer is the server API for HelloService service.
// All implementations must embed UnimplementedHelloServiceServer
// for forward compatibility.
type HelloServiceServer interface {
	// Unary
	Greet(context.Context, *v1.GreetRequest) (*v1.GreetResponse, error)
	mustEmbedUnimplementedHelloServiceServer()
}

// UnimplementedHelloServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedHelloServiceServer struct{}

func (UnimplementedHelloServiceServer) Greet(context.Context, *v1.GreetRequest) (*v1.GreetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Greet not implemented")
}
func (UnimplementedHelloServiceServer) mustEmbedUnimplementedHelloServiceServer() {}
func (UnimplementedHelloServiceServer) testEmbeddedByValue()                      {}

// UnsafeHelloServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HelloServiceServer will
// result in compilation errors.
type UnsafeHelloServiceServer interface {
	mustEmbedUnimplementedHelloServiceServer()
}

func RegisterHelloServiceServer(s grpc.ServiceRegistrar, srv HelloServiceServer) {
	// If the following call pancis, it indicates UnimplementedHelloServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&HelloService_ServiceDesc, srv)
}

func _HelloService_Greet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.GreetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HelloServiceServer).Greet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HelloService_Greet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HelloServiceServer).Greet(ctx, req.(*v1.GreetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// HelloService_ServiceDesc is the grpc.ServiceDesc for HelloService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var HelloService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "hello.v2.HelloService",
	HandlerType: (*HelloServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Greet",
			Handler:    _HelloService_Greet_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "hello/v2/hello.proto",
}