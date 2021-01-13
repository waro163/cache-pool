package cachepool

import (
	"errors"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	// ErrRedisClosed an error with message 'redis: already closed'
	ErrRedisClosed = errors.New("redis: already closed")
	// ErrKeyNotFound a type of error of non-existing redis keys.
	// The producers(the library) of this error will dynamically wrap this error(fmt.Errorf) with the key name.
	// Usage:
	// if err != nil && errors.Is(err, ErrKeyNotFound) {
	// [...]
	// }
	ErrKeyNotFound = errors.New("key not found")
)

// RedigoDriver implement CacheDriver interface, it is the redigo Redis go client,
// contains the config and the redis pool
type RedigoDriver struct {

	// Maximum number of idle connections in the pool.
	MaxIdle int
	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	MaxActive int

	// Close connections after remaining idle for this duration. If the value
	// is zero, then idle connections are not closed. Applications should set
	// the timeout to a value less than the server's timeout.
	IdleTimeout time.Duration

	// If Wait is true and the pool is at the MaxActive limit, then Get() waits
	// for a connection to be returned to the pool before returning.
	Wait bool

	//connected is true when the Service has already connected
	connected bool

	// Close connections older than this duration. If the value is zero, then
	// the pool does not close connections based on age.
	MaxConnLifetime time.Duration

	pool *redis.Pool
}

// NewDriver is for new redigo driver
func NewDriver(r ...RedigoDriver) *RedigoDriver {
	red := &RedigoDriver{}
	if len(r) > 0 {
		red = &r[0]
	}
	return red
}

// Connect connects to the redis by redigo driver pool, called only once.
func (r *RedigoDriver) Connect(c Configuration) error {
	pool := &redis.Pool{IdleTimeout: r.IdleTimeout, MaxIdle: r.MaxIdle, Wait: r.Wait, MaxConnLifetime: r.MaxConnLifetime, MaxActive: r.MaxActive}
	pool.TestOnBorrow = func(conn redis.Conn, t time.Time) error {
		_, err := conn.Do("PING")
		return err
	}
	var options []redis.DialOption

	if c.Timeout > 0 {
		options = append(options,
			redis.DialConnectTimeout(c.Timeout),
			redis.DialReadTimeout(c.Timeout),
			redis.DialWriteTimeout(c.Timeout))
	}

	if c.TLSConfig != nil {
		options = append(options,
			redis.DialTLSConfig(c.TLSConfig),
			redis.DialUseTLS(true),
		)
	}

	pool.Dial = func() (redis.Conn, error) {
		conn, err := redis.Dial(c.Network, c.Addr, options...)
		if err != nil {
			return nil, err
		}

		if c.Password != "" {
			if _, err = conn.Do("AUTH", c.Password); err != nil {
				conn.Close()
				return nil, err
			}
		}

		if c.Database != "" {
			if _, err = conn.Do("SELECT", c.Database); err != nil {
				conn.Close()
				return nil, err
			}
		}

		return conn, err
	}

	r.connected = true
	r.pool = pool
	return nil
}

// CloseConnection close driver pool
func (r *RedigoDriver) CloseConnection() error {
	if r.pool != nil {
		return r.pool.Close()
	}
	return ErrRedisClosed
}

// NativePool is to get native redis pool
func (r *RedigoDriver) NativePool() *redis.Pool {
	return r.pool
}

// PingPong sends a ping and receives a pong, if no pong received then returns false and filled error
func (r *RedigoDriver) PingPong() (bool, error) {
	c := r.pool.Get()
	defer c.Close()
	msg, err := c.Do("PING")
	if err != nil || msg == nil {
		return false, err
	}
	return (msg == "PONG"), nil
}

// Set sets a key-value to the redis store.
func (r *RedigoDriver) Set(key string, value interface{}) (err error) {
	c := r.pool.Get()
	defer c.Close()
	if c.Err() != nil {
		return c.Err()
	}

	_, err = c.Do("SET", key, value)

	return
}

// SetEX sets a key-value to the redis store with expiration time.
func (r *RedigoDriver) SetEX(key string, value interface{}, expirationtime int64) (err error) {
	c := r.pool.Get()
	defer c.Close()
	if c.Err() != nil {
		return c.Err()
	}
	if expirationtime > 0 {
		_, err = c.Do("SETEX", key, expirationtime, value)
	} else {
		_, err = c.Do("SET", key, value)
	}
	return
}

// Get returns value, err by its key
// returns nil and a filled error if something bad happened.
func (r *RedigoDriver) Get(key string) (interface{}, error) {
	c := r.pool.Get()
	defer c.Close()
	if err := c.Err(); err != nil {
		return nil, err
	}

	redisVal, err := c.Do("GET", key)
	if err != nil {
		return nil, err
	}
	if redisVal == nil {
		return nil, fmt.Errorf("%s: %w", key, ErrKeyNotFound)
	}
	return redisVal, nil
}
