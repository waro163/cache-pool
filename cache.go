package cachepool

const (
	// DefaultPrefix is cache key prefix, it is to distinguish which app the key belongs to.
	DefaultPrefix = ""
	// DefaultDelim is cache key delimitation option, "-".
	DefaultDelim = "-"
)

// Cache is cache pool client
type Cache struct {
	// redis configuration
	config *Configuration
	//
	driver CacheDriver
	// Prefix "myprefix-for-this-website". Defaults to "".
	Prefix string
	// Delim the delimeter for the keys on the sessiondb. Defaults to "-".
	Delim string
}

// DefaultCache default cache pool client
func DefaultCache() *Cache {
	cache := &Cache{
		config: DefaultConfig(),
		driver: NewDriver(),
		Prefix: DefaultPrefix,
		Delim:  DefaultDelim,
	}
	if err := cache.driver.Connect(*cache.config); err != nil {
		panic(err)
	}
	_, err := cache.driver.PingPong()
	if err != nil {
		panic(err)
	}
	return cache
}

// NewCache new a custom configuration cache
func NewCache(config *Configuration, driver CacheDriver) *Cache {

	cache := &Cache{}

	if config.Timeout < 0 {
		config.Timeout = DefaultRedisTimeout
	}

	if config.Network == "" {
		config.Network = DefaultRedisNetwork
	}

	if config.Addr == "" {
		config.Addr = DefaultRedisAddr
	}

	cache.config = config
	cache.driver = driver

	if err := cache.driver.Connect(*cache.config); err != nil {
		panic(err)
	}
	_, err := cache.driver.PingPong()
	if err != nil {
		panic(err)
	}
	return cache
}

// Close close cache pool
func (cache *Cache) Close() error {
	return cache.driver.CloseConnection()
}

// SetPrefix set the cache key prefix
func (cache *Cache) SetPrefix(prefix string) {
	cache.Prefix = prefix
}

// SetDelim set the cache key delimitation
func (cache *Cache) SetDelim(delim string) {
	cache.Delim = delim
}

func (cache *Cache) makeKey(key string) string {
	return cache.Prefix + cache.Delim + key
}

// Get retrive the key from cache storage
func (cache *Cache) Get(key string) (interface{}, error) {
	return cache.driver.Get(cache.makeKey(key))
}

// Set set the cache key
func (cache *Cache) Set(key string, value interface{}) error {
	return cache.driver.Set(cache.makeKey(key), value)
}
