package cache

import (
	"fmt"
)

var (
	errNotFound = fmt.Errorf("cache not found")
)

type CacheStorage interface {
	Set(key string, val []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
	Incr(key string) (int64, error)
	Expire(key string, expire int64) error
	FlushDB() error
}

func CacheNotFound() error {
	return errNotFound
}

func IsCacheNotFound(err error) bool {
	return err == errNotFound
}
