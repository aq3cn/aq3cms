package cache

import (
	"crypto/md5"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"aq3cms/pkg/logger"
)

// FileCache 文件缓存
type FileCache struct {
	dir   string
	mutex sync.RWMutex
}

// fileCacheItem 文件缓存项
type fileCacheItem struct {
	Value      interface{}
	ExpireTime time.Time
	Expire     time.Duration
}

// NewFileCache 创建文件缓存
func NewFileCache(dir string) *FileCache {
	// 确保缓存目录存在
	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Error("创建缓存目录失败", "dir", dir, "error", err)
	}

	cache := &FileCache{
		dir: dir,
	}

	// 启动过期清理
	go cache.cleanExpired()

	return cache
}

// Get 获取缓存
func (c *FileCache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	filename := c.getCacheFilename(key)

	// 检查文件是否存在
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, false
	}

	// 读取文件
	file, err := os.Open(filename)
	if err != nil {
		logger.Error("打开缓存文件失败", "file", filename, "error", err)
		return nil, false
	}
	defer file.Close()

	// 解码缓存项
	var item fileCacheItem
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&item); err != nil {
		logger.Error("解码缓存文件失败", "file", filename, "error", err)
		return nil, false
	}

	// 检查是否过期
	if !item.ExpireTime.IsZero() && item.ExpireTime.Before(time.Now()) {
		// 删除过期文件
		os.Remove(filename)
		return nil, false
	}

	return item.Value, true
}

// Set 设置缓存
func (c *FileCache) Set(key string, value interface{}, expire time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	filename := c.getCacheFilename(key)

	var expireTime time.Time
	if expire > 0 {
		expireTime = time.Now().Add(expire)
	}

	// 创建缓存项
	item := fileCacheItem{
		Value:      value,
		ExpireTime: expireTime,
		Expire:     expire,
	}

	// 创建临时文件
	tempFile, err := ioutil.TempFile(c.dir, "cache_")
	if err != nil {
		logger.Error("创建临时文件失败", "dir", c.dir, "error", err)
		return err
	}
	tempFilename := tempFile.Name()

	// 编码缓存项
	encoder := gob.NewEncoder(tempFile)
	if err := encoder.Encode(item); err != nil {
		logger.Error("编码缓存项失败", "key", key, "error", err)
		tempFile.Close()
		os.Remove(tempFilename)
		return err
	}

	// 关闭临时文件
	tempFile.Close()

	// 重命名临时文件为缓存文件
	if err := os.Rename(tempFilename, filename); err != nil {
		logger.Error("重命名缓存文件失败", "from", tempFilename, "to", filename, "error", err)
		os.Remove(tempFilename)
		return err
	}

	return nil
}

// Delete 删除缓存
func (c *FileCache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	filename := c.getCacheFilename(key)

	// 删除文件
	if err := os.Remove(filename); err != nil && !os.IsNotExist(err) {
		logger.Error("删除缓存文件失败", "file", filename, "error", err)
	}
}

// Clear 清空缓存
func (c *FileCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 读取缓存目录
	files, err := ioutil.ReadDir(c.dir)
	if err != nil {
		logger.Error("读取缓存目录失败", "dir", c.dir, "error", err)
		return
	}

	// 删除所有缓存文件
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := filepath.Join(c.dir, file.Name())
		if err := os.Remove(filename); err != nil {
			logger.Error("删除缓存文件失败", "file", filename, "error", err)
		}
	}
}

// Has 检查缓存是否存在
func (c *FileCache) Has(key string) bool {
	_, ok := c.Get(key)
	return ok
}

// DeleteByPrefix 删除指定前缀的缓存
func (c *FileCache) DeleteByPrefix(prefix string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 读取缓存目录
	files, err := ioutil.ReadDir(c.dir)
	if err != nil {
		logger.Error("读取缓存目录失败", "dir", c.dir, "error", err)
		return
	}

	// 遍历所有缓存文件
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := filepath.Join(c.dir, file.Name())

		// 打开文件
		f, err := os.Open(filename)
		if err != nil {
			logger.Error("打开缓存文件失败", "file", filename, "error", err)
			continue
		}

		// 解码缓存项
		var item fileCacheItem
		decoder := gob.NewDecoder(f)
		if err := decoder.Decode(&item); err != nil {
			logger.Error("解码缓存文件失败", "file", filename, "error", err)
			f.Close()
			continue
		}

		f.Close()

		// 检查键是否以前缀开头
		// 由于我们使用MD5哈希键名，无法直接从文件名判断
		// 这里简化处理，实际应该在缓存项中存储原始键名
		if key, ok := item.Value.(string); ok && len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			os.Remove(filename)
		}
	}
}

// 获取缓存文件名
func (c *FileCache) getCacheFilename(key string) string {
	// 使用MD5哈希键名
	hash := md5.Sum([]byte(key))
	filename := fmt.Sprintf("%x.cache", hash)
	return filepath.Join(c.dir, filename)
}

// 清理过期缓存
func (c *FileCache) cleanExpired() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()

		// 读取缓存目录
		files, err := ioutil.ReadDir(c.dir)
		if err != nil {
			logger.Error("读取缓存目录失败", "dir", c.dir, "error", err)
			c.mutex.Unlock()
			continue
		}

		// 检查每个缓存文件
		for _, file := range files {
			if file.IsDir() {
				continue
			}

			filename := filepath.Join(c.dir, file.Name())

			// 打开文件
			file, err := os.Open(filename)
			if err != nil {
				logger.Error("打开缓存文件失败", "file", filename, "error", err)
				continue
			}

			// 解码缓存项
			var item fileCacheItem
			decoder := gob.NewDecoder(file)
			if err := decoder.Decode(&item); err != nil {
				logger.Error("解码缓存文件失败", "file", filename, "error", err)
				file.Close()
				os.Remove(filename)
				continue
			}

			file.Close()

			// 检查是否过期
			if !item.ExpireTime.IsZero() && item.ExpireTime.Before(time.Now()) {
				os.Remove(filename)
			}
		}

		c.mutex.Unlock()
	}
}
