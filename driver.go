package cachepool

// CacheDriver is cache pool driver interface
type CacheDriver interface {
	Connect(Config) error
	CloseConnection() error
	PingPong() (bool, error)
	Set(key string, value interface{}) error
	SetEX(key string, value interface{}, expirationtime int64) error
	Get(key string) (interface{}, error)
}
