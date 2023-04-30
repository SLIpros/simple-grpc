package grpc

import (
	"google.golang.org/grpc"

	"github.com/SLIpros/simple-grpc/server/interceptor"
)

// Service сервис.
type Service interface {
	// Descriptor возвращает дескриптор.
	Descriptor() *grpc.ServiceDesc
}

// Services сервисы.
type Services []Service

// NewServiceWithUnaryInterceptors возвращает новый сервис с унарными интерцепторами.
func NewServiceWithUnaryInterceptors(service Service, interceptors ...grpc.UnaryServerInterceptor) *ServiceWithUnaryInterceptors {
	return &ServiceWithUnaryInterceptors{
		Service:      service,
		interceptors: interceptors,
	}
}

// ServiceWithUnaryInterceptors сервис с унарными интерцепторами.
type ServiceWithUnaryInterceptors struct {
	Service

	interceptors []grpc.UnaryServerInterceptor
}

// Interceptors возвращает интерцепторы.
func (s *ServiceWithUnaryInterceptors) Interceptors() []grpc.UnaryServerInterceptor {
	return s.interceptors
}

// collectServiceUnaryInterceptors возвращает сервисные унарные интерцепторы.
func collectServiceUnaryInterceptors(services Services) []grpc.UnaryServerInterceptor {
	//nolint:prealloc // не могу сделать prealloc
	var interceptors []grpc.UnaryServerInterceptor
	for _, service := range services {
		unary, ok := service.(*ServiceWithUnaryInterceptors)
		if !ok {
			continue
		}

		interceptors = append(interceptors, interceptor.Service(service.Descriptor(), unary.Interceptors()))
	}

	return interceptors
}
