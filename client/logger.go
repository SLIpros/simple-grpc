package grpc

import (
	"context"

	grpcLogging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	"google.golang.org/grpc"
)

type dummyLogger struct{}

func (d *dummyLogger) ReplaceGRPCLogger() {}
func (d *dummyLogger) UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req,
		reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
func (d *dummyLogger) PayloadUnaryClientInterceptor(
	grpcLogging.ClientPayloadLoggingDecider,
) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req,
		reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
