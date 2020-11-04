package main

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
	. "github.com/waro163/cache-pool"
)

func main() {
	cache := DefaultCache()
	defer cache.Close()
	cache.SetPrefix("myapp")
	cache.SetDelim(":")
	cache.Set("name", "waro163")
	value, err := cache.Get("name")
	if err != nil {
		panic(err)
	}
	fmt.Println(redis.String(value, nil))
}
