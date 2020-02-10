package redis

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func (r *Client) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	conn := r.pool.Get()
	defer conn.Close()
	return conn.Do(commandName, args...)
}

func (r *Client) Send(commandName string, args ...interface{}) error {
	conn := r.pool.Get()
	defer conn.Close()
	return conn.Send(commandName, args...)
}

func (r *Client) Flush() error {
	conn := r.pool.Get()
	defer conn.Close()
	return conn.Flush()
}

func (r *Client) GetString(key string) (string, error) {
	return redis.String(r.Do("GET", key))
}

func (r *Client) GetInt(key string) (int, error) {
	return redis.Int(r.Do("GET", key))
}

func (r *Client) GetInt64(key string) (int64, error) {
	return redis.Int64(r.Do("GET", key))
}

func (r *Client) GetBool(key string) (bool, error) {
	return redis.Bool(r.Do("GET", key))
}

func (r *Client) GetObject(key string, val interface{}) error {
	reply, err := r.GetString(key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(reply), val)
}

func (r *Client) Get(key string) (interface{}, error) {
	return r.Do("GET", key)
}

// Set 用法：Set("key", val, 60)，其中 expire 的单位为秒
func (r *Client) Set(key string, val interface{}, expire int) (interface{}, error) {
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
		return r.Do("SETEX", key, expire, value)
	} else {
		return r.Do("SET", key, value)
	}
}

// Exists 检查键是否存在
func (r *Client) Exists(key string) (bool, error) {
	return redis.Bool(r.Do("EXISTS", key))
}

//Del 删除键
func (r *Client) Del(key string) error {
	_, err := r.Do("DEL", key)
	return err
}

// TTL 以秒为单位。当 key 不存在时，返回 -2 。 当 key 存在但没有设置剩余生存时间时，返回 -1
func (r *Client) Ttl(key string) (ttl int64, err error) {
	return redis.Int64(r.Do("TTL", key))
}

// Expire 设置键过期时间，expire的单位为秒
func (r *Client) Expire(key string, expire int) error {
	_, err := redis.Bool(r.Do("EXPIRE", key, expire))
	return err
}

func (r *Client) Incr(key string) (val int64, err error) {
	return redis.Int64(r.Do("INCR", key))
}

func (r *Client) IncrBy(key string, amount int) (val int64, err error) {
	return redis.Int64(r.Do("INCRBY", key, amount))
}

func (r *Client) Decr(key string) (val int64, err error) {
	return redis.Int64(r.Do("DECR", key))
}

func (r *Client) DecrBy(key string, amount int) (val int64, err error) {
	return redis.Int64(r.Do("DECRBY", key, amount))
}

// Hmset 用法：cache.Myredis.Hmset("key", val, 60)
func (r *Client) Hmset(key string, val interface{}, expire int) (err error) {
	conn := r.pool.Get()
	defer conn.Close()
	err = conn.Send("HMSET", redis.Args{}.Add(key).AddFlat(val)...)
	if err != nil {
		return
	}
	if expire > 0 {
		err = conn.Send("EXPIRE", key, int64(expire))
	}
	if err != nil {
		return
	}
	err = conn.Flush()
	return
}

func (r *Client) Hset(key, field string, val interface{}) (interface{}, error) {
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
	return r.Do("HSET", key, field, value)
}

// Hmget 用法：cache.Myredis.Hget("key", "field_name")
func (r *Client) Hget(key, field string) (reply interface{}, err error) {
	reply, err = r.Do("HGET", key, field)
	return
}

func (r *Client) HgetString(key, field string) (reply string, err error) {
	reply, err = redis.String(r.Do("HGET", key, field))
	return
}

func (r *Client) HgetInt(key, field string) (reply int, err error) {
	reply, err = redis.Int(r.Do("HGET", key, field))
	return
}

func (r *Client) HgetInt64(key, field string) (reply int64, err error) {
	reply, err = redis.Int64(r.Do("HGET", key, field))
	return
}

func (r *Client) HgetBool(key, field string) (reply bool, err error) {
	reply, err = redis.Bool(r.Do("HGET", key, field))
	return
}

func (r *Client) HgetObject(key, field string, val interface{}) error {
	reply, err := r.HgetString(key, field)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(reply), val)
}

// Hmget 用法：cache.Myredis.HgetAll("key", &val)
func (r *Client) HgetAll(key string, val interface{}) error {
	v, err := redis.Values(r.Do("HGETALL", key))
	if err != nil {
		return err
	}

	if err := redis.ScanStruct(v, val); err != nil {
		fmt.Println(err)
	}
	//fmt.Printf("%+v\n", val)
	return err
}

// Zadd 将一个成员元素及其分数值加入到有序集当中
func (r *Client) Zadd(key string, score int64, member string) (reply interface{}, err error) {
	return r.Do("ZADD", key, score, member)
}

// Zrem 移除有序集中的一个，不存在的成员将被忽略。
func (r *Client) Zrem(key string, member string) (reply interface{}, err error) {
	return r.Do("ZREM", key, member)
}

// Zscore 返回有序集中，成员的分数值。 如果成员元素不是有序集 key 的成员，或 key 不存在，返回 nil
func (r *Client) Zscore(key string, member string) (int64, error) {
	return redis.Int64(r.Do("ZSCORE", key, member))
}

// Zrank 返回有序集中指定成员的排名。其中有序集成员按分数值递增(从小到大)顺序排列。score 值最小的成员排名为 0
func (r *Client) Zrank(key, member string) (int64, error) {
	return redis.Int64(r.Do("ZRANK", key, member))
}

// Zrevrank 返回有序集中成员的排名。其中有序集成员按分数值递减(从大到小)排序。分数值最大的成员排名为 0 。
func (r *Client) Zrevrank(key, member string) (int64, error) {
	return redis.Int64(r.Do("ZREVRANK", key, member))
}

// Zrange 返回有序集中，指定区间内的成员。其中成员的位置按分数值递增(从小到大)来排序。具有相同分数值的成员按字典序(lexicographical order )来排列。
// 以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。或 以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。
func (r *Client) Zrange(key string, from, to int64) (map[string]int64, error) {
	return redis.Int64Map(r.Do("ZRANGE", key, from, to, "WITHSCORES"))
}

// Zrevrange 返回有序集中，指定区间内的成员。其中成员的位置按分数值递减(从大到小)来排列。具有相同分数值的成员按字典序(lexicographical order )来排列。
// 以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。或 以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。
func (r *Client) Zrevrange(key string, from, to int64) (map[string]int64, error) {
	return redis.Int64Map(r.Do("ZREVRANGE", key, from, to, "WITHSCORES"))
}

// ZrangeByScore 返回有序集合中指定分数区间的成员列表。有序集成员按分数值递增(从小到大)次序排列。
// 具有相同分数值的成员按字典序来排列
func (r *Client) ZrangeByScore(key string, from, to, offset int64, count int) (map[string]int64, error) {
	return redis.Int64Map(r.Do("ZRANGEBYSCORE", key, from, to, "WITHSCORES", "LIMIT", offset, count))
}

// ZrevrangeByScore 返回有序集中指定分数区间内的所有的成员。有序集成员按分数值递减(从大到小)的次序排列。
// 具有相同分数值的成员按字典序来排列
func (r *Client) ZrevrangeByScore(key string, from, to, offset int64, count int) (map[string]int64, error) {
	return redis.Int64Map(r.Do("ZREVRANGEBYSCORE", key, from, to, "WITHSCORES", "LIMIT", offset, count))
}

// Publish 将信息发送到指定的频道，返回接收到信息的订阅者数量
func (r *Client) Publish(channel, message string) (int, error) {
	return redis.Int(r.Do("PUBLISH", channel, message))
}

func (r *Client) Blpop(channel string, time interface{}) (interface{}, error) {
	return r.Do("BLPOP", channel, time)
}

func (r *Client) Rpush(channel, msg string) (interface{}, error) {
	return r.Do("RPUSH", channel, msg)
}
