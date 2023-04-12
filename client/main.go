package main

import (
	"client/services/hello"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"log"
)

func main() {
	tracer.Start(
		tracer.WithService("grpc-client"),
		tracer.WithEnv("dev"),
	)
	defer tracer.Stop()
	// Create a gRPC client connection with the Datadog middleware.
	dd := grpctrace.UnaryClientInterceptor(grpctrace.WithServiceName("grpc-client"))

	conn, err := grpc.Dial(
		"localhost:3000",
		grpc.WithUnaryInterceptor(dd),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	// Create a gRPC client.
	client := hello.NewHelloServiceClient(conn)

	// Use the gRPC client.
	response, err := client.Hello(context.Background(), &hello.HelloRequest{
		Message: "Hello world",
	})
	if err != nil {
		log.Fatalf("failed to call MyMethod: %v", err)
	}

	log.Printf("response: %v", response)
}
