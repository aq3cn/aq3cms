package model

import (
	"fmt"
	"time"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// Plugin 插件
type Plugin struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`        // 插件名称
	Code        string    `json:"code"`        // 插件代码
	Version     string    `json:"version"`     // 版本
	Author      string    `json:"author"`      // 作者
	Description string    `json:"description"` // 描述
	Config      string    `json:"config"`      // 配置，JSON格式
	Status      int       `json:"status"`      // 状态：0禁用，1启用
	CreateTime  time.Time `json:"createtime"`  // 创建时间
	UpdateTime  time.Time `json:"updatetime"`  // 更新时间
}

// PluginHook 插件钩子
type PluginHook struct {
	ID         int64     `json:"id"`
	PluginID   int64     `json:"pluginid"`   // 插件ID
	Name       string    `json:"name"`       // 钩子名称
	Code       string    `json:"code"`       // 钩子代码
	Position   string    `json:"position"`   // 钩子位置
	OrderID    int       `json:"orderid"`    // 排序ID
	Status     int       `json:"status"`     // 状态：0禁用，1启用
	CreateTime time.Time `json:"createtime"` // 创建时间
	UpdateTime time.Time `json:"updatetime"` // 更新时间
}

// PluginModel 插件模型
type PluginModel struct {
	db *database.DB
}

// NewPluginModel 创建插件模型
func NewPluginModel(db *database.DB) *PluginModel {
	return &PluginModel{
		db: db,
	}
}

// GetByID 根据ID获取插件
func (m *PluginModel) GetByID(id int64) (*Plugin, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "plugin")
	qb.Where("id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取插件失败", "id", id, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("plugin not found: %d", id)
	}

	// 转换为插件
	plugin := &Plugin{}
	plugin.ID, _ = result["id"].(int64)
	plugin.Name, _ = result["name"].(string)
	plugin.Code, _ = result["code"].(string)
	plugin.Version, _ = result["version"].(string)
	plugin.Author, _ = result["author"].(string)
	plugin.Description, _ = result["description"].(string)
	plugin.Config, _ = result["config"].(string)
	plugin.Status, _ = result["status"].(int)
	plugin.CreateTime, _ = result["createtime"].(time.Time)
	plugin.UpdateTime, _ = result["updatetime"].(time.Time)

	return plugin, nil
}

// GetByCode 根据代码获取插件
func (m *PluginModel) GetByCode(code string) (*Plugin, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "plugin")
	qb.Where("code = ?", code)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取插件失败", "code", code, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("plugin not found: %s", code)
	}

	// 转换为插件
	plugin := &Plugin{}
	plugin.ID, _ = result["id"].(int64)
	plugin.Name, _ = result["name"].(string)
	plugin.Code, _ = result["code"].(string)
	plugin.Version, _ = result["version"].(string)
	plugin.Author, _ = result["author"].(string)
	plugin.Description, _ = result["description"].(string)
	plugin.Config, _ = result["config"].(string)
	plugin.Status, _ = result["status"].(int)
	plugin.CreateTime, _ = result["createtime"].(time.Time)
	plugin.UpdateTime, _ = result["updatetime"].(time.Time)

	return plugin, nil
}

// GetAll 获取所有插件
func (m *PluginModel) GetAll(status int) ([]*Plugin, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "plugin")
	if status >= 0 {
		qb.Where("status = ?", status)
	}
	qb.OrderBy("id ASC")

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取所有插件失败", "error", err)
		return nil, err
	}

	// 转换为插件列表
	plugins := make([]*Plugin, 0, len(results))
	for _, result := range results {
		plugin := &Plugin{}
		plugin.ID, _ = result["id"].(int64)
		plugin.Name, _ = result["name"].(string)
		plugin.Code, _ = result["code"].(string)
		plugin.Version, _ = result["version"].(string)
		plugin.Author, _ = result["author"].(string)
		plugin.Description, _ = result["description"].(string)
		plugin.Config, _ = result["config"].(string)
		plugin.Status, _ = result["status"].(int)
		plugin.CreateTime, _ = result["createtime"].(time.Time)
		plugin.UpdateTime, _ = result["updatetime"].(time.Time)
		plugins = append(plugins, plugin)
	}

	return plugins, nil
}

// Create 创建插件
func (m *PluginModel) Create(plugin *Plugin) (int64, error) {
	// 设置创建时间和更新时间
	now := time.Now()
	plugin.CreateTime = now
	plugin.UpdateTime = now

	// 执行插入
	result, err := m.db.Exec(
		"INSERT INTO "+m.db.TableName("plugin")+" (name, code, version, author, description, config, status, createtime, updatetime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		plugin.Name, plugin.Code, plugin.Version, plugin.Author, plugin.Description, plugin.Config, plugin.Status, plugin.CreateTime, plugin.UpdateTime,
	)
	if err != nil {
		logger.Error("创建插件失败", "error", err)
		return 0, err
	}

	// 获取插入ID
	id, err := result.LastInsertId()
	if err != nil {
		logger.Error("获取插入ID失败", "error", err)
		return 0, err
	}

	return id, nil
}

// Update 更新插件
func (m *PluginModel) Update(plugin *Plugin) error {
	// 设置更新时间
	plugin.UpdateTime = time.Now()

	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("plugin")+" SET name = ?, code = ?, version = ?, author = ?, description = ?, config = ?, status = ?, updatetime = ? WHERE id = ?",
		plugin.Name, plugin.Code, plugin.Version, plugin.Author, plugin.Description, plugin.Config, plugin.Status, plugin.UpdateTime, plugin.ID,
	)
	if err != nil {
		logger.Error("更新插件失败", "error", err)
		return err
	}

	return nil
}

// Delete 删除插件
func (m *PluginModel) Delete(id int64) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("plugin")+" WHERE id = ?", id)
	if err != nil {
		logger.Error("删除插件失败", "error", err)
		return err
	}

	// 删除插件钩子
	_, err = m.db.Exec("DELETE FROM "+m.db.TableName("plugin_hook")+" WHERE pluginid = ?", id)
	if err != nil {
		logger.Error("删除插件钩子失败", "error", err)
		return err
	}

	return nil
}

// UpdateStatus 更新插件状态
func (m *PluginModel) UpdateStatus(id int64, status int) error {
	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("plugin")+" SET status = ?, updatetime = ? WHERE id = ?",
		status, time.Now(), id,
	)
	if err != nil {
		logger.Error("更新插件状态失败", "error", err)
		return err
	}

	// 更新插件钩子状态
	_, err = m.db.Exec(
		"UPDATE "+m.db.TableName("plugin_hook")+" SET status = ?, updatetime = ? WHERE pluginid = ?",
		status, time.Now(), id,
	)
	if err != nil {
		logger.Error("更新插件钩子状态失败", "error", err)
		return err
	}

	return nil
}

// PluginHookModel 插件钩子模型
type PluginHookModel struct {
	db *database.DB
}

// NewPluginHookModel 创建插件钩子模型
func NewPluginHookModel(db *database.DB) *PluginHookModel {
	return &PluginHookModel{
		db: db,
	}
}

// GetByID 根据ID获取插件钩子
func (m *PluginHookModel) GetByID(id int64) (*PluginHook, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "plugin_hook")
	qb.Where("id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取插件钩子失败", "id", id, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("plugin hook not found: %d", id)
	}

	// 转换为插件钩子
	hook := &PluginHook{}
	hook.ID, _ = result["id"].(int64)
	hook.PluginID, _ = result["pluginid"].(int64)
	hook.Name, _ = result["name"].(string)
	hook.Code, _ = result["code"].(string)
	hook.Position, _ = result["position"].(string)
	hook.OrderID, _ = result["orderid"].(int)
	hook.Status, _ = result["status"].(int)
	hook.CreateTime, _ = result["createtime"].(time.Time)
	hook.UpdateTime, _ = result["updatetime"].(time.Time)

	return hook, nil
}

// GetByCode 根据代码获取插件钩子
func (m *PluginHookModel) GetByCode(code string) (*PluginHook, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "plugin_hook")
	qb.Where("code = ?", code)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取插件钩子失败", "code", code, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("plugin hook not found: %s", code)
	}

	// 转换为插件钩子
	hook := &PluginHook{}
	hook.ID, _ = result["id"].(int64)
	hook.PluginID, _ = result["pluginid"].(int64)
	hook.Name, _ = result["name"].(string)
	hook.Code, _ = result["code"].(string)
	hook.Position, _ = result["position"].(string)
	hook.OrderID, _ = result["orderid"].(int)
	hook.Status, _ = result["status"].(int)
	hook.CreateTime, _ = result["createtime"].(time.Time)
	hook.UpdateTime, _ = result["updatetime"].(time.Time)

	return hook, nil
}

// GetByPluginID 根据插件ID获取插件钩子
func (m *PluginHookModel) GetByPluginID(pluginID int64) ([]*PluginHook, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "plugin_hook")
	qb.Where("pluginid = ?", pluginID)
	qb.OrderBy("orderid ASC, id ASC")

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取插件钩子失败", "pluginid", pluginID, "error", err)
		return nil, err
	}

	// 转换为插件钩子列表
	hooks := make([]*PluginHook, 0, len(results))
	for _, result := range results {
		hook := &PluginHook{}
		hook.ID, _ = result["id"].(int64)
		hook.PluginID, _ = result["pluginid"].(int64)
		hook.Name, _ = result["name"].(string)
		hook.Code, _ = result["code"].(string)
		hook.Position, _ = result["position"].(string)
		hook.OrderID, _ = result["orderid"].(int)
		hook.Status, _ = result["status"].(int)
		hook.CreateTime, _ = result["createtime"].(time.Time)
		hook.UpdateTime, _ = result["updatetime"].(time.Time)
		hooks = append(hooks, hook)
	}

	return hooks, nil
}

// GetByPosition 根据位置获取插件钩子
func (m *PluginHookModel) GetByPosition(position string) ([]*PluginHook, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "plugin_hook")
	qb.Where("position = ?", position)
	qb.Where("status = ?", 1)
	qb.OrderBy("orderid ASC, id ASC")

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取插件钩子失败", "position", position, "error", err)
		return nil, err
	}

	// 转换为插件钩子列表
	hooks := make([]*PluginHook, 0, len(results))
	for _, result := range results {
		hook := &PluginHook{}
		hook.ID, _ = result["id"].(int64)
		hook.PluginID, _ = result["pluginid"].(int64)
		hook.Name, _ = result["name"].(string)
		hook.Code, _ = result["code"].(string)
		hook.Position, _ = result["position"].(string)
		hook.OrderID, _ = result["orderid"].(int)
		hook.Status, _ = result["status"].(int)
		hook.CreateTime, _ = result["createtime"].(time.Time)
		hook.UpdateTime, _ = result["updatetime"].(time.Time)
		hooks = append(hooks, hook)
	}

	return hooks, nil
}

// GetAll 获取所有插件钩子
func (m *PluginHookModel) GetAll(status int) ([]*PluginHook, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "plugin_hook")
	if status >= 0 {
		qb.Where("status = ?", status)
	}
	qb.OrderBy("position ASC, orderid ASC, id ASC")

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取所有插件钩子失败", "error", err)
		return nil, err
	}

	// 转换为插件钩子列表
	hooks := make([]*PluginHook, 0, len(results))
	for _, result := range results {
		hook := &PluginHook{}
		hook.ID, _ = result["id"].(int64)
		hook.PluginID, _ = result["pluginid"].(int64)
		hook.Name, _ = result["name"].(string)
		hook.Code, _ = result["code"].(string)
		hook.Position, _ = result["position"].(string)
		hook.OrderID, _ = result["orderid"].(int)
		hook.Status, _ = result["status"].(int)
		hook.CreateTime, _ = result["createtime"].(time.Time)
		hook.UpdateTime, _ = result["updatetime"].(time.Time)
		hooks = append(hooks, hook)
	}

	return hooks, nil
}

// Create 创建插件钩子
func (m *PluginHookModel) Create(hook *PluginHook) (int64, error) {
	// 设置创建时间和更新时间
	now := time.Now()
	hook.CreateTime = now
	hook.UpdateTime = now

	// 执行插入
	result, err := m.db.Exec(
		"INSERT INTO "+m.db.TableName("plugin_hook")+" (pluginid, name, code, position, orderid, status, createtime, updatetime) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		hook.PluginID, hook.Name, hook.Code, hook.Position, hook.OrderID, hook.Status, hook.CreateTime, hook.UpdateTime,
	)
	if err != nil {
		logger.Error("创建插件钩子失败", "error", err)
		return 0, err
	}

	// 获取插入ID
	id, err := result.LastInsertId()
	if err != nil {
		logger.Error("获取插入ID失败", "error", err)
		return 0, err
	}

	return id, nil
}

// Update 更新插件钩子
func (m *PluginHookModel) Update(hook *PluginHook) error {
	// 设置更新时间
	hook.UpdateTime = time.Now()

	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("plugin_hook")+" SET pluginid = ?, name = ?, code = ?, position = ?, orderid = ?, status = ?, updatetime = ? WHERE id = ?",
		hook.PluginID, hook.Name, hook.Code, hook.Position, hook.OrderID, hook.Status, hook.UpdateTime, hook.ID,
	)
	if err != nil {
		logger.Error("更新插件钩子失败", "error", err)
		return err
	}

	return nil
}

// Delete 删除插件钩子
func (m *PluginHookModel) Delete(id int64) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("plugin_hook")+" WHERE id = ?", id)
	if err != nil {
		logger.Error("删除插件钩子失败", "error", err)
		return err
	}

	return nil
}

// DeleteByPluginID 根据插件ID删除插件钩子
func (m *PluginHookModel) DeleteByPluginID(pluginID int64) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("plugin_hook")+" WHERE pluginid = ?", pluginID)
	if err != nil {
		logger.Error("删除插件钩子失败", "pluginid", pluginID, "error", err)
		return err
	}

	return nil
}

// UpdateStatus 更新插件钩子状态
func (m *PluginHookModel) UpdateStatus(id int64, status int) error {
	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("plugin_hook")+" SET status = ?, updatetime = ? WHERE id = ?",
		status, time.Now(), id,
	)
	if err != nil {
		logger.Error("更新插件钩子状态失败", "error", err)
		return err
	}

	return nil
}

// UpdateStatusByPluginID 根据插件ID更新插件钩子状态
func (m *PluginHookModel) UpdateStatusByPluginID(pluginID int64, status int) error {
	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("plugin_hook")+" SET status = ?, updatetime = ? WHERE pluginid = ?",
		status, time.Now(), pluginID,
	)
	if err != nil {
		logger.Error("更新插件钩子状态失败", "pluginid", pluginID, "error", err)
		return err
	}

	return nil
}
