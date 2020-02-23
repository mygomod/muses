package leveldb

import (
	"encoding/json"
	"github.com/i2eco/muses/pkg/cache/mixcache/standard"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
)

type Client struct {
	c *leveldb.DB
}

// todo defer
func NewMixCache(path string) (client standard.MixCache, err error) {
	db, err := leveldb.OpenFile(path, nil)
	//defer db.Close()
	client = &Client{
		c: db,
	}
	return
}

// todo 不支持数字
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

	switch putStr := value.(type) {
	case []byte:
		err = c.c.Put([]byte(key), putStr, nil)
		return
	case string:
		err = c.c.Put([]byte(key), []byte(putStr), nil)
		return
	}
	err = errors.New("type not exist")
	return
}

func (c *Client) Get(key string) (resp interface{}, err error) {
	var data []byte
	data, err = c.c.Get([]byte(key), nil)
	resp = data
	return
}
