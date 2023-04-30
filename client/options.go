package grpc

// ConnectionOptions - настройки подключения.
type ConnectionOptions struct {
	Addr           string
	Insecure       bool
	TLSSkipVerify  bool
	LogBody        bool
	LogBodyMethods []string
}
