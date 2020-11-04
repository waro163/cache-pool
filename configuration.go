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

// Config the redis configuration used inside sessions
type Config struct {
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
	// MaxActive. Defaults to 10.
	MaxActive int
	// Timeout for connect, write and read, defaults to 30 seconds, 0 means no timeout.
	Timeout time.Duration

	// TLSConfig will cause Dial to perform a TLS handshake using the provided
	// config. If is nil then no TLS is used.
	// See https://golang.org/pkg/crypto/tls/#Config
	TLSConfig *tls.Config
}

// DefaultConfig is default cache pool configuration
func DefaultConfig() Config {
	return Config{
		Network:   DefaultRedisNetwork,
		Addr:      DefaultRedisAddr,
		Password:  "",
		Database:  "",
		MaxActive: 10,
		Timeout:   DefaultRedisTimeout,
		TLSConfig: nil,
	}
}