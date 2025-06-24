package model

import (
	"fmt"
	"time"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// Ad 广告
type Ad struct {
	ID         int64     `json:"id"`
	PositionID int64     `json:"positionid"` // 广告位ID
	Title      string    `json:"title"`      // 广告标题
	Type       int       `json:"type"`       // 广告类型：0图片，1Flash，2代码，3文字
	Image      string    `json:"image"`      // 图片地址
	Flash      string    `json:"flash"`      // Flash地址
	Code       string    `json:"code"`       // 代码内容
	Text       string    `json:"text"`       // 文字内容
	URL        string    `json:"url"`        // 链接地址
	StartTime  time.Time `json:"starttime"`  // 开始时间
	EndTime    time.Time `json:"endtime"`    // 结束时间
	OrderID    int       `json:"orderid"`    // 排序ID
	Status     int       `json:"status"`     // 状态：0禁用，1启用
	Target     string    `json:"target"`     // 打开方式：_blank, _self
	Width      int       `json:"width"`      // 宽度
	Height     int       `json:"height"`     // 高度
	Click      int       `json:"click"`      // 点击次数
	CreateTime time.Time `json:"createtime"` // 创建时间
	UpdateTime time.Time `json:"updatetime"` // 更新时间
}

// AdPosition 广告位
type AdPosition struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`       // 广告位名称
	Code       string    `json:"code"`       // 广告位代码
	Width      int       `json:"width"`      // 宽度
	Height     int       `json:"height"`     // 高度
	Template   string    `json:"template"`   // 模板
	Description string   `json:"description"` // 描述
	Status     int       `json:"status"`     // 状态：0禁用，1启用
	CreateTime time.Time `json:"createtime"` // 创建时间
	UpdateTime time.Time `json:"updatetime"` // 更新时间
}

// AdModel 广告模型
type AdModel struct {
	db *database.DB
}

// NewAdModel 创建广告模型
func NewAdModel(db *database.DB) *AdModel {
	return &AdModel{
		db: db,
	}
}

// GetByID 根据ID获取广告
func (m *AdModel) GetByID(id int64) (*Ad, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "ad")
	qb.Where("id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取广告失败", "id", id, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("ad not found: %d", id)
	}

	// 转换为广告
	ad := &Ad{}
	ad.ID, _ = result["id"].(int64)
	ad.PositionID, _ = result["positionid"].(int64)
	ad.Title, _ = result["title"].(string)
	ad.Type, _ = result["type"].(int)
	ad.Image, _ = result["image"].(string)
	ad.Flash, _ = result["flash"].(string)
	ad.Code, _ = result["code"].(string)
	ad.Text, _ = result["text"].(string)
	ad.URL, _ = result["url"].(string)
	ad.StartTime, _ = result["starttime"].(time.Time)
	ad.EndTime, _ = result["endtime"].(time.Time)
	ad.OrderID, _ = result["orderid"].(int)
	ad.Status, _ = result["status"].(int)
	ad.Target, _ = result["target"].(string)
	ad.Width, _ = result["width"].(int)
	ad.Height, _ = result["height"].(int)
	ad.Click, _ = result["click"].(int)
	ad.CreateTime, _ = result["createtime"].(time.Time)
	ad.UpdateTime, _ = result["updatetime"].(time.Time)

	return ad, nil
}

// GetByPositionID 根据广告位ID获取广告
func (m *AdModel) GetByPositionID(positionID int64) ([]*Ad, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "ad")
	qb.Where("positionid = ?", positionID)
	qb.Where("status = ?", 1)
	qb.Where("starttime <= ?", time.Now())
	qb.Where("endtime >= ?", time.Now())
	qb.OrderBy("orderid ASC, id ASC")

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取广告失败", "positionid", positionID, "error", err)
		return nil, err
	}

	// 转换为广告列表
	ads := make([]*Ad, 0, len(results))
	for _, result := range results {
		ad := &Ad{}
		ad.ID, _ = result["id"].(int64)
		ad.PositionID, _ = result["positionid"].(int64)
		ad.Title, _ = result["title"].(string)
		ad.Type, _ = result["type"].(int)
		ad.Image, _ = result["image"].(string)
		ad.Flash, _ = result["flash"].(string)
		ad.Code, _ = result["code"].(string)
		ad.Text, _ = result["text"].(string)
		ad.URL, _ = result["url"].(string)
		ad.StartTime, _ = result["starttime"].(time.Time)
		ad.EndTime, _ = result["endtime"].(time.Time)
		ad.OrderID, _ = result["orderid"].(int)
		ad.Status, _ = result["status"].(int)
		ad.Target, _ = result["target"].(string)
		ad.Width, _ = result["width"].(int)
		ad.Height, _ = result["height"].(int)
		ad.Click, _ = result["click"].(int)
		ad.CreateTime, _ = result["createtime"].(time.Time)
		ad.UpdateTime, _ = result["updatetime"].(time.Time)
		ads = append(ads, ad)
	}

	return ads, nil
}

// GetAll 获取所有广告
func (m *AdModel) GetAll(positionID int64, status int) ([]*Ad, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "ad")
	if positionID > 0 {
		qb.Where("positionid = ?", positionID)
	}
	if status >= 0 {
		qb.Where("status = ?", status)
	}
	qb.OrderBy("orderid ASC, id ASC")

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取所有广告失败", "error", err)
		return nil, err
	}

	// 转换为广告列表
	ads := make([]*Ad, 0, len(results))
	for _, result := range results {
		ad := &Ad{}
		ad.ID, _ = result["id"].(int64)
		ad.PositionID, _ = result["positionid"].(int64)
		ad.Title, _ = result["title"].(string)
		ad.Type, _ = result["type"].(int)
		ad.Image, _ = result["image"].(string)
		ad.Flash, _ = result["flash"].(string)
		ad.Code, _ = result["code"].(string)
		ad.Text, _ = result["text"].(string)
		ad.URL, _ = result["url"].(string)
		ad.StartTime, _ = result["starttime"].(time.Time)
		ad.EndTime, _ = result["endtime"].(time.Time)
		ad.OrderID, _ = result["orderid"].(int)
		ad.Status, _ = result["status"].(int)
		ad.Target, _ = result["target"].(string)
		ad.Width, _ = result["width"].(int)
		ad.Height, _ = result["height"].(int)
		ad.Click, _ = result["click"].(int)
		ad.CreateTime, _ = result["createtime"].(time.Time)
		ad.UpdateTime, _ = result["updatetime"].(time.Time)
		ads = append(ads, ad)
	}

	return ads, nil
}

// Create 创建广告
func (m *AdModel) Create(ad *Ad) (int64, error) {
	// 设置创建时间和更新时间
	now := time.Now()
	ad.CreateTime = now
	ad.UpdateTime = now

	// 执行插入
	result, err := m.db.Exec(
		"INSERT INTO "+m.db.TableName("ad")+" (positionid, title, type, image, flash, code, text, url, starttime, endtime, orderid, status, target, width, height, click, createtime, updatetime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		ad.PositionID, ad.Title, ad.Type, ad.Image, ad.Flash, ad.Code, ad.Text, ad.URL, ad.StartTime, ad.EndTime, ad.OrderID, ad.Status, ad.Target, ad.Width, ad.Height, ad.Click, ad.CreateTime, ad.UpdateTime,
	)
	if err != nil {
		logger.Error("创建广告失败", "error", err)
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

// Update 更新广告
func (m *AdModel) Update(ad *Ad) error {
	// 设置更新时间
	ad.UpdateTime = time.Now()

	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("ad")+" SET positionid = ?, title = ?, type = ?, image = ?, flash = ?, code = ?, text = ?, url = ?, starttime = ?, endtime = ?, orderid = ?, status = ?, target = ?, width = ?, height = ?, updatetime = ? WHERE id = ?",
		ad.PositionID, ad.Title, ad.Type, ad.Image, ad.Flash, ad.Code, ad.Text, ad.URL, ad.StartTime, ad.EndTime, ad.OrderID, ad.Status, ad.Target, ad.Width, ad.Height, ad.UpdateTime, ad.ID,
	)
	if err != nil {
		logger.Error("更新广告失败", "error", err)
		return err
	}

	return nil
}

// Delete 删除广告
func (m *AdModel) Delete(id int64) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("ad")+" WHERE id = ?", id)
	if err != nil {
		logger.Error("删除广告失败", "error", err)
		return err
	}

	return nil
}

// UpdateStatus 更新广告状态
func (m *AdModel) UpdateStatus(id int64, status int) error {
	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("ad")+" SET status = ?, updatetime = ? WHERE id = ?",
		status, time.Now(), id,
	)
	if err != nil {
		logger.Error("更新广告状态失败", "error", err)
		return err
	}

	return nil
}

// IncrementClick 增加点击次数
func (m *AdModel) IncrementClick(id int64) error {
	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("ad")+" SET click = click + 1 WHERE id = ?",
		id,
	)
	if err != nil {
		logger.Error("增加点击次数失败", "error", err)
		return err
	}

	return nil
}

// AdPositionModel 广告位模型
type AdPositionModel struct {
	db *database.DB
}

// NewAdPositionModel 创建广告位模型
func NewAdPositionModel(db *database.DB) *AdPositionModel {
	return &AdPositionModel{
		db: db,
	}
}

// GetByID 根据ID获取广告位
func (m *AdPositionModel) GetByID(id int64) (*AdPosition, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "ad_position")
	qb.Where("id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取广告位失败", "id", id, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("ad position not found: %d", id)
	}

	// 转换为广告位
	position := &AdPosition{}
	position.ID, _ = result["id"].(int64)
	position.Name, _ = result["name"].(string)
	position.Code, _ = result["code"].(string)
	position.Width, _ = result["width"].(int)
	position.Height, _ = result["height"].(int)
	position.Template, _ = result["template"].(string)
	position.Description, _ = result["description"].(string)
	position.Status, _ = result["status"].(int)
	position.CreateTime, _ = result["createtime"].(time.Time)
	position.UpdateTime, _ = result["updatetime"].(time.Time)

	return position, nil
}

// GetByCode 根据代码获取广告位
func (m *AdPositionModel) GetByCode(code string) (*AdPosition, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "ad_position")
	qb.Where("code = ?", code)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取广告位失败", "code", code, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("ad position not found: %s", code)
	}

	// 转换为广告位
	position := &AdPosition{}
	position.ID, _ = result["id"].(int64)
	position.Name, _ = result["name"].(string)
	position.Code, _ = result["code"].(string)
	position.Width, _ = result["width"].(int)
	position.Height, _ = result["height"].(int)
	position.Template, _ = result["template"].(string)
	position.Description, _ = result["description"].(string)
	position.Status, _ = result["status"].(int)
	position.CreateTime, _ = result["createtime"].(time.Time)
	position.UpdateTime, _ = result["updatetime"].(time.Time)

	return position, nil
}

// GetAll 获取所有广告位
func (m *AdPositionModel) GetAll(status int) ([]*AdPosition, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "ad_position")
	if status >= 0 {
		qb.Where("status = ?", status)
	}
	qb.OrderBy("id ASC")

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取所有广告位失败", "error", err)
		return nil, err
	}

	// 转换为广告位列表
	positions := make([]*AdPosition, 0, len(results))
	for _, result := range results {
		position := &AdPosition{}
		position.ID, _ = result["id"].(int64)
		position.Name, _ = result["name"].(string)
		position.Code, _ = result["code"].(string)
		position.Width, _ = result["width"].(int)
		position.Height, _ = result["height"].(int)
		position.Template, _ = result["template"].(string)
		position.Description, _ = result["description"].(string)
		position.Status, _ = result["status"].(int)
		position.CreateTime, _ = result["createtime"].(time.Time)
		position.UpdateTime, _ = result["updatetime"].(time.Time)
		positions = append(positions, position)
	}

	return positions, nil
}

// Create 创建广告位
func (m *AdPositionModel) Create(position *AdPosition) (int64, error) {
	// 设置创建时间和更新时间
	now := time.Now()
	position.CreateTime = now
	position.UpdateTime = now

	// 执行插入
	result, err := m.db.Exec(
		"INSERT INTO "+m.db.TableName("ad_position")+" (name, code, width, height, template, description, status, createtime, updatetime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		position.Name, position.Code, position.Width, position.Height, position.Template, position.Description, position.Status, position.CreateTime, position.UpdateTime,
	)
	if err != nil {
		logger.Error("创建广告位失败", "error", err)
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

// Update 更新广告位
func (m *AdPositionModel) Update(position *AdPosition) error {
	// 设置更新时间
	position.UpdateTime = time.Now()

	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("ad_position")+" SET name = ?, code = ?, width = ?, height = ?, template = ?, description = ?, status = ?, updatetime = ? WHERE id = ?",
		position.Name, position.Code, position.Width, position.Height, position.Template, position.Description, position.Status, position.UpdateTime, position.ID,
	)
	if err != nil {
		logger.Error("更新广告位失败", "error", err)
		return err
	}

	return nil
}

// Delete 删除广告位
func (m *AdPositionModel) Delete(id int64) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("ad_position")+" WHERE id = ?", id)
	if err != nil {
		logger.Error("删除广告位失败", "error", err)
		return err
	}

	return nil
}

// UpdateStatus 更新广告位状态
func (m *AdPositionModel) UpdateStatus(id int64, status int) error {
	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("ad_position")+" SET status = ?, updatetime = ? WHERE id = ?",
		status, time.Now(), id,
	)
	if err != nil {
		logger.Error("更新广告位状态失败", "error", err)
		return err
	}

	return nil
}

// HasAds 检查广告位是否有广告
func (m *AdPositionModel) HasAds(id int64) (bool, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "ad")
	qb.Where("positionid = ?", id)

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("检查广告位是否有广告失败", "error", err)
		return false, err
	}

	return count > 0, nil
}
