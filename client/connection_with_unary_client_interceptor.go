package grpc

import (
	"context"

	"google.golang.org/grpc"
)

// InvokerFunc - функция вызова.
type InvokerFunc func(ctx context.Context, method string, in, out any, opts ...grpc.CallOption) error

// ConnectionWithUnaryClientInterceptor - подключение с grpc.UnaryClientInterceptor.
type ConnectionWithUnaryClientInterceptor struct {
	*Connection

	invokeFunc grpc.UnaryInvoker
}

// Invoke - исполняет вызов.
func (c *ConnectionWithUnaryClientInterceptor) Invoke(
	ctx context.Context,
	method string,
	in,
	out any,
	opts ...grpc.CallOption,
) error {
	return c.invokeFunc(ctx, method, in, out, c.Connection.ClientConn, opts...)
}
