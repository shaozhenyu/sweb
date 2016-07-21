package cache

import (
	"time"

	"gopkg.in/redis.v2"
)

var Cache CacheStorage = nil

type CacheRedis struct {
	conn *redis.Client
}

func New(addr string, db, maxPoolSize int) (CacheStorage, error) {
	opt := redis.Options{}
	opt.Addr = addr
	opt.DB = int64(db)
	opt.PoolSize = maxPoolSize

	conn := redis.NewTCPClient(&opt)
	err := conn.Ping().Err()
	Cache = &CacheRedis{conn}
	return Cache, err
}

func (this *CacheRedis) Set(key string, val []byte) error {
	return this.conn.Set(key, string(val)).Err()
}

func (this *CacheRedis) Get(key string) ([]byte, error) {
	s, err := this.conn.Get(key).Result()
	if err == redis.Nil {
		err = CacheNotFound()
	}
	return []byte(s), err
}

func (this *CacheRedis) Delete(key string) error {
	return this.conn.Del(key).Err()
}

func (this *CacheRedis) Incr(key string) (int64, error) {
	return this.conn.Incr(key).Result()
}

func (this *CacheRedis) Expire(key string, expire int64) error {
	return this.conn.Expire(key, time.Duration(expire)*time.Second).Err()
}

func (this *CacheRedis) FlushDB() error {
	return this.conn.FlushDb().Err()
}
