package hello

import (
	context "context"
)

type ServiceImpl struct {
}

func (h ServiceImpl) Hello(ctx context.Context, request *HelloRequest) (*HelloResponse, error) {
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
