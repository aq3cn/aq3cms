package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"aq3cms/config"
	"aq3cms/pkg/logger"
	"github.com/gomodule/redigo/redis"
)

// RedisCache Redis缓存
type RedisCache struct {
	pool   *redis.Pool
	prefix string
}

// NewRedisCache 创建Redis缓存
func NewRedisCache(cfg config.CacheConfig) (*RedisCache, error) {
	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
			if err != nil {
				return nil, err
			}

			// 如果有密码，进行认证
			if cfg.Password != "" {
				if _, err := c.Do("AUTH", cfg.Password); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	// 测试连接
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("PING")
	if err != nil {
		return nil, fmt.Errorf("Redis连接失败: %v", err)
	}

	return &RedisCache{
		pool:   pool,
		prefix: "aq3cms:",
	}, nil
}

// Get 获取缓存
func (c *RedisCache) Get(key string) (interface{}, bool) {
	conn := c.pool.Get()
	defer conn.Close()

	// 获取缓存
	data, err := redis.Bytes(conn.Do("GET", c.prefix+key))
	if err != nil {
		if err != redis.ErrNil {
			logger.Error("Redis获取缓存失败", "key", key, "error", err)
		}
		return nil, false
	}

	// 解码数据
	var value interface{}
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&value); err != nil {
		logger.Error("解码缓存数据失败", "key", key, "error", err)
		return nil, false
	}

	return value, true
}

// Set 设置缓存
func (c *RedisCache) Set(key string, value interface{}, expire time.Duration) error {
	conn := c.pool.Get()
	defer conn.Close()

	// 编码数据
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(value); err != nil {
		logger.Error("编码缓存数据失败", "key", key, "error", err)
		return err
	}

	// 设置缓存
	if expire > 0 {
		_, err := conn.Do("SETEX", c.prefix+key, int64(expire.Seconds()), buffer.Bytes())
		if err != nil {
			logger.Error("Redis设置缓存失败", "key", key, "error", err)
			return err
		}
	} else {
		_, err := conn.Do("SET", c.prefix+key, buffer.Bytes())
		if err != nil {
			logger.Error("Redis设置缓存失败", "key", key, "error", err)
			return err
		}
	}

	return nil
}

// Delete 删除缓存
func (c *RedisCache) Delete(key string) {
	conn := c.pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", c.prefix+key)
	if err != nil {
		logger.Error("Redis删除缓存失败", "key", key, "error", err)
	}
}

// Clear 清空缓存
func (c *RedisCache) Clear() {
	conn := c.pool.Get()
	defer conn.Close()

	// 获取所有键
	keys, err := redis.Strings(conn.Do("KEYS", c.prefix+"*"))
	if err != nil {
		logger.Error("Redis获取所有键失败", "error", err)
		return
	}

	// 如果没有键，直接返回
	if len(keys) == 0 {
		return
	}

	// 删除所有键
	_, err = conn.Do("DEL", redis.Args{}.AddFlat(keys)...)
	if err != nil {
		logger.Error("Redis删除所有键失败", "error", err)
	}
}

// Has 检查缓存是否存在
func (c *RedisCache) Has(key string) bool {
	conn := c.pool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", c.prefix+key))
	if err != nil {
		logger.Error("Redis检查键是否存在失败", "key", key, "error", err)
		return false
	}

	return exists
}

// DeleteByPrefix 删除指定前缀的缓存
func (c *RedisCache) DeleteByPrefix(prefix string) {
	conn := c.pool.Get()
	defer conn.Close()

	// 获取所有匹配的键
	keys, err := redis.Strings(conn.Do("KEYS", c.prefix+prefix+"*"))
	if err != nil {
		logger.Error("Redis获取前缀键失败", "prefix", prefix, "error", err)
		return
	}

	// 如果没有键，直接返回
	if len(keys) == 0 {
		return
	}

	// 删除所有键
	_, err = conn.Do("DEL", redis.Args{}.AddFlat(keys)...)
	if err != nil {
		logger.Error("Redis删除前缀键失败", "prefix", prefix, "error", err)
	}
}
