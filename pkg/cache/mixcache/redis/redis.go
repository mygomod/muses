package redis

import (
	"encoding/json"
	"github.com/goecology/muses/pkg/cache/mixcache/standard"
	"github.com/gomodule/redigo/redis"
	"log"
	"os"
	"time"
)

type Client struct {
	c *redis.Pool
}

// todo 过期
func NewMixCache(addr string, debug bool, maxActive int, idleTimeout time.Duration, wait bool, dialOptions ...redis.DialOption) (client standard.MixCache, err error) {
	c := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr, dialOptions...)
			if err != nil {
				return nil, err
			}

			if debug {
				return redis.NewLoggingConn(c, log.New(os.Stderr, "", log.LstdFlags), "redis"), nil
			}
			return c, nil
		},
		// Use the TestOnBorrow function to check the health of an idle connection
		// before the connection is returned to the application. This example PINGs
		// connections that have been idle more than a minute:
		//
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:     maxActive,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout,
		Wait:        wait, // wait until getting connection from pool
	}
	client = &Client{
		c: c,
	}
	return
}

func (c *Client) Set(key string, val interface{}, expire int) (resp interface{}, err error) {
	var value interface{}
	switch v := val.(type) {
	case string, int, uint, int8, int16, int32, int64, float32, float64, bool:
		value = v
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		value = string(b)
	}
	if expire > 0 {
		return c.Do("SETEX", key, expire, value)
	} else {
		return c.Do("SET", key, value)
	}
}

func (c *Client) Get(key string) (resp interface{}, err error) {
	return c.Do("GET", key)
}

func (c *Client) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	conn := c.c.Get()
	defer conn.Close()
	return conn.Do(commandName, args...)
}
