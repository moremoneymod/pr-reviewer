package config

import (
	"errors"
	"net"
	"os"
	"time"
)

const (
	httpHostName        = "HTTP_HOST"
	httpPortName        = "HTTP_PORT"
	httpTimeoutName     = "HTTP_TIMEOUT_SECONDS"
	httpIdleTimeoutName = "HTTP_IDLE_TIMEOUT_SECONDS"
)

type HTTPConfig struct {
	host        string
	port        string
	timeout     time.Duration
	idleTimeout time.Duration
}

func NewHTTPConfig() (HTTPConfig, error) {
	host := os.Getenv(httpHostName)
	if len(host) == 0 {
		return HTTPConfig{}, errors.New("http host not found")
	}
	port := os.Getenv(httpPortName)
	if len(port) == 0 {
		return HTTPConfig{}, errors.New("http port not found")
	}

	timeout, err := time.ParseDuration(os.Getenv(httpTimeoutName))
	if err != nil {
		return HTTPConfig{}, err
	}

	idleTimeout, err := time.ParseDuration(os.Getenv(httpIdleTimeoutName))
	if err != nil {
		return HTTPConfig{}, err
	}

	return HTTPConfig{host, port, timeout, idleTimeout}, nil
}

func (cfg *HTTPConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

func (cfg *HTTPConfig) IDLETimeout() time.Duration {
	return cfg.idleTimeout
}

func (cfg *HTTPConfig) Timeout() time.Duration {
	return cfg.timeout
}
