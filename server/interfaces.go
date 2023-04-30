package grpc

import (
	grpcLogging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	"google.golang.org/grpc"
)

type Logger interface {
	ReplaceGRPCLogger()
	UnaryServerInterceptor() grpc.UnaryServerInterceptor
	PayloadUnaryServerInterceptor(grpcLogging.ServerPayloadLoggingDecider) grpc.UnaryServerInterceptor
	Info(...any)
	Debugf(string, ...any)
}

type PayloadLoggingDecider interface {
	// PayloadLoggingDecider решает логгировать тела метода или нет.
	PayloadLoggingDecider(fullMethod string) bool
}
