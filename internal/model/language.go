package model

import (
	"fmt"
	"time"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// Language 语言
type Language struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`        // 语言名称
	Code        string    `json:"code"`        // 语言代码
	Flag        string    `json:"flag"`        // 国旗图标
	IsDefault   int       `json:"isdefault"`   // 是否默认：0否，1是
	Status      int       `json:"status"`      // 状态：0禁用，1启用
	OrderID     int       `json:"orderid"`     // 排序ID
	CreateTime  time.Time `json:"createtime"`  // 创建时间
	UpdateTime  time.Time `json:"updatetime"`  // 更新时间
}

// LanguageText 语言文本
type LanguageText struct {
	ID         int64     `json:"id"`
	LanguageID int64     `json:"languageid"` // 语言ID
	Key        string    `json:"key"`        // 键
	Value      string    `json:"value"`      // 值
	Module     string    `json:"module"`     // 模块
	CreateTime time.Time `json:"createtime"` // 创建时间
	UpdateTime time.Time `json:"updatetime"` // 更新时间
}

// LanguageModel 语言模型
type LanguageModel struct {
	db *database.DB
}

// NewLanguageModel 创建语言模型
func NewLanguageModel(db *database.DB) *LanguageModel {
	return &LanguageModel{
		db: db,
	}
}

// GetByID 根据ID获取语言
func (m *LanguageModel) GetByID(id int64) (*Language, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "language")
	qb.Where("id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取语言失败", "id", id, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("language not found: %d", id)
	}

	// 转换为语言
	language := &Language{}
	language.ID, _ = result["id"].(int64)
	language.Name, _ = result["name"].(string)
	language.Code, _ = result["code"].(string)
	language.Flag, _ = result["flag"].(string)
	language.IsDefault, _ = result["isdefault"].(int)
	language.Status, _ = result["status"].(int)
	language.OrderID, _ = result["orderid"].(int)
	language.CreateTime, _ = result["createtime"].(time.Time)
	language.UpdateTime, _ = result["updatetime"].(time.Time)

	return language, nil
}

// GetByCode 根据代码获取语言
func (m *LanguageModel) GetByCode(code string) (*Language, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "language")
	qb.Where("code = ?", code)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取语言失败", "code", code, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("language not found: %s", code)
	}

	// 转换为语言
	language := &Language{}
	language.ID, _ = result["id"].(int64)
	language.Name, _ = result["name"].(string)
	language.Code, _ = result["code"].(string)
	language.Flag, _ = result["flag"].(string)
	language.IsDefault, _ = result["isdefault"].(int)
	language.Status, _ = result["status"].(int)
	language.OrderID, _ = result["orderid"].(int)
	language.CreateTime, _ = result["createtime"].(time.Time)
	language.UpdateTime, _ = result["updatetime"].(time.Time)

	return language, nil
}

// GetDefault 获取默认语言
func (m *LanguageModel) GetDefault() (*Language, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "language")
	qb.Where("isdefault = ?", 1)
	qb.Where("status = ?", 1)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取默认语言失败", "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("default language not found")
	}

	// 转换为语言
	language := &Language{}
	language.ID, _ = result["id"].(int64)
	language.Name, _ = result["name"].(string)
	language.Code, _ = result["code"].(string)
	language.Flag, _ = result["flag"].(string)
	language.IsDefault, _ = result["isdefault"].(int)
	language.Status, _ = result["status"].(int)
	language.OrderID, _ = result["orderid"].(int)
	language.CreateTime, _ = result["createtime"].(time.Time)
	language.UpdateTime, _ = result["updatetime"].(time.Time)

	return language, nil
}

// GetAll 获取所有语言
func (m *LanguageModel) GetAll(status int) ([]*Language, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "language")
	if status >= 0 {
		qb.Where("status = ?", status)
	}
	qb.OrderBy("orderid ASC, id ASC")

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取所有语言失败", "error", err)
		return nil, err
	}

	// 转换为语言列表
	languages := make([]*Language, 0, len(results))
	for _, result := range results {
		language := &Language{}
		language.ID, _ = result["id"].(int64)
		language.Name, _ = result["name"].(string)
		language.Code, _ = result["code"].(string)
		language.Flag, _ = result["flag"].(string)
		language.IsDefault, _ = result["isdefault"].(int)
		language.Status, _ = result["status"].(int)
		language.OrderID, _ = result["orderid"].(int)
		language.CreateTime, _ = result["createtime"].(time.Time)
		language.UpdateTime, _ = result["updatetime"].(time.Time)
		languages = append(languages, language)
	}

	return languages, nil
}

// Create 创建语言
func (m *LanguageModel) Create(language *Language) (int64, error) {
	// 设置创建时间和更新时间
	now := time.Now()
	language.CreateTime = now
	language.UpdateTime = now

	// 如果是默认语言，更新其他语言为非默认
	if language.IsDefault == 1 {
		_, err := m.db.Exec("UPDATE "+m.db.TableName("language")+" SET isdefault = 0, updatetime = ? WHERE id != ?", now, language.ID)
		if err != nil {
			logger.Error("更新其他语言为非默认失败", "error", err)
			return 0, err
		}
	}

	// 执行插入
	result, err := m.db.Exec(
		"INSERT INTO "+m.db.TableName("language")+" (name, code, flag, isdefault, status, orderid, createtime, updatetime) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		language.Name, language.Code, language.Flag, language.IsDefault, language.Status, language.OrderID, language.CreateTime, language.UpdateTime,
	)
	if err != nil {
		logger.Error("创建语言失败", "error", err)
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

// Update 更新语言
func (m *LanguageModel) Update(language *Language) error {
	// 设置更新时间
	language.UpdateTime = time.Now()

	// 如果是默认语言，更新其他语言为非默认
	if language.IsDefault == 1 {
		_, err := m.db.Exec("UPDATE "+m.db.TableName("language")+" SET isdefault = 0, updatetime = ? WHERE id != ?", language.UpdateTime, language.ID)
		if err != nil {
			logger.Error("更新其他语言为非默认失败", "error", err)
			return err
		}
	}

	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("language")+" SET name = ?, code = ?, flag = ?, isdefault = ?, status = ?, orderid = ?, updatetime = ? WHERE id = ?",
		language.Name, language.Code, language.Flag, language.IsDefault, language.Status, language.OrderID, language.UpdateTime, language.ID,
	)
	if err != nil {
		logger.Error("更新语言失败", "error", err)
		return err
	}

	return nil
}

// Delete 删除语言
func (m *LanguageModel) Delete(id int64) error {
	// 检查是否是默认语言
	language, err := m.GetByID(id)
	if err != nil {
		return err
	}
	if language.IsDefault == 1 {
		return fmt.Errorf("cannot delete default language")
	}

	// 执行删除
	_, err = m.db.Exec("DELETE FROM "+m.db.TableName("language")+" WHERE id = ?", id)
	if err != nil {
		logger.Error("删除语言失败", "error", err)
		return err
	}

	// 删除语言文本
	_, err = m.db.Exec("DELETE FROM "+m.db.TableName("language_text")+" WHERE languageid = ?", id)
	if err != nil {
		logger.Error("删除语言文本失败", "error", err)
		return err
	}

	return nil
}

// UpdateStatus 更新语言状态
func (m *LanguageModel) UpdateStatus(id int64, status int) error {
	// 检查是否是默认语言
	language, err := m.GetByID(id)
	if err != nil {
		return err
	}
	if language.IsDefault == 1 && status == 0 {
		return fmt.Errorf("cannot disable default language")
	}

	// 执行更新
	_, err = m.db.Exec(
		"UPDATE "+m.db.TableName("language")+" SET status = ?, updatetime = ? WHERE id = ?",
		status, time.Now(), id,
	)
	if err != nil {
		logger.Error("更新语言状态失败", "error", err)
		return err
	}

	return nil
}

// LanguageTextModel 语言文本模型
type LanguageTextModel struct {
	db *database.DB
}

// NewLanguageTextModel 创建语言文本模型
func NewLanguageTextModel(db *database.DB) *LanguageTextModel {
	return &LanguageTextModel{
		db: db,
	}
}

// GetByID 根据ID获取语言文本
func (m *LanguageTextModel) GetByID(id int64) (*LanguageText, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "language_text")
	qb.Where("id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取语言文本失败", "id", id, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("language text not found: %d", id)
	}

	// 转换为语言文本
	text := &LanguageText{}
	text.ID, _ = result["id"].(int64)
	text.LanguageID, _ = result["languageid"].(int64)
	text.Key, _ = result["key"].(string)
	text.Value, _ = result["value"].(string)
	text.Module, _ = result["module"].(string)
	text.CreateTime, _ = result["createtime"].(time.Time)
	text.UpdateTime, _ = result["updatetime"].(time.Time)

	return text, nil
}

// GetByKey 根据键获取语言文本
func (m *LanguageTextModel) GetByKey(languageID int64, key string) (*LanguageText, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "language_text")
	qb.Where("languageid = ?", languageID)
	qb.Where("key = ?", key)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取语言文本失败", "languageid", languageID, "key", key, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("language text not found: %d, %s", languageID, key)
	}

	// 转换为语言文本
	text := &LanguageText{}
	text.ID, _ = result["id"].(int64)
	text.LanguageID, _ = result["languageid"].(int64)
	text.Key, _ = result["key"].(string)
	text.Value, _ = result["value"].(string)
	text.Module, _ = result["module"].(string)
	text.CreateTime, _ = result["createtime"].(time.Time)
	text.UpdateTime, _ = result["updatetime"].(time.Time)

	return text, nil
}

// GetByLanguageID 根据语言ID获取语言文本
func (m *LanguageTextModel) GetByLanguageID(languageID int64, module string) ([]*LanguageText, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "language_text")
	qb.Where("languageid = ?", languageID)
	if module != "" {
		qb.Where("module = ?", module)
	}
	qb.OrderBy("module ASC, key ASC")

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取语言文本失败", "languageid", languageID, "module", module, "error", err)
		return nil, err
	}

	// 转换为语言文本列表
	texts := make([]*LanguageText, 0, len(results))
	for _, result := range results {
		text := &LanguageText{}
		text.ID, _ = result["id"].(int64)
		text.LanguageID, _ = result["languageid"].(int64)
		text.Key, _ = result["key"].(string)
		text.Value, _ = result["value"].(string)
		text.Module, _ = result["module"].(string)
		text.CreateTime, _ = result["createtime"].(time.Time)
		text.UpdateTime, _ = result["updatetime"].(time.Time)
		texts = append(texts, text)
	}

	return texts, nil
}

// GetAll 获取所有语言文本
func (m *LanguageTextModel) GetAll(module string) ([]*LanguageText, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "language_text")
	if module != "" {
		qb.Where("module = ?", module)
	}
	qb.OrderBy("languageid ASC, module ASC, key ASC")

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取所有语言文本失败", "module", module, "error", err)
		return nil, err
	}

	// 转换为语言文本列表
	texts := make([]*LanguageText, 0, len(results))
	for _, result := range results {
		text := &LanguageText{}
		text.ID, _ = result["id"].(int64)
		text.LanguageID, _ = result["languageid"].(int64)
		text.Key, _ = result["key"].(string)
		text.Value, _ = result["value"].(string)
		text.Module, _ = result["module"].(string)
		text.CreateTime, _ = result["createtime"].(time.Time)
		text.UpdateTime, _ = result["updatetime"].(time.Time)
		texts = append(texts, text)
	}

	return texts, nil
}

// Create 创建语言文本
func (m *LanguageTextModel) Create(text *LanguageText) (int64, error) {
	// 设置创建时间和更新时间
	now := time.Now()
	text.CreateTime = now
	text.UpdateTime = now

	// 执行插入
	result, err := m.db.Exec(
		"INSERT INTO "+m.db.TableName("language_text")+" (languageid, `key`, value, module, createtime, updatetime) VALUES (?, ?, ?, ?, ?, ?)",
		text.LanguageID, text.Key, text.Value, text.Module, text.CreateTime, text.UpdateTime,
	)
	if err != nil {
		logger.Error("创建语言文本失败", "error", err)
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

// Update 更新语言文本
func (m *LanguageTextModel) Update(text *LanguageText) error {
	// 设置更新时间
	text.UpdateTime = time.Now()

	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("language_text")+" SET languageid = ?, `key` = ?, value = ?, module = ?, updatetime = ? WHERE id = ?",
		text.LanguageID, text.Key, text.Value, text.Module, text.UpdateTime, text.ID,
	)
	if err != nil {
		logger.Error("更新语言文本失败", "error", err)
		return err
	}

	return nil
}

// Delete 删除语言文本
func (m *LanguageTextModel) Delete(id int64) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("language_text")+" WHERE id = ?", id)
	if err != nil {
		logger.Error("删除语言文本失败", "error", err)
		return err
	}

	return nil
}

// DeleteByLanguageID 根据语言ID删除语言文本
func (m *LanguageTextModel) DeleteByLanguageID(languageID int64) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("language_text")+" WHERE languageid = ?", languageID)
	if err != nil {
		logger.Error("删除语言文本失败", "languageid", languageID, "error", err)
		return err
	}

	return nil
}

// DeleteByKey 根据键删除语言文本
func (m *LanguageTextModel) DeleteByKey(key string) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("language_text")+" WHERE `key` = ?", key)
	if err != nil {
		logger.Error("删除语言文本失败", "key", key, "error", err)
		return err
	}

	return nil
}

// ImportTexts 导入语言文本
func (m *LanguageTextModel) ImportTexts(languageID int64, texts map[string]string, module string) error {
	// 开始事务
	tx, err := m.db.Begin()
	if err != nil {
		logger.Error("开始事务失败", "error", err)
		return err
	}
	defer tx.Rollback()

	// 设置时间
	now := time.Now()

	// 导入语言文本
	for key, value := range texts {
		// 检查是否存在
		qb := database.NewQueryBuilder(m.db, "language_text")
		qb.Where("languageid = ?", languageID)
		qb.Where("key = ?", key)
		qb.Where("module = ?", module)
		result, err := qb.First()
		if err != nil {
			logger.Error("检查语言文本是否存在失败", "error", err)
			return err
		}

		if result == nil {
			// 不存在，创建
			_, err = tx.Exec(
				"INSERT INTO "+m.db.TableName("language_text")+" (languageid, `key`, value, module, createtime, updatetime) VALUES (?, ?, ?, ?, ?, ?)",
				languageID, key, value, module, now, now,
			)
			if err != nil {
				logger.Error("创建语言文本失败", "error", err)
				return err
			}
		} else {
			// 存在，更新
			id, _ := result["id"].(int64)
			_, err = tx.Exec(
				"UPDATE "+m.db.TableName("language_text")+" SET value = ?, updatetime = ? WHERE id = ?",
				value, now, id,
			)
			if err != nil {
				logger.Error("更新语言文本失败", "error", err)
				return err
			}
		}
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		logger.Error("提交事务失败", "error", err)
		return err
	}

	return nil
}

// ExportTexts 导出语言文本
func (m *LanguageTextModel) ExportTexts(languageID int64, module string) (map[string]string, error) {
	// 获取语言文本
	texts, err := m.GetByLanguageID(languageID, module)
	if err != nil {
		return nil, err
	}

	// 转换为map
	result := make(map[string]string)
	for _, text := range texts {
		result[text.Key] = text.Value
	}

	return result, nil
}
