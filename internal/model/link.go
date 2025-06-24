package model

import (
	"fmt"
	"time"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// Link 友情链接
type Link struct {
	ID          int64     `json:"id"`
	TypeID      int64     `json:"typeid"`      // 分类ID
	Title       string    `json:"title"`       // 链接标题
	URL         string    `json:"url"`         // 链接URL
	Logo        string    `json:"logo"`        // 链接LOGO
	Description string    `json:"description"` // 链接描述
	Email       string    `json:"email"`       // 联系邮箱
	OrderID     int       `json:"orderid"`     // 排序ID
	IsCheck     int       `json:"ischeck"`     // 状态：0禁用，1启用
	IsLogo      int       `json:"islogo"`      // 是否显示LOGO：0否，1是
	CreateTime  time.Time `json:"createtime"`  // 创建时间
	UpdateTime  time.Time `json:"updatetime"`  // 更新时间
}

// LinkType 友情链接分类
type LinkType struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`       // 分类名称
	OrderID    int       `json:"orderid"`    // 排序ID
	CreateTime time.Time `json:"createtime"` // 创建时间
	UpdateTime time.Time `json:"updatetime"` // 更新时间
}

// LinkModel 友情链接模型
type LinkModel struct {
	db *database.DB
}

// NewLinkModel 创建友情链接模型
func NewLinkModel(db *database.DB) *LinkModel {
	return &LinkModel{
		db: db,
	}
}

// GetByID 根据ID获取友情链接
func (m *LinkModel) GetByID(id int64) (*Link, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "flink")
	qb.Where("id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取友情链接失败", "id", id, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("link not found: %d", id)
	}

	// 转换为友情链接
	link := &Link{}
	link.ID, _ = result["id"].(int64)
	link.TypeID, _ = result["typeid"].(int64)
	link.Title, _ = result["title"].(string)
	link.URL, _ = result["url"].(string)
	link.Logo, _ = result["logo"].(string)
	link.Description, _ = result["description"].(string)
	link.Email, _ = result["email"].(string)
	link.OrderID, _ = result["orderid"].(int)
	link.IsCheck, _ = result["ischeck"].(int)
	link.IsLogo, _ = result["islogo"].(int)
	link.CreateTime, _ = result["createtime"].(time.Time)
	link.UpdateTime, _ = result["updatetime"].(time.Time)

	return link, nil
}

// GetAll 获取所有友情链接
func (m *LinkModel) GetAll(typeID int64, status int) ([]*Link, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "flink")
	if typeID > 0 {
		qb.Where("typeid = ?", typeID)
	}
	if status >= 0 {
		qb.Where("ischeck = ?", status)
	}
	qb.OrderBy("id ASC")

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取所有友情链接失败", "error", err)
		return nil, err
	}

	// 转换为友情链接列表
	links := make([]*Link, 0, len(results))
	for _, result := range results {
		link := &Link{}
		link.ID, _ = result["id"].(int64)
		link.TypeID, _ = result["typeid"].(int64)
		link.Title, _ = result["title"].(string)
		link.URL, _ = result["url"].(string)
		link.Logo, _ = result["logo"].(string)
		link.Description, _ = result["description"].(string)
		link.Email, _ = result["email"].(string)
		link.OrderID, _ = result["orderid"].(int)
		link.IsCheck, _ = result["ischeck"].(int)
		link.IsLogo, _ = result["islogo"].(int)
		link.CreateTime, _ = result["createtime"].(time.Time)
		link.UpdateTime, _ = result["updatetime"].(time.Time)
		links = append(links, link)
	}

	return links, nil
}

// Create 创建友情链接
func (m *LinkModel) Create(link *Link) (int64, error) {
	// 设置创建时间和更新时间
	now := time.Now()
	link.CreateTime = now
	link.UpdateTime = now

	// 执行插入
	result, err := m.db.Exec(
		"INSERT INTO "+m.db.TableName("flink")+" (typeid, title, url, logo, description, email, orderid, ischeck, islogo, createtime, updatetime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		link.TypeID, link.Title, link.URL, link.Logo, link.Description, link.Email, link.OrderID, link.IsCheck, link.IsLogo, link.CreateTime, link.UpdateTime,
	)
	if err != nil {
		logger.Error("创建友情链接失败", "error", err)
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

// Update 更新友情链接
func (m *LinkModel) Update(link *Link) error {
	// 设置更新时间
	link.UpdateTime = time.Now()

	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("flink")+" SET typeid = ?, title = ?, url = ?, logo = ?, description = ?, email = ?, orderid = ?, ischeck = ?, islogo = ?, updatetime = ? WHERE id = ?",
		link.TypeID, link.Title, link.URL, link.Logo, link.Description, link.Email, link.OrderID, link.IsCheck, link.IsLogo, link.UpdateTime, link.ID,
	)
	if err != nil {
		logger.Error("更新友情链接失败", "error", err)
		return err
	}

	return nil
}

// Delete 删除友情链接
func (m *LinkModel) Delete(id int64) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("flink")+" WHERE id = ?", id)
	if err != nil {
		logger.Error("删除友情链接失败", "error", err)
		return err
	}

	return nil
}

// UpdateStatus 更新友情链接状态
func (m *LinkModel) UpdateStatus(id int64, status int) error {
	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("flink")+" SET ischeck = ?, updatetime = ? WHERE id = ?",
		status, time.Now(), id,
	)
	if err != nil {
		logger.Error("更新友情链接状态失败", "error", err)
		return err
	}

	return nil
}

// LinkTypeModel 友情链接分类模型
type LinkTypeModel struct {
	db *database.DB
}

// NewLinkTypeModel 创建友情链接分类模型
func NewLinkTypeModel(db *database.DB) *LinkTypeModel {
	return &LinkTypeModel{
		db: db,
	}
}

// GetByID 根据ID获取友情链接分类
func (m *LinkTypeModel) GetByID(id int64) (*LinkType, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "link_type")
	qb.Where("id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取友情链接分类失败", "id", id, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("link type not found: %d", id)
	}

	// 转换为友情链接分类
	linkType := &LinkType{}
	linkType.ID, _ = result["id"].(int64)
	linkType.Name, _ = result["name"].(string)
	linkType.OrderID, _ = result["orderid"].(int)
	linkType.CreateTime, _ = result["createtime"].(time.Time)
	linkType.UpdateTime, _ = result["updatetime"].(time.Time)

	return linkType, nil
}

// GetAll 获取所有友情链接分类
func (m *LinkTypeModel) GetAll() ([]*LinkType, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "link_type")
	qb.OrderBy("orderid ASC, id ASC")

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取所有友情链接分类失败", "error", err)
		return nil, err
	}

	// 转换为友情链接分类列表
	linkTypes := make([]*LinkType, 0, len(results))
	for _, result := range results {
		linkType := &LinkType{}
		linkType.ID, _ = result["id"].(int64)
		linkType.Name, _ = result["name"].(string)
		linkType.OrderID, _ = result["orderid"].(int)
		linkType.CreateTime, _ = result["createtime"].(time.Time)
		linkType.UpdateTime, _ = result["updatetime"].(time.Time)
		linkTypes = append(linkTypes, linkType)
	}

	return linkTypes, nil
}

// Create 创建友情链接分类
func (m *LinkTypeModel) Create(linkType *LinkType) (int64, error) {
	// 设置创建时间和更新时间
	now := time.Now()
	linkType.CreateTime = now
	linkType.UpdateTime = now

	// 执行插入
	result, err := m.db.Exec(
		"INSERT INTO "+m.db.TableName("link_type")+" (name, orderid, createtime, updatetime) VALUES (?, ?, ?, ?)",
		linkType.Name, linkType.OrderID, linkType.CreateTime, linkType.UpdateTime,
	)
	if err != nil {
		logger.Error("创建友情链接分类失败", "error", err)
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

// Update 更新友情链接分类
func (m *LinkTypeModel) Update(linkType *LinkType) error {
	// 设置更新时间
	linkType.UpdateTime = time.Now()

	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("link_type")+" SET name = ?, orderid = ?, updatetime = ? WHERE id = ?",
		linkType.Name, linkType.OrderID, linkType.UpdateTime, linkType.ID,
	)
	if err != nil {
		logger.Error("更新友情链接分类失败", "error", err)
		return err
	}

	return nil
}

// Delete 删除友情链接分类
func (m *LinkTypeModel) Delete(id int64) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("link_type")+" WHERE id = ?", id)
	if err != nil {
		logger.Error("删除友情链接分类失败", "error", err)
		return err
	}

	return nil
}

// HasLinks 检查分类是否有友情链接
func (m *LinkTypeModel) HasLinks(id int64) (bool, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "flink")
	qb.Where("typeid = ?", id)

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("检查分类是否有友情链接失败", "error", err)
		return false, err
	}

	return count > 0, nil
}

// GetLinks 获取友情链接列表
func (m *LinkModel) GetLinks(limit int) ([]*Link, error) {
	return m.GetAll(0, 1)
}

// GetByType 根据类型获取友情链接
func (m *LinkModel) GetByType(typeID int64, limit int) ([]*Link, error) {
	return m.GetAll(typeID, 1)
}
