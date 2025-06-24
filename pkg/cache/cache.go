package cache

import (
	"time"

	"aq3cms/pkg/logger"
)

// Cache 缓存接口
type Cache interface {
	// Get 获取缓存
	Get(key string) (interface{}, bool)

	// Set 设置缓存
	Set(key string, value interface{}, expire time.Duration) error

	// Delete 删除缓存
	Delete(key string)

	// Clear 清空缓存
	Clear()

	// Has 检查缓存是否存在
	Has(key string) bool

	// DeleteByPrefix 删除指定前缀的缓存
	DeleteByPrefix(prefix string)
}

// SafeSet 安全设置缓存，自动处理错误日志
func SafeSet(cache Cache, key string, value interface{}, expire time.Duration) {
	if err := cache.Set(key, value, expire); err != nil {
		logger.Error("缓存设置失败", "key", key, "error", err)
	}
}
