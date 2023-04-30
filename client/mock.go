package grpc

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.nhat.io/grpcmock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ServerMock - функция mock сервера.
type ServerMock func(s *grpcmock.Server)

// NewMockClient - создает новый mock client.
//
//nolint:ireturn,nolintlint // должен быть интерфейс
func NewMockClient[T any](t *testing.T, serviceRegistrarFunc any, clientFunc NewClientFunc[T], mockFunc ServerMock) T {
	t.Helper()

	require.NotNil(t, serviceRegistrarFunc, "nil service registrar func")
	require.NotNil(t, clientFunc, "nil client func")

	options := []grpcmock.ServerOption{
		grpcmock.RegisterService(serviceRegistrarFunc),
	}

	if mockFunc != nil {
		options = append(options, grpcmock.ServerOption(mockFunc))
	}

	_, d := grpcmock.MockServerWithBufConn(options...)(t)

	dial, err := grpc.Dial(
		"",
		grpc.WithContextDialer(d),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err, "unable to execute grpc dial")

	return clientFunc(dial)
}
