package grpc

type PayloadLoggingDecider interface {
	// PayloadLoggingDecider решает логгировать тела метода или нет.
	PayloadLoggingDecider(fullMethod string) bool
}
