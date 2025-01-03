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
		resource.WithAttributes(
			attribute.Key("service.name").String("grpc-greeter-client-otel"), // Add custom resource attributes.
			attribute.Key("deployment.environment").String("development"),    // this is specifically for elastic so that env is properly mapped
		),
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
	// when using batcher as the exporter we need to shutdown the exporter to assure that all traces are sent
	// this is because the batcher will not send traces until the buffer is full or the exporter is shutdown or batcher timeout is reached
	// for testing purposes we just do this here. an alternative would be to use Syncer exporter instead of Batcher
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		/*
			The attribute `attribute.Key("http.url")`` would populate the elastic
			span attribute `span.destination.service.name`. it is optional!
			its kinda weird that we need to put it here given that we start a span with the doUnary function
			ideally we would put it in the doUnary span, but that didnt seem to have affected anything
		*/
		grpc.WithStatsHandler(otelgrpc.NewClientHandler(
			otelgrpc.WithSpanAttributes(
				attribute.Key("http.url").String("http://127.0.0.1:50051"),
				attribute.Key("deployment.environment").String("development"), // this is specifically for elastic so that env is properly mapped
			),
		)),
	}

	cc, err := grpc.NewClient("localhost:50051", opts...)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	doUnary(c)

}
