package grpc

import (
	"context"

	"google.golang.org/grpc"
)

// NewClientFunc - функция регистрации нового клиента.
type NewClientFunc[T any] func(cc grpc.ClientConnInterface) T

// NewClient - возвращает новый клиент.
//
//nolint:ireturn,nolintlint // должен быть интерфейс
func NewClient[T any](
	connection *Connection,
	clientFunc NewClientFunc[T],
	interceptors ...grpc.UnaryClientInterceptor,
) T {
	if len(interceptors) == 0 {
		return clientFunc(connection)
	}

	wrappedConnection := &ConnectionWithUnaryClientInterceptor{
		Connection: connection,
		invokeFunc: chainInvoke(interceptors, connection.Invoke),
	}

	return clientFunc(wrappedConnection)
}

// chaiInvokerFunc - генерирует рекурсивный вызов для slice grpc.UnaryClientInterceptor.
func chainInvokerFunc(interceptors []grpc.UnaryClientInterceptor, curr int, finalInvoker InvokerFunc) grpc.UnaryInvoker {
	if curr == len(interceptors)-1 {
		return func(ctx context.Context, method string, in, out any, _ *grpc.ClientConn, opts ...grpc.CallOption) error {
			return finalInvoker(ctx, method, in, out, opts...)
		}
	}

	return func(ctx context.Context, method string, in, out any, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return interceptors[curr+1](ctx, method, in, out, cc, chainInvokerFunc(interceptors, curr+1, finalInvoker), opts...)
	}
}

// chainInvoke - возвращает цепочку вызовов из grpc.UnaryClientInterceptor для grpc.UnaryInvoker.
func chainInvoke(interceptors []grpc.UnaryClientInterceptor, invokerFunc InvokerFunc) grpc.UnaryInvoker {
	nextInterceptorFunc := chainInvokerFunc(interceptors, 0, invokerFunc)
	return func(ctx context.Context, method string, in, out any, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return interceptors[0](ctx, method, in, out, cc, nextInterceptorFunc, opts...)
	}
}
