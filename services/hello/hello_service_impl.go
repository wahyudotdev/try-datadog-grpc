package hello

import (
	context "context"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type ServiceImpl struct {
}

func (h ServiceImpl) Hello(ctx context.Context, request *HelloRequest) (*HelloResponse, error) {
	tracer.SpanFromContext(ctx)
	return &HelloResponse{
		Message: "Hello",
	}, nil
}

func (h ServiceImpl) mustEmbedUnimplementedHelloServiceServer() {
	//TODO implement me
	panic("implement me")
}

func NewHelloService() HelloServiceServer {
	return ServiceImpl{}
}
