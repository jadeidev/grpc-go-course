package main

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/credentials/insecure"

	greetpb "github.com/jadeidev/grpc-go-course/greet-observability/gen/greet/v1"

	"google.golang.org/grpc"
)

func InitTracer() (*sdktrace.TracerProvider, error) {
	exporter, err := otlptracegrpc.New(
		context.Background(),
		otlptracegrpc.WithEndpointURL("https://localhost:8200"),
	)
	/* can also use use the http exporter, one use case would be to export to http istead of https (grpc requires https)
		for this import "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp" 
		if endpointurl ends with extra path (e.g http://localhost:8200/apmsink), we need to add the WithURLPath to specify the path to send traces
		```
		otlptracehttp.WithEndpoint("localhost:8200")
		otlptracehttp.WithURLPath("/apm-sink/v1/traces")
		```
		or can also do
		```
		otlptracehttp.WithEndpointURL("http://localhost:8200/apm-sink/v1/traces")
		```
	*/ 
	// exporter, err := otlptracehttp.New(
	// 	context.Background(),
	// 	otlptracehttp.WithEndpointURL("http://localhost:8200"),
	// )


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

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	// creating span here is critical to assure that elastic service map will work in the same way it does with 
	tracer := otel.Tracer("")
	ctx, span := tracer.Start(
		context.Background(),
		"greet.v1.GreetService",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("rpc.system", "grpc"),
			attribute.String("server.address", "localhost"),
			attribute.String("server.port", "50051"),
		),
	)
	defer span.End()
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Stephane",
			LastName:  "Maarek",
		},
	}
	res, err := c.Greet(ctx, req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}


func main() {
	tp, err := InitTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	}

	cc, err := grpc.NewClient("localhost:50051", opts...)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	doUnary(c)

}

