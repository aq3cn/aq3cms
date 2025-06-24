package model

import (
	"encoding/json"
	"fmt"
	"time"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// FormDataEntry 表单数据条目
type FormDataEntry struct {
	ID         int64     `json:"id"`
	FormID     int64     `json:"formid"`     // 表单ID
	Data       string    `json:"data"`       // 表单数据
	IP         string    `json:"ip"`         // IP地址
	Status     int       `json:"status"`     // 状态：0未处理，1已处理
	CreateTime time.Time `json:"createtime"` // 创建时间
	UpdateTime time.Time `json:"updatetime"` // 更新时间
}

// FormDataModel 表单数据模型
type FormDataModel struct {
	db *database.DB
}

// NewFormDataModel 创建表单数据模型
func NewFormDataModel(db *database.DB) *FormDataModel {
	return &FormDataModel{
		db: db,
	}
}

// GetByID 根据ID获取表单数据
func (m *FormDataModel) GetByID(id int64) (*FormDataEntry, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "form_data")
	qb.Where("id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取表单数据失败", "id", id, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("form data not found: %d", id)
	}

	// 转换为表单数据
	formData := &FormDataEntry{}
	formData.ID, _ = result["id"].(int64)
	formData.FormID, _ = result["formid"].(int64)
	formData.Data, _ = result["data"].(string)
	formData.IP, _ = result["ip"].(string)
	formData.Status, _ = result["status"].(int)
	formData.CreateTime, _ = result["createtime"].(time.Time)
	formData.UpdateTime, _ = result["updatetime"].(time.Time)

	return formData, nil
}

// GetByFormID 根据表单ID获取表单数据
func (m *FormDataModel) GetByFormID(formID int64, status int, page, pageSize int) ([]*FormDataEntry, int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "form_data")
	qb.Where("formid = ?", formID)
	if status >= 0 {
		qb.Where("status = ?", status)
	}
	qb.OrderBy("id DESC")
	qb.Limit(pageSize)
	qb.Offset((page - 1) * pageSize)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取表单数据失败", "formid", formID, "error", err)
		return nil, 0, err
	}

	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("获取表单数据总数失败", "formid", formID, "error", err)
		return nil, 0, err
	}

	// 转换为表单数据列表
	formDataList := make([]*FormDataEntry, 0, len(results))
	for _, result := range results {
		formData := &FormDataEntry{}
		formData.ID, _ = result["id"].(int64)
		formData.FormID, _ = result["formid"].(int64)
		formData.Data, _ = result["data"].(string)
		formData.IP, _ = result["ip"].(string)
		formData.Status, _ = result["status"].(int)
		formData.CreateTime, _ = result["createtime"].(time.Time)
		formData.UpdateTime, _ = result["updatetime"].(time.Time)
		formDataList = append(formDataList, formData)
	}

	return formDataList, total, nil
}

// Create 创建表单数据
func (m *FormDataModel) Create(formData *FormDataEntry) (int64, error) {
	// 设置创建时间和更新时间
	now := time.Now()
	formData.CreateTime = now
	formData.UpdateTime = now

	// 执行插入
	result, err := m.db.Exec(
		"INSERT INTO "+m.db.TableName("form_data")+" (formid, data, ip, status, createtime, updatetime) VALUES (?, ?, ?, ?, ?, ?)",
		formData.FormID, formData.Data, formData.IP, formData.Status, formData.CreateTime, formData.UpdateTime,
	)
	if err != nil {
		logger.Error("创建表单数据失败", "error", err)
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

// Update 更新表单数据
func (m *FormDataModel) Update(formData *FormDataEntry) error {
	// 设置更新时间
	formData.UpdateTime = time.Now()

	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("form_data")+" SET formid = ?, data = ?, ip = ?, status = ?, updatetime = ? WHERE id = ?",
		formData.FormID, formData.Data, formData.IP, formData.Status, formData.UpdateTime, formData.ID,
	)
	if err != nil {
		logger.Error("更新表单数据失败", "error", err)
		return err
	}

	return nil
}

// UpdateStatus 更新表单数据状态
func (m *FormDataModel) UpdateStatus(id int64, status int) error {
	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("form_data")+" SET status = ?, updatetime = ? WHERE id = ?",
		status, time.Now(), id,
	)
	if err != nil {
		logger.Error("更新表单数据状态失败", "error", err)
		return err
	}

	return nil
}

// Delete 删除表单数据
func (m *FormDataModel) Delete(id int64) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("form_data")+" WHERE id = ?", id)
	if err != nil {
		logger.Error("删除表单数据失败", "error", err)
		return err
	}

	return nil
}

// DeleteByFormID 根据表单ID删除表单数据
func (m *FormDataModel) DeleteByFormID(formID int64) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("form_data")+" WHERE formid = ?", formID)
	if err != nil {
		logger.Error("删除表单数据失败", "formid", formID, "error", err)
		return err
	}

	return nil
}

// GetData 获取表单数据
func (m *FormDataModel) GetData(id int64) (map[string]interface{}, error) {
	// 获取表单数据
	formData, err := m.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 解析表单数据
	var data map[string]interface{}
	err = json.Unmarshal([]byte(formData.Data), &data)
	if err != nil {
		logger.Error("解析表单数据失败", "error", err)
		return nil, err
	}

	return data, nil
}
