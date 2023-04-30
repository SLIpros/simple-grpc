package grpc

import (
	"context"
	"net"

	grpcLogging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcCTXTags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Logger interface {
	ReplaceGRPCLogger()
	UnaryServerInterceptor() grpc.UnaryServerInterceptor
	PayloadUnaryServerInterceptor(grpcLogging.ServerPayloadLoggingDecider) grpc.UnaryServerInterceptor
	Info(...any)
	Debugf(string, ...any)
}

// API grpc API.
type API struct {
	options *Options
	logger  Logger
	server  *grpc.Server
	listen  net.Listener
}

// NewAPI возвращает обработчик grpc запросов.
func NewAPI(opt *Options, logger Logger, services ...Service) (*API, error) {
	if opt == nil || logger == nil {
		return nil, errors.New("invalid arguments")
	}

	listen, err := net.Listen("tcp", opt.Addr)
	if err != nil {
		return nil, err
	}

	if logger == nil {
		logger = &dummyLogger{}
	}

	// подменяем логер на logrus
	logger.ReplaceGRPCLogger()

	// включаем сохранение времени исполнения запроса в метрики
	grpcPrometheus.EnableHandlingTimeHistogram()

	interceptors := []grpc.UnaryServerInterceptor{
		// tracing
		otelgrpc.UnaryServerInterceptor(),
		// считаем метрики
		grpcPrometheus.UnaryServerInterceptor,
		// добавляем теги для телеметрии в context
		grpcCTXTags.UnaryServerInterceptor(
			grpcCTXTags.WithFieldExtractor(grpcCTXTags.CodeGenRequestFieldExtractor),
		),
		// логируем исполнение
		logger.UnaryServerInterceptor(),
		// логируем тело запроса
		logger.PayloadUnaryServerInterceptor(func(ctx context.Context, fullMethodName string, servingObject any) bool {
			if decider, ok := servingObject.(PayloadLoggingDecider); ok {
				return decider.PayloadLoggingDecider(fullMethodName)
			}

			// тут можно решать для какого именно запроса логировать тело
			return true
		}),
		// используем recovery
		grpcRecovery.UnaryServerInterceptor(),
	}

	// собираем сервисные интерцепторы
	if serviceInterceptors := collectServiceUnaryInterceptors(services); len(serviceInterceptors) > 0 {
		interceptors = append(interceptors, serviceInterceptors...)
	}

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptors...))

	if opt.EnableReflection {
		// используем рефлексию дабы можно было видеть все доступные методы и их свойства в postman
		reflection.Register(server)
	}

	// регистрируем сервисы
	registerServices(logger, server, services)

	// обнуляем метрики для всех сервисов
	grpcPrometheus.Register(server)

	a := API{
		options: opt,
		logger:  logger,
		server:  server,
		listen:  listen,
	}

	return &a, nil
}

// Serve запускает сервер.
func (a *API) Serve() error {
	a.logger.Info("Running GRPC Server")
	return a.server.Serve(a.listen)
}

// Close закрывает сервер.
func (a *API) Close() error {
	a.logger.Info("Shutting down GRPC Server")
	a.server.GracefulStop() // сам закрывает net.Listener
	return nil
}

// registerServices регистрирует сервисы для сервера.
func registerServices(logger Logger, server *grpc.Server, services []Service) {
	// регистрируем сервисы
	for _, service := range services {
		descriptor := service.Descriptor()

		if logger != nil {
			logger.Debugf("Register GRPC service %q", descriptor.ServiceName)
		}

		switch t := service.(type) {
		case *ServiceWithUnaryInterceptors:
			server.RegisterService(descriptor, t.Service)
		default:
			server.RegisterService(descriptor, service)
		}
	}
}
