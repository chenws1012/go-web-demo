package service

import "context"

type HelloService interface {
	SayHello(ctx context.Context) (string, error)
}

type helloService struct{}

func NewHelloService() HelloService {
	return &helloService{}
}

func (s *helloService) SayHello(ctx context.Context) (string, error) {
	return "Hello, Enterprise Go Web!", nil
}
