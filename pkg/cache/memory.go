package cache

import (
	"sync"
	"time"
)

// MemoryCache 内存缓存
type MemoryCache struct {
	data  map[string]*cacheItem
	mutex sync.RWMutex
}

// cacheItem 缓存项
type cacheItem struct {
	value      interface{}
	expireTime time.Time
	expire     time.Duration
}

// NewMemoryCache 创建内存缓存
func NewMemoryCache() *MemoryCache {
	cache := &MemoryCache{
		data: make(map[string]*cacheItem),
	}

	// 启动过期清理
	go cache.cleanExpired()

	return cache
}

// Get 获取缓存
func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, ok := c.data[key]
	if !ok {
		return nil, false
	}

	// 检查是否过期
	if !item.expireTime.IsZero() && item.expireTime.Before(time.Now()) {
		return nil, false
	}

	return item.value, true
}

// Set 设置缓存
func (c *MemoryCache) Set(key string, value interface{}, expire time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var expireTime time.Time
	if expire > 0 {
		expireTime = time.Now().Add(expire)
	}

	c.data[key] = &cacheItem{
		value:      value,
		expireTime: expireTime,
		expire:     expire,
	}

	return nil
}

// Delete 删除缓存
func (c *MemoryCache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, key)
}

// Clear 清空缓存
func (c *MemoryCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = make(map[string]*cacheItem)
}

// Has 检查缓存是否存在
func (c *MemoryCache) Has(key string) bool {
	_, ok := c.Get(key)
	return ok
}

// DeleteByPrefix 删除指定前缀的缓存
func (c *MemoryCache) DeleteByPrefix(prefix string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for key := range c.data {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			delete(c.data, key)
		}
	}
}

// 清理过期缓存
func (c *MemoryCache) cleanExpired() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()

		for key, item := range c.data {
			if !item.expireTime.IsZero() && item.expireTime.Before(time.Now()) {
				delete(c.data, key)
			}
		}

		c.mutex.Unlock()
	}
}
