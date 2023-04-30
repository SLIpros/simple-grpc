package grpc

import (
	"context"
	"net"
	"testing"

	"go.nhat.io/grpcmock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

// ContextDialer номеронабиратель для создания соединений.
type ContextDialer = func(context.Context, string) (net.Conn, error)

// NewMockServer mock grpc server.
func NewMockServer(t *testing.T, services ...Service) ContextDialer {
	t.Helper()

	const bufSize = 1024 * 1024
	listener := bufconn.Listen(bufSize)

	server := grpc.NewServer()
	t.Cleanup(server.GracefulStop)

	// регистрируем сервисы
	registerServices(nil, server, services)

	go func() {
		if err := server.Serve(listener); err != nil {
			t.Errorf("unable to start grpc server %v", err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

// InvokeUnary вызывает унарный запрос.
func InvokeUnary[Response any](
	ctx context.Context,
	dialer ContextDialer,
	method string,
	request any,
) (*Response, error) {
	options := []grpcmock.InvokeOption{
		grpcmock.WithContextDialer(dialer),
		grpcmock.WithInsecure(),
	}

	return invokeUnary[Response](ctx, method, request, options)
}

// InvokeUnaryWithHeaders вызывает унарный запрос с указанными заголовками.
func InvokeUnaryWithHeaders[Response any](
	ctx context.Context,
	dialer ContextDialer,
	method string,
	request any,
	headers map[string]string,
) (*Response, error) {
	options := []grpcmock.InvokeOption{
		grpcmock.WithContextDialer(dialer),
		grpcmock.WithInsecure(),
	}

	for key, value := range headers {
		options = append(options, grpcmock.WithHeader(key, value))
	}

	return invokeUnary[Response](ctx, method, request, options)
}

// invokeUnary вызывает унарный запрос.
func invokeUnary[Response any](
	ctx context.Context,
	method string,
	request any,
	options []grpcmock.InvokeOption,
) (*Response, error) {
	var response Response
	if err := grpcmock.InvokeUnary(ctx, method, request, &response, options...); err != nil {
		return nil, err
	}

	return &response, nil
}
