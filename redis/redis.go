package redis

import (
	// log "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
	"github.com/jiangzhx/receive/bloom"
	"time"
)

var (
	redisServer   string             = "x01:6379"
	redisPassword string             = ""
	pool          *redis.Pool        = newPool(redisServer)
	filter        *bloom.BloomFilter = bloom.NewWithEstimates(100000, 0.01)
)

func newPool(address string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		MaxActive:   0,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
func GetConn() redis.Conn {
	return pool.Get()
}

func Badd(key string, value string) {
	hashes := filter.Hashes(value)
	conn := pool.Get()
	defer conn.Close()
	for i := 0; i < len(hashes); i++ {
		conn.Do("setbit", key, hashes[i], true)
	}
}

func Bexist(key string, value string) bool {
	hashes := filter.Hashes(value)
	conn := pool.Get()
	defer conn.Close()
	for i := 0; i < len(hashes); i++ {
		exist, _ := conn.Do("getbit", key, hashes[i])
		if exist == false {
			return false
		}
	}
	return true
}
