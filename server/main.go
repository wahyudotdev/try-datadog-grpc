package main

import (
	"context"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpclogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"log"
	"net"
	"server/services/hello"
	"time"
)

func reqBodyInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	span, _ := tracer.StartSpanFromContext(ctx, info.FullMethod, tracer.SpanType("grpc"),
		tracer.ServiceName("grpc-server"),
		tracer.ResourceName(info.FullMethod))
	span.SetTag("body", req)

	span.Finish()
	return handler(ctx, req)
}

func main() {
	tracer.Start(
		tracer.WithService("grpc-server"),
		tracer.WithEnv("dev"),
	)
	defer tracer.Stop()

	logrusEntry := logrus.NewEntry(logrus.New())
	logrusOpts := []grpclogrus.Option{
		grpclogrus.WithDurationField(func(duration time.Duration) (key string, value interface{}) {
			return "grpc.time_s", duration.Nanoseconds()
		}),
	}

	dd := grpctrace.UnaryServerInterceptor(grpctrace.WithServiceName("grpc-server"))
	server := grpc.NewServer(grpc.UnaryInterceptor(middleware.ChainUnaryServer(
		dd,
		grpclogrus.UnaryServerInterceptor(logrusEntry, logrusOpts...),
		reqBodyInterceptor,
	)))
	l, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	service := hello.NewHelloService()
	reflection.Register(server)
	hello.RegisterHelloServiceServer(server, service)
	log.Println("grpc server started")
	log.Fatal(server.Serve(l))
}
