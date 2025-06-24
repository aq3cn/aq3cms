package plugin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"plugin"
	"sync"

	"aq3cms/pkg/logger"
)

// Plugin 插件接口
type Plugin interface {
	// Name 插件名称
	Name() string

	// Version 插件版本
	Version() string

	// Description 插件描述
	Description() string

	// Author 插件作者
	Author() string

	// Init 初始化插件
	Init(manager *Manager) error

	// Start 启动插件
	Start() error

	// Stop 停止插件
	Stop() error
}

// PluginInfo 插件信息
type PluginInfo struct {
	Name        string          `json:"name"`
	Version     string          `json:"version"`
	Description string          `json:"description"`
	Author      string          `json:"author"`
	Enabled     bool            `json:"enabled"`
	Config      json.RawMessage `json:"config"`
}

// Manager 插件管理器
type Manager struct {
	plugins     map[string]Plugin
	pluginInfos map[string]*PluginInfo
	pluginDir   string
	configFile  string
	mutex       sync.RWMutex
	hooks       map[string][]Hook
}

// Hook 钩子函数
type Hook struct {
	Name     string
	Priority int
	Callback func(args ...interface{}) interface{}
}

// NewManager 创建插件管理器
func NewManager(pluginDir, configFile string) *Manager {
	return &Manager{
		plugins:     make(map[string]Plugin),
		pluginInfos: make(map[string]*PluginInfo),
		pluginDir:   pluginDir,
		configFile:  configFile,
		hooks:       make(map[string][]Hook),
	}
}

// LoadPlugins 加载插件
func (m *Manager) LoadPlugins() error {
	// 检查插件目录是否为空
	if m.pluginDir == "" {
		logger.Warn("插件目录为空，跳过加载插件")
		return nil
	}

	// 创建插件目录
	if _, err := os.Stat(m.pluginDir); os.IsNotExist(err) {
		if err := os.MkdirAll(m.pluginDir, 0755); err != nil {
			return fmt.Errorf("创建插件目录失败: %w", err)
		}
	}

	// 加载插件配置
	if err := m.loadConfig(); err != nil {
		return fmt.Errorf("加载插件配置失败: %w", err)
	}

	// 加载插件
	files, err := ioutil.ReadDir(m.pluginDir)
	if err != nil {
		return fmt.Errorf("读取插件目录失败: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// 检查文件扩展名
		if filepath.Ext(file.Name()) != ".so" {
			continue
		}

		// 加载插件
		pluginPath := filepath.Join(m.pluginDir, file.Name())
		if err := m.loadPlugin(pluginPath); err != nil {
			logger.Error("加载插件失败", "path", pluginPath, "error", err)
			continue
		}
	}

	// 启动已启用的插件
	for name, info := range m.pluginInfos {
		if info.Enabled {
			if p, ok := m.plugins[name]; ok {
				if err := p.Start(); err != nil {
					logger.Error("启动插件失败", "name", name, "error", err)
				}
			}
		}
	}

	return nil
}

// loadPlugin 加载插件
func (m *Manager) loadPlugin(path string) error {
	// 加载插件
	p, err := plugin.Open(path)
	if err != nil {
		return err
	}

	// 获取插件实例
	symPlugin, err := p.Lookup("Plugin")
	if err != nil {
		return err
	}

	// 转换为插件接口
	plugin, ok := symPlugin.(Plugin)
	if !ok {
		return fmt.Errorf("invalid plugin: %s", path)
	}

	// 获取插件名称
	name := plugin.Name()

	// 检查插件是否已加载
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.plugins[name]; ok {
		return fmt.Errorf("plugin already loaded: %s", name)
	}

	// 获取插件信息
	info, ok := m.pluginInfos[name]
	if !ok {
		// 创建插件信息
		info = &PluginInfo{
			Name:        name,
			Version:     plugin.Version(),
			Description: plugin.Description(),
			Author:      plugin.Author(),
			Enabled:     false,
			Config:      json.RawMessage("{}"),
		}
		m.pluginInfos[name] = info
	}

	// 初始化插件
	if err := plugin.Init(m); err != nil {
		return err
	}

	// 添加到插件列表
	m.plugins[name] = plugin

	return nil
}

// loadConfig 加载插件配置
func (m *Manager) loadConfig() error {
	// 检查配置文件路径是否为空
	if m.configFile == "" {
		logger.Warn("插件配置文件路径为空，使用默认配置")
		return nil
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(m.configFile); os.IsNotExist(err) {
		// 创建默认配置
		return m.saveConfig()
	}

	// 读取配置文件
	data, err := ioutil.ReadFile(m.configFile)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析配置
	var pluginInfos map[string]*PluginInfo
	if err := json.Unmarshal(data, &pluginInfos); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 更新插件信息
	m.mutex.Lock()
	m.pluginInfos = pluginInfos
	m.mutex.Unlock()

	return nil
}

// saveConfig 保存插件配置
func (m *Manager) saveConfig() error {
	// 检查配置文件路径是否为空
	if m.configFile == "" {
		logger.Warn("插件配置文件路径为空，跳过保存配置")
		return nil
	}

	// 获取插件信息
	m.mutex.RLock()
	pluginInfos := m.pluginInfos
	m.mutex.RUnlock()

	// 序列化配置
	data, err := json.MarshalIndent(pluginInfos, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	// 确保配置文件的父目录存在
	configDir := filepath.Dir(m.configFile)
	if configDir != "" && configDir != "." {
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			if err := os.MkdirAll(configDir, 0755); err != nil {
				logger.Error("创建插件配置目录失败", "dir", configDir, "error", err)
				return fmt.Errorf("创建插件配置目录失败: %w", err)
			}
		}
	}

	// 保存配置文件
	if err := ioutil.WriteFile(m.configFile, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

// GetPlugin 获取插件
func (m *Manager) GetPlugin(name string) (Plugin, error) {
	// 获取插件
	m.mutex.RLock()
	plugin, ok := m.plugins[name]
	m.mutex.RUnlock()

	if !ok {
		return nil, fmt.Errorf("plugin not found: %s", name)
	}

	return plugin, nil
}

// GetPluginInfo 获取插件信息
func (m *Manager) GetPluginInfo(name string) (*PluginInfo, error) {
	// 获取插件信息
	m.mutex.RLock()
	info, ok := m.pluginInfos[name]
	m.mutex.RUnlock()

	if !ok {
		return nil, fmt.Errorf("plugin info not found: %s", name)
	}

	return info, nil
}

// GetPlugins 获取所有插件
func (m *Manager) GetPlugins() map[string]Plugin {
	// 获取插件列表
	m.mutex.RLock()
	plugins := make(map[string]Plugin, len(m.plugins))
	for name, plugin := range m.plugins {
		plugins[name] = plugin
	}
	m.mutex.RUnlock()

	return plugins
}

// GetPluginInfos 获取所有插件信息
func (m *Manager) GetPluginInfos() map[string]*PluginInfo {
	// 获取插件信息列表
	m.mutex.RLock()
	infos := make(map[string]*PluginInfo, len(m.pluginInfos))
	for name, info := range m.pluginInfos {
		infos[name] = info
	}
	m.mutex.RUnlock()

	return infos
}

// EnablePlugin 启用插件
func (m *Manager) EnablePlugin(name string) error {
	// 获取插件
	plugin, err := m.GetPlugin(name)
	if err != nil {
		return err
	}

	// 获取插件信息
	info, err := m.GetPluginInfo(name)
	if err != nil {
		return err
	}

	// 检查插件是否已启用
	if info.Enabled {
		return nil
	}

	// 启动插件
	if err := plugin.Start(); err != nil {
		return err
	}

	// 更新插件信息
	m.mutex.Lock()
	info.Enabled = true
	m.mutex.Unlock()

	// 保存配置
	return m.saveConfig()
}

// DisablePlugin 禁用插件
func (m *Manager) DisablePlugin(name string) error {
	// 获取插件
	plugin, err := m.GetPlugin(name)
	if err != nil {
		return err
	}

	// 获取插件信息
	info, err := m.GetPluginInfo(name)
	if err != nil {
		return err
	}

	// 检查插件是否已禁用
	if !info.Enabled {
		return nil
	}

	// 停止插件
	if err := plugin.Stop(); err != nil {
		return err
	}

	// 更新插件信息
	m.mutex.Lock()
	info.Enabled = false
	m.mutex.Unlock()

	// 保存配置
	return m.saveConfig()
}

// UpdatePluginConfig 更新插件配置
func (m *Manager) UpdatePluginConfig(name string, config json.RawMessage) error {
	// 获取插件信息
	info, err := m.GetPluginInfo(name)
	if err != nil {
		return err
	}

	// 更新插件配置
	m.mutex.Lock()
	info.Config = config
	m.mutex.Unlock()

	// 保存配置
	return m.saveConfig()
}

// AddHook 添加钩子
func (m *Manager) AddHook(hookName string, hook Hook) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 获取钩子列表
	hooks, ok := m.hooks[hookName]
	if !ok {
		hooks = make([]Hook, 0)
	}

	// 添加钩子
	hooks = append(hooks, hook)

	// 按优先级排序
	for i := 0; i < len(hooks)-1; i++ {
		for j := i + 1; j < len(hooks); j++ {
			if hooks[i].Priority < hooks[j].Priority {
				hooks[i], hooks[j] = hooks[j], hooks[i]
			}
		}
	}

	// 更新钩子列表
	m.hooks[hookName] = hooks
}

// RemoveHook 移除钩子
func (m *Manager) RemoveHook(hookName, name string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 获取钩子列表
	hooks, ok := m.hooks[hookName]
	if !ok {
		return
	}

	// 移除钩子
	newHooks := make([]Hook, 0, len(hooks))
	for _, hook := range hooks {
		if hook.Name != name {
			newHooks = append(newHooks, hook)
		}
	}

	// 更新钩子列表
	m.hooks[hookName] = newHooks
}

// ApplyHooks 应用钩子
func (m *Manager) ApplyHooks(hookName string, args ...interface{}) interface{} {
	m.mutex.RLock()
	hooks, ok := m.hooks[hookName]
	m.mutex.RUnlock()

	if !ok {
		return nil
	}

	var result interface{}
	for _, hook := range hooks {
		result = hook.Callback(args...)
	}

	return result
}

// HasHook 检查钩子是否存在
func (m *Manager) HasHook(hookName string) bool {
	m.mutex.RLock()
	_, ok := m.hooks[hookName]
	m.mutex.RUnlock()

	return ok
}

// GetHooks 获取钩子列表
func (m *Manager) GetHooks(hookName string) []Hook {
	m.mutex.RLock()
	hooks, ok := m.hooks[hookName]
	m.mutex.RUnlock()

	if !ok {
		return nil
	}

	return hooks
}
