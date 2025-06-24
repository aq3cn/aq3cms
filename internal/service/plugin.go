package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"aq3cms/config"
	"aq3cms/internal/model"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// PluginService 插件服务
type PluginService struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	pluginModel     *model.PluginModel
	pluginHookModel *model.PluginHookModel
	pluginDir       string
}

// NewPluginService 创建插件服务
func NewPluginService(db *database.DB, cache cache.Cache, config *config.Config) *PluginService {
	pluginDir := config.Plugin.Dir
	if pluginDir == "" {
		pluginDir = "plugins"
	}
	return &PluginService{
		db:              db,
		cache:           cache,
		config:          config,
		pluginModel:     model.NewPluginModel(db),
		pluginHookModel: model.NewPluginHookModel(db),
		pluginDir:       pluginDir,
	}
}

// GetPlugin 获取插件
func (s *PluginService) GetPlugin(id int64) (*model.Plugin, error) {
	// 从缓存获取
	cacheKey := fmt.Sprintf("plugin:%d", id)
	if cached, ok := s.cache.Get(cacheKey); ok {
		if plugin, ok := cached.(*model.Plugin); ok {
			return plugin, nil
		}
	}

	// 获取插件
	plugin, err := s.pluginModel.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 缓存插件
	cache.SafeSet(s.cache, cacheKey, plugin, time.Hour)

	return plugin, nil
}

// GetPluginByCode 根据代码获取插件
func (s *PluginService) GetPluginByCode(code string) (*model.Plugin, error) {
	// 从缓存获取
	cacheKey := fmt.Sprintf("plugin:code:%s", code)
	if cached, ok := s.cache.Get(cacheKey); ok {
		if plugin, ok := cached.(*model.Plugin); ok {
			return plugin, nil
		}
	}

	// 获取插件
	plugin, err := s.pluginModel.GetByCode(code)
	if err != nil {
		return nil, err
	}

	// 缓存插件
	cache.SafeSet(s.cache, cacheKey, plugin, time.Hour)

	return plugin, nil
}

// GetAllPlugins 获取所有插件
func (s *PluginService) GetAllPlugins(status int) ([]*model.Plugin, error) {
	// 从缓存获取
	cacheKey := fmt.Sprintf("plugin:all:%d", status)
	if cached, ok := s.cache.Get(cacheKey); ok {
		if plugins, ok := cached.([]*model.Plugin); ok {
			return plugins, nil
		}
	}

	// 获取所有插件
	plugins, err := s.pluginModel.GetAll(status)
	if err != nil {
		return nil, err
	}

	// 缓存插件
	cache.SafeSet(s.cache, cacheKey, plugins, time.Hour)

	return plugins, nil
}

// CreatePlugin 创建插件
func (s *PluginService) CreatePlugin(plugin *model.Plugin) (int64, error) {
	// 创建插件
	id, err := s.pluginModel.Create(plugin)
	if err != nil {
		return 0, err
	}

	// 清除缓存
	s.ClearPluginCache()

	return id, nil
}

// UpdatePlugin 更新插件
func (s *PluginService) UpdatePlugin(plugin *model.Plugin) error {
	// 更新插件
	err := s.pluginModel.Update(plugin)
	if err != nil {
		return err
	}

	// 清除缓存
	s.ClearPluginCache()

	return nil
}

// DeletePlugin 删除插件
func (s *PluginService) DeletePlugin(id int64) error {
	// 删除插件
	err := s.pluginModel.Delete(id)
	if err != nil {
		return err
	}

	// 清除缓存
	s.ClearPluginCache()

	return nil
}

// UpdatePluginStatus 更新插件状态
func (s *PluginService) UpdatePluginStatus(id int64, status int) error {
	// 更新插件状态
	err := s.pluginModel.UpdateStatus(id, status)
	if err != nil {
		return err
	}

	// 清除缓存
	s.ClearPluginCache()

	return nil
}

// GetPluginHook 获取插件钩子
func (s *PluginService) GetPluginHook(id int64) (*model.PluginHook, error) {
	// 从缓存获取
	cacheKey := fmt.Sprintf("plugin:hook:%d", id)
	if cached, ok := s.cache.Get(cacheKey); ok {
		if hook, ok := cached.(*model.PluginHook); ok {
			return hook, nil
		}
	}

	// 获取插件钩子
	hook, err := s.pluginHookModel.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 缓存插件钩子
	cache.SafeSet(s.cache, cacheKey, hook, time.Hour)

	return hook, nil
}

// GetPluginHookByCode 根据代码获取插件钩子
func (s *PluginService) GetPluginHookByCode(code string) (*model.PluginHook, error) {
	// 从缓存获取
	cacheKey := fmt.Sprintf("plugin:hook:code:%s", code)
	if cached, ok := s.cache.Get(cacheKey); ok {
		if hook, ok := cached.(*model.PluginHook); ok {
			return hook, nil
		}
	}

	// 获取插件钩子
	hook, err := s.pluginHookModel.GetByCode(code)
	if err != nil {
		return nil, err
	}

	// 缓存插件钩子
	cache.SafeSet(s.cache, cacheKey, hook, time.Hour)

	return hook, nil
}

// GetPluginHooksByPluginID 根据插件ID获取插件钩子
func (s *PluginService) GetPluginHooksByPluginID(pluginID int64) ([]*model.PluginHook, error) {
	// 从缓存获取
	cacheKey := fmt.Sprintf("plugin:hook:pluginid:%d", pluginID)
	if cached, ok := s.cache.Get(cacheKey); ok {
		if hooks, ok := cached.([]*model.PluginHook); ok {
			return hooks, nil
		}
	}

	// 获取插件钩子
	hooks, err := s.pluginHookModel.GetByPluginID(pluginID)
	if err != nil {
		return nil, err
	}

	// 缓存插件钩子
	cache.SafeSet(s.cache, cacheKey, hooks, time.Hour)

	return hooks, nil
}

// GetPluginHooksByPosition 根据位置获取插件钩子
func (s *PluginService) GetPluginHooksByPosition(position string) ([]*model.PluginHook, error) {
	// 从缓存获取
	cacheKey := fmt.Sprintf("plugin:hook:position:%s", position)
	if cached, ok := s.cache.Get(cacheKey); ok {
		if hooks, ok := cached.([]*model.PluginHook); ok {
			return hooks, nil
		}
	}

	// 获取插件钩子
	hooks, err := s.pluginHookModel.GetByPosition(position)
	if err != nil {
		return nil, err
	}

	// 缓存插件钩子
	cache.SafeSet(s.cache, cacheKey, hooks, time.Hour)

	return hooks, nil
}

// GetAllPluginHooks 获取所有插件钩子
func (s *PluginService) GetAllPluginHooks(status int) ([]*model.PluginHook, error) {
	// 从缓存获取
	cacheKey := fmt.Sprintf("plugin:hook:all:%d", status)
	if cached, ok := s.cache.Get(cacheKey); ok {
		if hooks, ok := cached.([]*model.PluginHook); ok {
			return hooks, nil
		}
	}

	// 获取所有插件钩子
	hooks, err := s.pluginHookModel.GetAll(status)
	if err != nil {
		return nil, err
	}

	// 缓存插件钩子
	cache.SafeSet(s.cache, cacheKey, hooks, time.Hour)

	return hooks, nil
}

// CreatePluginHook 创建插件钩子
func (s *PluginService) CreatePluginHook(hook *model.PluginHook) (int64, error) {
	// 创建插件钩子
	id, err := s.pluginHookModel.Create(hook)
	if err != nil {
		return 0, err
	}

	// 清除缓存
	s.ClearPluginHookCache()

	return id, nil
}

// UpdatePluginHook 更新插件钩子
func (s *PluginService) UpdatePluginHook(hook *model.PluginHook) error {
	// 更新插件钩子
	err := s.pluginHookModel.Update(hook)
	if err != nil {
		return err
	}

	// 清除缓存
	s.ClearPluginHookCache()

	return nil
}

// DeletePluginHook 删除插件钩子
func (s *PluginService) DeletePluginHook(id int64) error {
	// 删除插件钩子
	err := s.pluginHookModel.Delete(id)
	if err != nil {
		return err
	}

	// 清除缓存
	s.ClearPluginHookCache()

	return nil
}

// UpdatePluginHookStatus 更新插件钩子状态
func (s *PluginService) UpdatePluginHookStatus(id int64, status int) error {
	// 更新插件钩子状态
	err := s.pluginHookModel.UpdateStatus(id, status)
	if err != nil {
		return err
	}

	// 清除缓存
	s.ClearPluginHookCache()

	return nil
}

// ClearPluginCache 清除插件缓存
func (s *PluginService) ClearPluginCache() {
	// 清除插件缓存
	s.cache.Delete("plugin:all:0")
	s.cache.Delete("plugin:all:1")
	s.cache.Delete("plugin:all:-1")
}

// ClearPluginHookCache 清除插件钩子缓存
func (s *PluginService) ClearPluginHookCache() {
	// 清除插件钩子缓存
	s.cache.Delete("plugin:hook:all:0")
	s.cache.Delete("plugin:hook:all:1")
	s.cache.Delete("plugin:hook:all:-1")
}

// ExecuteHook 执行钩子
func (s *PluginService) ExecuteHook(position string, data map[string]interface{}) (string, error) {
	// 获取插件钩子
	hooks, err := s.GetPluginHooksByPosition(position)
	if err != nil {
		return "", err
	}

	// 执行钩子
	var result string
	for _, hook := range hooks {
		// 获取插件
		plugin, err := s.GetPlugin(hook.PluginID)
		if err != nil {
			logger.Error("获取插件失败", "id", hook.PluginID, "error", err)
			continue
		}

		// 检查插件状态
		if plugin.Status != 1 {
			continue
		}

		// 执行钩子
		hookResult, err := s.executePluginHook(plugin, hook, data)
		if err != nil {
			logger.Error("执行钩子失败", "plugin", plugin.Code, "hook", hook.Code, "error", err)
			continue
		}

		// 添加结果
		result += hookResult
	}

	return result, nil
}

// executePluginHook 执行插件钩子
func (s *PluginService) executePluginHook(plugin *model.Plugin, hook *model.PluginHook, data map[string]interface{}) (string, error) {
	// 构建插件路径
	pluginPath := filepath.Join(s.pluginDir, plugin.Code)
	hookPath := filepath.Join(pluginPath, "hooks", hook.Code+".go")

	// 检查钩子文件是否存在
	_, err := os.Stat(hookPath)
	if err != nil {
		return "", err
	}

	// 读取钩子文件
	hookContent, err := ioutil.ReadFile(hookPath)
	if err != nil {
		return "", err
	}

	// 执行钩子
	// 这里应该实现插件钩子的执行逻辑
	// 例如：使用Go插件机制、脚本引擎等
	// 这里简单返回钩子内容
	return string(hookContent), nil
}

// InstallPlugin 安装插件
func (s *PluginService) InstallPlugin(pluginCode string) error {
	// 构建插件路径
	pluginPath := filepath.Join(s.pluginDir, pluginCode)

	// 检查插件目录是否存在
	_, err := os.Stat(pluginPath)
	if err != nil {
		return err
	}

	// 读取插件配置
	configPath := filepath.Join(pluginPath, "config.json")
	configContent, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	// 解析插件配置
	var pluginConfig struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		Version     string `json:"version"`
		Author      string `json:"author"`
		Description string `json:"description"`
		Config      string `json:"config"`
		Hooks       []struct {
			Name     string `json:"name"`
			Code     string `json:"code"`
			Position string `json:"position"`
			OrderID  int    `json:"orderid"`
		} `json:"hooks"`
	}
	err = json.Unmarshal(configContent, &pluginConfig)
	if err != nil {
		return err
	}

	// 检查插件代码是否一致
	if pluginConfig.Code != pluginCode {
		return fmt.Errorf("plugin code mismatch: %s != %s", pluginConfig.Code, pluginCode)
	}

	// 检查插件是否已存在
	existingPlugin, err := s.pluginModel.GetByCode(pluginCode)
	if err == nil && existingPlugin != nil {
		// 插件已存在，更新
		existingPlugin.Name = pluginConfig.Name
		existingPlugin.Version = pluginConfig.Version
		existingPlugin.Author = pluginConfig.Author
		existingPlugin.Description = pluginConfig.Description
		existingPlugin.Config = pluginConfig.Config
		err = s.pluginModel.Update(existingPlugin)
		if err != nil {
			return err
		}

		// 删除插件钩子
		err = s.pluginHookModel.DeleteByPluginID(existingPlugin.ID)
		if err != nil {
			return err
		}

		// 创建插件钩子
		for _, hookConfig := range pluginConfig.Hooks {
			hook := &model.PluginHook{
				PluginID:  existingPlugin.ID,
				Name:      hookConfig.Name,
				Code:      hookConfig.Code,
				Position:  hookConfig.Position,
				OrderID:   hookConfig.OrderID,
				Status:    existingPlugin.Status,
			}
			_, err = s.pluginHookModel.Create(hook)
			if err != nil {
				return err
			}
		}
	} else {
		// 插件不存在，创建
		plugin := &model.Plugin{
			Name:        pluginConfig.Name,
			Code:        pluginConfig.Code,
			Version:     pluginConfig.Version,
			Author:      pluginConfig.Author,
			Description: pluginConfig.Description,
			Config:      pluginConfig.Config,
			Status:      1,
		}
		pluginID, err := s.pluginModel.Create(plugin)
		if err != nil {
			return err
		}

		// 创建插件钩子
		for _, hookConfig := range pluginConfig.Hooks {
			hook := &model.PluginHook{
				PluginID:  pluginID,
				Name:      hookConfig.Name,
				Code:      hookConfig.Code,
				Position:  hookConfig.Position,
				OrderID:   hookConfig.OrderID,
				Status:    1,
			}
			_, err = s.pluginHookModel.Create(hook)
			if err != nil {
				return err
			}
		}
	}

	// 清除缓存
	s.ClearPluginCache()
	s.ClearPluginHookCache()

	return nil
}

// UninstallPlugin 卸载插件
func (s *PluginService) UninstallPlugin(pluginCode string) error {
	// 获取插件
	plugin, err := s.pluginModel.GetByCode(pluginCode)
	if err != nil {
		return err
	}

	// 删除插件
	err = s.pluginModel.Delete(plugin.ID)
	if err != nil {
		return err
	}

	// 清除缓存
	s.ClearPluginCache()
	s.ClearPluginHookCache()

	return nil
}

// GetPluginConfig 获取插件配置
func (s *PluginService) GetPluginConfig(pluginCode string) (map[string]interface{}, error) {
	// 获取插件
	plugin, err := s.pluginModel.GetByCode(pluginCode)
	if err != nil {
		return nil, err
	}

	// 解析配置
	var config map[string]interface{}
	if plugin.Config != "" {
		err = json.Unmarshal([]byte(plugin.Config), &config)
		if err != nil {
			return nil, err
		}
	} else {
		config = make(map[string]interface{})
	}

	return config, nil
}

// SetPluginConfig 设置插件配置
func (s *PluginService) SetPluginConfig(pluginCode string, config map[string]interface{}) error {
	// 获取插件
	plugin, err := s.pluginModel.GetByCode(pluginCode)
	if err != nil {
		return err
	}

	// 序列化配置
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	// 更新配置
	plugin.Config = string(configJSON)
	err = s.pluginModel.Update(plugin)
	if err != nil {
		return err
	}

	// 清除缓存
	s.ClearPluginCache()

	return nil
}
