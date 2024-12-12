package main

import (
	"context"
	"fmt"
	"log"
	"net"

	greetpb "github.com/jadeidev/grpc-go-course/greet-observability/gen/greet/v1"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v\n", req)
	// Add span attributes for better visualization (not mandatory)
	tracer := otel.Tracer("")
	_, span := tracer.Start(
		ctx,
		"greet.v1.GreetService",
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(
			attribute.String("rpc.system", "grpc"),
			attribute.String("server.address", "localhost"),
			attribute.String("server.port", "50051"),
		),
	)
	defer span.End()
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func initTracer() (*sdktrace.TracerProvider, error) {
	exporter, err := otlptracegrpc.New(
		context.Background(),
		otlptracegrpc.WithEndpointURL("https://localhost:8200"),
	)
	// can also use use the http exporter, one use case would be to export to http istead of https (grpc requires https)
	// for this import "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp" 
	// exporter, err := otlptracehttp.New(
	// 	context.Background(),
	// 	otlptracehttp.WithEndpointURL("http://localhost:8200"),
	if err != nil {
		return nil, err
	}


	resource, err := resource.New(
		context.Background(),
		resource.WithTelemetrySDK(), // Discover and provide information about the OpenTelemetry SDK used.
		resource.WithProcess(),      // Discover and provide process information.
		resource.WithOS(),           // Discover and provide OS information.
		resource.WithContainer(),    // Discover and provide container information.
		resource.WithHost(),         // Discover and provide host information.
		resource.WithAttributes(attribute.String("service.name", "grpc-greeter-server-otel")), // Add custom resource attributes.
	)
	if err != nil {
		log.Fatalln(err) // The error may be fatal.
	}
	// can also do this manually
	// resource := resource.NewWithAttributes(
	// 	semconv.SchemaURL,
	// 	semconv.ServiceName("grpc-greeter-server-otel"),
	// 	semconv.DeploymentEnvironmentName("development"),
	// 	semconv.ServiceVersion("1.0.0"),
	// )


	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	return tp, nil
}

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	tp, err := initTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
	greetpb.RegisterGreetServiceServer(s, &server{})

	fmt.Println("Server started")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
