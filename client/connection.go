package grpc

import (
	"context"
	"crypto/tls"

	grpcLogging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
)

type Logger interface {
	ReplaceGRPCLogger()
	UnaryClientInterceptor() grpc.UnaryClientInterceptor
	PayloadUnaryClientInterceptor(grpcLogging.ClientPayloadLoggingDecider) grpc.UnaryClientInterceptor
}

// NewConnection возвращает новое подключение.
func NewConnection(
	options *ConnectionOptions,
	logger Logger,
	interceptors ...grpc.UnaryClientInterceptor,
) (*Connection, error) {
	if options == nil {
		return nil, errors.New("invalid arguments")
	}

	if logger == nil {
		logger = &dummyLogger{}
	}

	logger.ReplaceGRPCLogger()

	// включаем сохранение времени исполнения запроса в метрики
	grpcPrometheus.EnableClientHandlingTimeHistogram()

	systemInterceptors := []grpc.UnaryClientInterceptor{
		// tracing
		otelgrpc.UnaryClientInterceptor(),
		// считаем метрики
		grpcPrometheus.UnaryClientInterceptor,
		// логируем исполнение
		logger.UnaryClientInterceptor(),
	}

	// логируем тело запроса
	if options.LogBody {
		hasMethods := len(options.LogBodyMethods) > 0
		methods := make(map[string]struct{}, len(options.LogBodyMethods))
		for i := range options.LogBodyMethods {
			methods[options.LogBodyMethods[i]] = struct{}{}
		}

		systemInterceptors = append(systemInterceptors,
			logger.PayloadUnaryClientInterceptor(func(ctx context.Context, fullMethodName string) bool {
				if !hasMethods {
					// так как у нас нет ограничений, то логируем тело всех запросов
					return true
				}

				_, found := methods[fullMethodName]
				return found
			}),
		)
	}

	dialOptions := []grpc.DialOption{grpc.WithChainUnaryInterceptor(systemInterceptors...)}

	if options.Insecure {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(
			credentials.NewTLS(&tls.Config{
				//nolint:gosec //intentional insecure
				InsecureSkipVerify: options.TLSSkipVerify,
			}),
		))
	}

	if len(interceptors) > 0 {
		dialOptions = append(dialOptions, grpc.WithChainUnaryInterceptor(interceptors...))
	}

	clientConn, err := grpc.Dial(options.Addr, dialOptions...)
	if err != nil {
		return nil, err
	}

	connection := Connection{
		ClientConn: clientConn,
	}

	return &connection, nil
}

// Connection - подключение.
type Connection struct {
	*grpc.ClientConn
}
