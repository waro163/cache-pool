package cachepool

import "github.com/gomodule/redigo/redis"

// CacheDriver is cache pool driver interface
type CacheDriver interface {
	Connect(Configuration) error
	CloseConnection() error
	NativePool() *redis.Pool
	PingPong() (bool, error)
	Set(key string, value interface{}) error
	SetEX(key string, value interface{}, expirationtime int64) error
	Get(key string) (interface{}, error)
}
