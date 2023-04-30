package logrus

import (
	grpcLogging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	grpcLogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Logger struct {
	logger *logrus.Entry
}

func NewLogger(logger *logrus.Entry) *Logger {
	return &Logger{logger: logger}
}

func (l *Logger) UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return grpcLogrus.UnaryClientInterceptor(l.logger)
}

func (l *Logger) PayloadUnaryClientInterceptor(decider grpcLogging.ClientPayloadLoggingDecider) grpc.UnaryClientInterceptor {
	return grpcLogrus.PayloadUnaryClientInterceptor(l.logger, decider)
}

func (l *Logger) ReplaceGRPCLogger() {
	grpcLogrus.ReplaceGrpcLogger(l.logger)
}

func (l *Logger) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return grpcLogrus.UnaryServerInterceptor(l.logger)
}

func (l *Logger) PayloadUnaryServerInterceptor(decider grpcLogging.ServerPayloadLoggingDecider) grpc.UnaryServerInterceptor {
	return grpcLogrus.PayloadUnaryServerInterceptor(l.logger, decider)
}

func (l *Logger) Info(args ...any) {
	l.logger.Info(args...)
}

func (l *Logger) Debugf(format string, args ...any) {
	l.logger.Debugf(format, args...)
}
