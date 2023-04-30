package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	// MetaDataFieldAuthorization - ключ авторизации из метаданных.
	MetaDataFieldAuthorization = "authorization"
	// BearerPrefix - префикс перед токеном аутентификации/авторизации в заголовке.
	BearerPrefix = "Bearer "
)

// AuthorizeWithBearerToken - записывает bearer token в метаданные подключения.
func AuthorizeWithBearerToken(token string) grpc.UnaryClientInterceptor {
	metaDataValue := BearerPrefix + token
	return func(
		ctx context.Context,
		method string,
		req,
		reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		ctxWithAuthorization := metadata.AppendToOutgoingContext(ctx, MetaDataFieldAuthorization, metaDataValue)
		return invoker(ctxWithAuthorization, method, req, reply, cc, opts...)
	}
}
