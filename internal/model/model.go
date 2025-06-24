package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// ContentModel 内容模型
type ContentModel struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`        // 模型名称
	TableName   string    `json:"tablename"`   // 表名
	Description string    `json:"description"` // 描述
	State       int       `json:"state"`       // 状态：0禁用，1启用
	Fields      string    `json:"fields"`      // 字段定义，JSON格式
	CreateTime  time.Time `json:"createtime"`  // 创建时间
	UpdateTime  time.Time `json:"updatetime"`  // 更新时间
}

// Field 字段定义
type Field struct {
	Name        string `json:"name"`        // 字段名
	Title       string `json:"title"`       // 字段标题
	Type        string `json:"type"`        // 字段类型
	Length      int    `json:"length"`      // 字段长度
	Default     string `json:"default"`     // 默认值
	Description string `json:"description"` // 描述
	Required    bool   `json:"required"`    // 是否必填
	Options     string `json:"options"`     // 选项，JSON格式
}

// ContentModelModel 内容模型模型
type ContentModelModel struct {
	db *database.DB
}

// NewContentModelModel 创建内容模型模型
func NewContentModelModel(db *database.DB) *ContentModelModel {
	return &ContentModelModel{
		db: db,
	}
}

// GetByID 根据ID获取内容模型
func (m *ContentModelModel) GetByID(id int64) (*ContentModel, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "content_model")
	qb.Where("id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取内容模型失败", "id", id, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("content model not found: %d", id)
	}

	// 转换为内容模型
	model := &ContentModel{}
	model.ID, _ = result["id"].(int64)
	model.Name, _ = result["name"].(string)
	model.TableName, _ = result["tablename"].(string)
	model.Description, _ = result["description"].(string)
	model.State, _ = result["state"].(int)
	model.Fields, _ = result["fields"].(string)
	model.CreateTime, _ = result["createtime"].(time.Time)
	model.UpdateTime, _ = result["updatetime"].(time.Time)

	return model, nil
}

// GetByName 根据名称获取内容模型
func (m *ContentModelModel) GetByName(name string) (*ContentModel, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "content_model")
	qb.Where("name = ?", name)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取内容模型失败", "name", name, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("content model not found: %s", name)
	}

	// 转换为内容模型
	model := &ContentModel{}
	model.ID, _ = result["id"].(int64)
	model.Name, _ = result["name"].(string)
	model.TableName, _ = result["tablename"].(string)
	model.Description, _ = result["description"].(string)
	model.State, _ = result["state"].(int)
	model.Fields, _ = result["fields"].(string)
	model.CreateTime, _ = result["createtime"].(time.Time)
	model.UpdateTime, _ = result["updatetime"].(time.Time)

	return model, nil
}

// GetAll 获取所有内容模型
func (m *ContentModelModel) GetAll() ([]*ContentModel, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "content_model")
	qb.OrderBy("id ASC")

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取所有内容模型失败", "error", err)
		return nil, err
	}

	// 转换为内容模型列表
	models := make([]*ContentModel, 0, len(results))
	for _, result := range results {
		model := &ContentModel{}
		model.ID, _ = result["id"].(int64)
		model.Name, _ = result["name"].(string)
		model.TableName, _ = result["tablename"].(string)
		model.Description, _ = result["description"].(string)
		model.State, _ = result["state"].(int)
		model.Fields, _ = result["fields"].(string)
		model.CreateTime, _ = result["createtime"].(time.Time)
		model.UpdateTime, _ = result["updatetime"].(time.Time)
		models = append(models, model)
	}

	return models, nil
}

// Create 创建内容模型
func (m *ContentModelModel) Create(model *ContentModel) (int64, error) {
	// 检查表名是否已存在
	exists, err := m.TableExists(model.TableName)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, fmt.Errorf("table already exists: %s", model.TableName)
	}

	// 设置创建时间和更新时间
	now := time.Now()
	model.CreateTime = now
	model.UpdateTime = now

	// 执行插入
	result, err := m.db.Exec(
		"INSERT INTO "+m.db.TableName("content_model")+" (name, tablename, description, state, fields, createtime, updatetime) VALUES (?, ?, ?, ?, ?, ?, ?)",
		model.Name, model.TableName, model.Description, model.State, model.Fields, model.CreateTime, model.UpdateTime,
	)
	if err != nil {
		logger.Error("创建内容模型失败", "error", err)
		return 0, err
	}

	// 获取插入ID
	id, err := result.LastInsertId()
	if err != nil {
		logger.Error("获取插入ID失败", "error", err)
		return 0, err
	}

	// 创建数据表
	err = m.CreateTable(model)
	if err != nil {
		logger.Error("创建数据表失败", "error", err)
		// 删除内容模型
		m.db.Exec("DELETE FROM "+m.db.TableName("content_model")+" WHERE id = ?", id)
		return 0, err
	}

	return id, nil
}

// Update 更新内容模型
func (m *ContentModelModel) Update(model *ContentModel) error {
	// 获取原内容模型
	oldModel, err := m.GetByID(model.ID)
	if err != nil {
		return err
	}

	// 设置更新时间
	model.UpdateTime = time.Now()

	// 执行更新
	_, err = m.db.Exec(
		"UPDATE "+m.db.TableName("content_model")+" SET name = ?, description = ?, state = ?, fields = ?, updatetime = ? WHERE id = ?",
		model.Name, model.Description, model.State, model.Fields, model.UpdateTime, model.ID,
	)
	if err != nil {
		logger.Error("更新内容模型失败", "error", err)
		return err
	}

	// 更新数据表
	err = m.UpdateTable(oldModel, model)
	if err != nil {
		logger.Error("更新数据表失败", "error", err)
		return err
	}

	return nil
}

// Delete 删除内容模型
func (m *ContentModelModel) Delete(id int64) error {
	// 获取内容模型
	model, err := m.GetByID(id)
	if err != nil {
		return err
	}

	// 删除数据表
	err = m.DropTable(model)
	if err != nil {
		logger.Error("删除数据表失败", "error", err)
		return err
	}

	// 执行删除
	_, err = m.db.Exec("DELETE FROM "+m.db.TableName("content_model")+" WHERE id = ?", id)
	if err != nil {
		logger.Error("删除内容模型失败", "error", err)
		return err
	}

	return nil
}

// TableExists 检查表是否存在
func (m *ContentModelModel) TableExists(tableName string) (bool, error) {
	// 构建查询
	var count int
	err := m.db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?", m.db.TableName(tableName)).Scan(&count)
	if err != nil {
		logger.Error("检查表是否存在失败", "tableName", tableName, "error", err)
		return false, err
	}

	return count > 0, nil
}

// CreateTable 创建数据表
func (m *ContentModelModel) CreateTable(model *ContentModel) error {
	// 解析字段定义
	var fields []Field
	err := json.Unmarshal([]byte(model.Fields), &fields)
	if err != nil {
		logger.Error("解析字段定义失败", "error", err)
		return err
	}

	// 构建SQL
	var sql strings.Builder
	sql.WriteString("CREATE TABLE " + m.db.TableName(model.TableName) + " (\n")
	sql.WriteString("  id BIGINT(20) NOT NULL AUTO_INCREMENT,\n")
	sql.WriteString("  aid BIGINT(20) NOT NULL DEFAULT 0,\n") // 关联文章ID

	// 添加字段
	for _, field := range fields {
		sql.WriteString("  " + field.Name + " ")
		switch field.Type {
		case "varchar":
			sql.WriteString("VARCHAR(" + fmt.Sprintf("%d", field.Length) + ")")
		case "int":
			sql.WriteString("INT(" + fmt.Sprintf("%d", field.Length) + ")")
		case "float":
			sql.WriteString("FLOAT")
		case "text":
			sql.WriteString("TEXT")
		case "longtext":
			sql.WriteString("LONGTEXT")
		case "datetime":
			sql.WriteString("DATETIME")
		case "date":
			sql.WriteString("DATE")
		default:
			sql.WriteString("VARCHAR(255)")
		}
		if field.Default != "" {
			sql.WriteString(" DEFAULT '" + field.Default + "'")
		}
		if field.Required {
			sql.WriteString(" NOT NULL")
		} else {
			sql.WriteString(" NULL")
		}
		sql.WriteString(",\n")
	}

	// 添加主键和索引
	sql.WriteString("  PRIMARY KEY (id),\n")
	sql.WriteString("  KEY idx_aid (aid)\n")
	sql.WriteString(") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")

	// 执行SQL
	_, err = m.db.Exec(sql.String())
	if err != nil {
		logger.Error("创建数据表失败", "error", err)
		return err
	}

	return nil
}

// UpdateTable 更新数据表
func (m *ContentModelModel) UpdateTable(oldModel, newModel *ContentModel) error {
	// 解析旧字段定义
	var oldFields []Field
	err := json.Unmarshal([]byte(oldModel.Fields), &oldFields)
	if err != nil {
		logger.Error("解析旧字段定义失败", "error", err)
		return err
	}

	// 解析新字段定义
	var newFields []Field
	err = json.Unmarshal([]byte(newModel.Fields), &newFields)
	if err != nil {
		logger.Error("解析新字段定义失败", "error", err)
		return err
	}

	// 创建字段映射
	oldFieldMap := make(map[string]Field)
	for _, field := range oldFields {
		oldFieldMap[field.Name] = field
	}

	newFieldMap := make(map[string]Field)
	for _, field := range newFields {
		newFieldMap[field.Name] = field
	}

	// 查找需要添加的字段
	for name, field := range newFieldMap {
		if _, ok := oldFieldMap[name]; !ok {
			// 添加字段
			var sql strings.Builder
			sql.WriteString("ALTER TABLE " + m.db.TableName(oldModel.TableName) + " ADD COLUMN " + field.Name + " ")
			switch field.Type {
			case "varchar":
				sql.WriteString("VARCHAR(" + fmt.Sprintf("%d", field.Length) + ")")
			case "int":
				sql.WriteString("INT(" + fmt.Sprintf("%d", field.Length) + ")")
			case "float":
				sql.WriteString("FLOAT")
			case "text":
				sql.WriteString("TEXT")
			case "longtext":
				sql.WriteString("LONGTEXT")
			case "datetime":
				sql.WriteString("DATETIME")
			case "date":
				sql.WriteString("DATE")
			default:
				sql.WriteString("VARCHAR(255)")
			}
			if field.Default != "" {
				sql.WriteString(" DEFAULT '" + field.Default + "'")
			}
			if field.Required {
				sql.WriteString(" NOT NULL")
			} else {
				sql.WriteString(" NULL")
			}

			// 执行SQL
			_, err = m.db.Exec(sql.String())
			if err != nil {
				logger.Error("添加字段失败", "error", err)
				return err
			}
		}
	}

	// 查找需要修改的字段
	for name, oldField := range oldFieldMap {
		if newField, ok := newFieldMap[name]; ok {
			// 检查字段是否需要修改
			if oldField.Type != newField.Type || oldField.Length != newField.Length || oldField.Default != newField.Default || oldField.Required != newField.Required {
				// 修改字段
				var sql strings.Builder
				sql.WriteString("ALTER TABLE " + m.db.TableName(oldModel.TableName) + " MODIFY COLUMN " + newField.Name + " ")
				switch newField.Type {
				case "varchar":
					sql.WriteString("VARCHAR(" + fmt.Sprintf("%d", newField.Length) + ")")
				case "int":
					sql.WriteString("INT(" + fmt.Sprintf("%d", newField.Length) + ")")
				case "float":
					sql.WriteString("FLOAT")
				case "text":
					sql.WriteString("TEXT")
				case "longtext":
					sql.WriteString("LONGTEXT")
				case "datetime":
					sql.WriteString("DATETIME")
				case "date":
					sql.WriteString("DATE")
				default:
					sql.WriteString("VARCHAR(255)")
				}
				if newField.Default != "" {
					sql.WriteString(" DEFAULT '" + newField.Default + "'")
				}
				if newField.Required {
					sql.WriteString(" NOT NULL")
				} else {
					sql.WriteString(" NULL")
				}

				// 执行SQL
				_, err = m.db.Exec(sql.String())
				if err != nil {
					logger.Error("修改字段失败", "error", err)
					return err
				}
			}
		} else {
			// 删除字段
			sql := "ALTER TABLE " + m.db.TableName(oldModel.TableName) + " DROP COLUMN " + oldField.Name
			_, err = m.db.Exec(sql)
			if err != nil {
				logger.Error("删除字段失败", "error", err)
				return err
			}
		}
	}

	return nil
}

// DropTable 删除数据表
func (m *ContentModelModel) DropTable(model *ContentModel) error {
	// 执行SQL
	_, err := m.db.Exec("DROP TABLE IF EXISTS " + m.db.TableName(model.TableName))
	if err != nil {
		logger.Error("删除数据表失败", "error", err)
		return err
	}

	return nil
}

// GetContent 获取内容
func (m *ContentModelModel) GetContent(modelID int64, contentID int64) (map[string]interface{}, error) {
	// 获取内容模型
	model, err := m.GetByID(modelID)
	if err != nil {
		return nil, err
	}

	// 构建查询
	qb := database.NewQueryBuilder(m.db, model.TableName)
	qb.Where("aid = ?", contentID)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取内容失败", "error", err)
		return nil, err
	}

	return result, nil
}

// SaveContent 保存内容
func (m *ContentModelModel) SaveContent(modelID int64, contentID int64, data map[string]interface{}) error {
	// 获取内容模型
	model, err := m.GetByID(modelID)
	if err != nil {
		return err
	}

	// 检查内容是否存在
	exists, err := m.ContentExists(model.TableName, contentID)
	if err != nil {
		return err
	}

	// 设置文章ID
	data["aid"] = contentID

	if exists {
		// 更新内容
		qb := database.NewQueryBuilder(m.db, model.TableName)
		qb.Where("aid = ?", contentID)
		_, err = qb.Update(data)
		if err != nil {
			logger.Error("更新内容失败", "error", err)
			return err
		}
	} else {
		// 创建内容
		qb := database.NewQueryBuilder(m.db, model.TableName)
		_, err = qb.Insert(data)
		if err != nil {
			logger.Error("创建内容失败", "error", err)
			return err
		}
	}

	return nil
}

// DeleteContent 删除内容
func (m *ContentModelModel) DeleteContent(modelID int64, contentID int64) error {
	// 获取内容模型
	model, err := m.GetByID(modelID)
	if err != nil {
		return err
	}

	// 执行删除
	_, err = m.db.Exec("DELETE FROM "+m.db.TableName(model.TableName)+" WHERE aid = ?", contentID)
	if err != nil {
		logger.Error("删除内容失败", "error", err)
		return err
	}

	return nil
}

// ContentExists 检查内容是否存在
func (m *ContentModelModel) ContentExists(tableName string, contentID int64) (bool, error) {
	// 构建查询
	var count int
	err := m.db.QueryRow("SELECT COUNT(*) FROM "+m.db.TableName(tableName)+" WHERE aid = ?", contentID).Scan(&count)
	if err != nil {
		logger.Error("检查内容是否存在失败", "error", err)
		return false, err
	}

	return count > 0, nil
}
