package gocache

import (
	"encoding/json"
	"github.com/goecology/muses/pkg/cache/mixcache/standard"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"time"
)

type Client struct {
	c *cache.Cache
}

// todo 过期
func NewMixCache() (client standard.MixCache, err error) {
	c := cache.New(time.Duration(0), time.Duration(0))
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

	c.c.Set(key, value, time.Duration(expire))
	return
}

func (c *Client) Get(key string) (resp interface{}, err error) {
	output, flag := c.c.Get(key)
	if !flag {
		err = errors.New("not exist")
		return
	}
	resp = output
	return
}
