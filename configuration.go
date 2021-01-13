package cachepool

import (
	"crypto/tls"
	"time"
)

const (
	// DefaultRedisNetwork the redis network option, "tcp".
	DefaultRedisNetwork = "tcp"
	// DefaultRedisAddr the redis address option, "127.0.0.1:6379".
	DefaultRedisAddr = "127.0.0.1:6379"
	// DefaultRedisTimeout the redis idle timeout option, time.Duration(30) * time.Second
	DefaultRedisTimeout = time.Duration(30) * time.Second
)

// Configuration the redis configuration
type Configuration struct {
	// Network protocol. Defaults to "tcp".
	Network string
	// Addr of a single redis server instance.
	// See "Clusters" field for clusters support.
	// Defaults to "127.0.0.1:6379".
	Addr string
	// Clusters a list of network addresses for clusters.
	// If not empty "Addr" is ignored.
	// Currently only Radix() Driver supports it.
	Clusters []string
	// Password string .If no password then no 'AUTH'. Defaults to "".
	Password string
	// If Database is empty "" then no 'SELECT'. Defaults to "".
	Database string

	// Timeout for connect, write and read, defaults to 30 seconds, 0 means no timeout.
	Timeout time.Duration

	// TLSConfig will cause Dial to perform a TLS handshake using the provided
	// config. If is nil then no TLS is used.
	// See https://golang.org/pkg/crypto/tls/#Config
	TLSConfig *tls.Config
}

type Configurator func(conf *Configuration)

// DefaultConfig is default cache pool configuration
func DefaultConfig() *Configuration {
	return &Configuration{
		Network:   DefaultRedisNetwork,
		Addr:      DefaultRedisAddr,
		Password:  "",
		Database:  "",
		Timeout:   DefaultRedisTimeout,
		TLSConfig: nil,
	}
}

func (config *Configuration) SetNetWork(network string) {
	config.Network = network
}

func (config *Configuration) SetAddr(addr string) {
	config.Addr = addr
}

func (config *Configuration) SetPassword(password string) {
	config.Password = password
}

func (config *Configuration) SetDatabase(db string) {
	config.Database = db
}

func (config *Configuration) SetTimeOut(timeout time.Duration) {
	config.Timeout = timeout
}

func (config *Configuration) Configure(configurators ...Configurator) {
	for _, cfg := range configurators {
		if cfg != nil {
			cfg(config)
		}
	}
}

func WithNetWork(network string) Configurator {
	return func(conf *Configuration) {
		conf.Network = network
	}
}

func WithAddr(addr string) Configurator {
	return func(conf *Configuration) {
		conf.Addr = addr
	}
}

func WithPassword(password string) Configurator {
	return func(conf *Configuration) {
		conf.Password = password
	}
}

func WithTimeOut(timeout time.Duration) Configurator {
	return func(conf *Configuration) {
		conf.Timeout = timeout
	}
}

func WithDatabase(db string) Configurator {
	return func(conf *Configuration) {
		conf.Database = db
	}
}
