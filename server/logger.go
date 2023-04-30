package grpc

import (
	"context"

	grpcLogging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	"google.golang.org/grpc"
)

type dummyLogger struct{}

func (d *dummyLogger) PayloadUnaryServerInterceptor(
	grpcLogging.ServerPayloadLoggingDecider,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		return handler(ctx, req)
	}
}

func (d *dummyLogger) ReplaceGRPCLogger() {}
func (d *dummyLogger) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		return handler(ctx, req)
	}
}
func (d *dummyLogger) Info(...any)           {}
func (d *dummyLogger) Debugf(string, ...any) {}
