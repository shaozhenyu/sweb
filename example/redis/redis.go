package main

import (
	"fmt"

	"libs/cache"
)

func main() {
	RedisHost := "127.0.0.1:6379"
	rd_cache, err := cache.New(RedisHost, 0, 100)
	if err != nil {
		return
	}

	if ok := Set("a", "szy", rd_cache); !ok {
		fmt.Println("set k-v error")
		return
	}

	var v string
	var ok bool
	if v, ok = Get("a", rd_cache); !ok {
		fmt.Println("get k-v error")
		return
	}

	fmt.Println("value: ", v)
}

func Get(key string, rds cache.CacheStorage) (string, bool) {
	bs, err := rds.Get(key)
	if err != nil {
		return "", false
	}
	value := string(bs)
	return value, true
}

func Set(key, value string, rds cache.CacheStorage) bool {
	if err := rds.Set(key, []byte(value)); err != nil {
		return false
	}
	return true
}
