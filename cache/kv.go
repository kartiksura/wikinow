package cache

import (
	"fmt"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	Pool *redis.Pool
)

func init() {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = ":6379"
	}
	Pool = newPool(redisHost)
}

func newPool(server string) *redis.Pool {

	return &redis.Pool{

		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
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

func Get(key string) (string, error) {

	conn := Pool.Get()
	defer conn.Close()

	var data []byte
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return string(data), fmt.Errorf("error getting key %s: %v", key, err)
	}
	return string(data), err
}

func Set(key string, value string) error {

	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
	}
	return err
}

func EnQueue(key, value string) error {
	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("LPUSH", key, value)
	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("error LPUSH key %s to %s: %v", key, v, err)
	}
	return err
}

func DeQueue(key string) (string, error) {
	conn := Pool.Get()
	defer conn.Close()

	var data []byte
	data, err := redis.Bytes(conn.Do("LPOP", key))
	if err != nil {
		return string(data), fmt.Errorf("error LPOP key %s: %v", key, err)
	}
	return string(data), err
}
