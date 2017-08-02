package cache

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/kartiksura/wikinow/conf"
)

var (
	pool *redis.Pool
)

//GetString returns the value for the key in the redis
func GetString(key string) (string, error) {
	conn := pool.Get()
	defer conn.Close()

	data, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return "", fmt.Errorf("error getting key %s: %v", key, err)
	}
	return data, err
}

//SetString stores the value against the key in redis
func SetString(key string, value string) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		return fmt.Errorf("error setting key %s to %v: %v", key, value, err)
	}
	return err
}

//EnQueue inserts an object in the redis list
func EnQueue(key, value interface{}) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("LPUSH", key, value)
	if err != nil {
		return fmt.Errorf("error LPUSH key %s to %v: %v", key, value, err)
	}
	return err
}

//DeQueue retrieves the head of the list in redis
func DeQueue(key string) (interface{}, error) {
	conn := pool.Get()
	defer conn.Close()

	var data []byte
	data, err := redis.Bytes(conn.Do("LPOP", key))
	if err != nil {
		return nil, fmt.Errorf("error LPOP key %s: %v", key, err)
	}
	return data, err
}

func init() {
	pool = newPool(conf.WikiNowConfig.Redis.HostPort)
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
