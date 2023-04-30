package interceptor

import (
	"context"
	"strings"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// Service исполняет сервисные унарные интерцепторы.
func Service(descriptor *grpc.ServiceDesc, interceptors []grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	chain := grpcMiddleware.ChainUnaryServer(interceptors...)
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		service, _, err := splitMethodName(info.FullMethod)
		if err != nil {
			return nil, err
		}

		if service != descriptor.ServiceName {
			return handler(ctx, req)
		}

		return chain(ctx, req, info, handler)
	}
}

// splitMethodName разделяет название сервиса от метода.
//
//nolint:gocritic // не нужен тут именованный возврат
func splitMethodName(fullMethod string) (string, string, error) {
	fullMethod = strings.TrimPrefix(fullMethod, "/")
	service, method, found := strings.Cut(fullMethod, "/")
	if !found {
		return service, method, errors.New("unable to split method name")
	}

	return service, method, nil
}
