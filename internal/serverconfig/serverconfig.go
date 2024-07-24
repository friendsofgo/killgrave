package serverconfig

import (
	"io"

	"github.com/gorilla/handlers"
)

type ServerConfig struct {
	CORSOptions []handlers.CORSOption
	LogLevel    int
	LogBodyMax  int
	LogWriter   io.Writer
}

type ServerOption func(*ServerConfig)

func WithCORSOptions(corsOptions []handlers.CORSOption) ServerOption {
	return func(cfg *ServerConfig) {
		cfg.CORSOptions = corsOptions
	}
}

func WithLogLevel(logLevel int) ServerOption {
	return func(cfg *ServerConfig) {
		cfg.LogLevel = logLevel
	}
}

func WithLogBodyMax(logBodyMax int) ServerOption {
	return func(cfg *ServerConfig) {
		cfg.LogBodyMax = logBodyMax
	}
}

func WithLogWriter(logWriter io.Writer) ServerOption {
	return func(cfg *ServerConfig) {
		cfg.LogWriter = logWriter
	}
}
