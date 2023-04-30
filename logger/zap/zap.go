package zap

import (
	"fmt"

	grpcLogging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	grpcZap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Logger struct {
	logger *zap.Logger
}

func NewLogger(logger *zap.Logger) *Logger {
	return &Logger{logger: logger}
}

func (l *Logger) UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return grpcZap.UnaryClientInterceptor(l.logger)
}

func (l *Logger) PayloadUnaryClientInterceptor(decider grpcLogging.ClientPayloadLoggingDecider) grpc.UnaryClientInterceptor {
	return grpcZap.PayloadUnaryClientInterceptor(l.logger, decider)
}

func (l *Logger) ReplaceGRPCLogger() {
	grpcZap.ReplaceGrpcLoggerV2(l.logger)
}

func (l *Logger) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return grpcZap.UnaryServerInterceptor(l.logger)
}

func (l *Logger) PayloadUnaryServerInterceptor(decider grpcLogging.ServerPayloadLoggingDecider) grpc.UnaryServerInterceptor {
	return grpcZap.PayloadUnaryServerInterceptor(l.logger, decider)
}

func (l *Logger) Info(args ...any) {
	l.logger.Info(fmt.Sprint(args...))
}

func (l *Logger) Debugf(format string, args ...any) {
	args = append([]any{format}, args...)
	l.logger.Debug(fmt.Sprint(args...))
}
