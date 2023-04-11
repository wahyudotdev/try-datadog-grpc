package main

import (
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpclogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"log"
	"net"
	"time"
	"try-datadog-grpc/services/hello"
)

func main() {
	tracer.Start(
		tracer.WithService("grpc-server"),
		tracer.WithEnv("dev"),
	)
	defer tracer.Stop()

	var logrusEntry = logrus.NewEntry(logrus.New())
	var logrusOpts = []grpclogrus.Option{
		grpclogrus.WithDurationField(func(duration time.Duration) (key string, value interface{}) {
			return "grpc.time_s", duration.Nanoseconds()
		}),
	}

	dd := grpctrace.UnaryServerInterceptor(grpctrace.WithServiceName("grpc-server"))
	registrar := grpc.NewServer(grpc.UnaryInterceptor(middleware.ChainUnaryServer(
		dd,
		grpclogrus.UnaryServerInterceptor(logrusEntry, logrusOpts...),
	)))
	l, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := hello.NewHelloService()
	hello.RegisterHelloServiceServer(registrar, server)
	log.Println("grpc server started")
	log.Fatal(registrar.Serve(l))
}
