package cachepool

const (
	// DefaultPrefix is cache key prefix, it is to distinguish which app the key belongs to.
	DefaultPrefix = ""
	// DefaultDelim is cache key delimitation option, "-".
	DefaultDelim = "-"
)

// Cache is cache pool client
type Cache struct {
	// redis driver's configuration
	cfg Config
	//
	Driver CacheDriver
	// Prefix "myprefix-for-this-website". Defaults to "".
	Prefix string
	// Delim the delimeter for the keys on the sessiondb. Defaults to "-".
	Delim string
}

// DefaultCache default cache pool client
func DefaultCache() *Cache {
	cache := &Cache{
		cfg:    DefaultConfig(),
		Driver: NewDriver(),
		Prefix: DefaultPrefix,
		Delim:  DefaultDelim,
	}
	if err := cache.Driver.Connect(cache.cfg); err != nil {
		panic(err)
	}
	_, err := cache.Driver.PingPong()
	if err != nil {
		panic(err)
	}
	return cache
}

// NewCache new a custom configuration cache
func NewCache(cfg ...Config) *Cache {
	cache := DefaultCache()
	if len(cfg) > 0 {
		c := cfg[0]

		if c.Timeout < 0 {
			c.Timeout = DefaultRedisTimeout
		}

		if c.Network == "" {
			c.Network = DefaultRedisNetwork
		}

		if c.Addr == "" {
			c.Addr = DefaultRedisAddr
		}

		if c.MaxActive == 0 {
			c.MaxActive = 10
		}

		cache.cfg = c
	}

	if err := cache.Driver.Connect(cache.cfg); err != nil {
		panic(err)
	}
	_, err := cache.Driver.PingPong()
	if err != nil {
		panic(err)
	}
	return cache
}

// Close close cache pool
func (cache *Cache) Close() error {
	return cache.Driver.CloseConnection()
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
	return cache.Driver.Get(cache.makeKey(key))
}

// Set set the cache key
func (cache *Cache) Set(key string, value interface{}) error {
	return cache.Driver.Set(cache.makeKey(key), value)
}
