package model

import (
	"fmt"
	"time"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
	"aq3cms/pkg/security"
)

// Form 表单模型
type Form struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`       // 表单标题
	Description string    `json:"description"` // 表单描述
	Template    string    `json:"template"`    // 表单模板
	Status      int       `json:"status"`      // 状态
	AddTime     time.Time `json:"addtime"`     // 添加时间
	Fields      []*FormField `json:"fields"`   // 表单字段
	SuccessURL  string    `json:"successurl"`  // 成功跳转URL
}

// FormField 表单字段
type FormField struct {
	ID          int64  `json:"id"`
	FormID      int64  `json:"formid"`      // 表单ID
	FieldName   string `json:"fieldname"`   // 字段名称
	FieldType   string `json:"fieldtype"`   // 字段类型
	FieldTitle  string `json:"fieldtitle"`  // 字段标题
	FieldValue  string `json:"fieldvalue"`  // 字段默认值
	IsRequired  int    `json:"isrequired"`  // 是否必填
	SortRank    int    `json:"sortrank"`    // 排序
}

// FormData 表单数据
type FormData struct {
	ID          int64     `json:"id"`
	FormID      int64     `json:"formid"`      // 表单ID
	IP          string    `json:"ip"`          // IP地址
	AddTime     time.Time `json:"addtime"`     // 添加时间
	Status      int       `json:"status"`      // 状态
	Data        map[string]string `json:"data"` // 表单数据
}

// FormModel 表单模型操作
type FormModel struct {
	db *database.DB
}

// NewFormModel 创建表单模型
func NewFormModel(db *database.DB) *FormModel {
	return &FormModel{
		db: db,
	}
}

// GetByID 根据ID获取表单
func (m *FormModel) GetByID(id int64) (*Form, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "diyform")
	qb.Select("*")
	qb.Where("id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("查询表单失败", "id", id, "error", err)
		return nil, err
	}

	if result == nil {
		return nil, fmt.Errorf("表单不存在")
	}

	// 转换为表单对象
	form := &Form{}
	form.ID, _ = result["id"].(int64)
	form.Title, _ = result["title"].(string)
	form.Description, _ = result["description"].(string)
	form.Template, _ = result["template"].(string)

	// 处理整数字段
	if status, ok := result["status"].(int64); ok {
		form.Status = int(status)
	}

	// 处理日期
	if addtime, ok := result["addtime"].(time.Time); ok {
		form.AddTime = addtime
	}

	// 获取表单字段
	fields, err := m.GetFormFields(form.ID)
	if err != nil {
		logger.Error("获取表单字段失败", "formid", form.ID, "error", err)
	}
	form.Fields = fields

	return form, nil
}

// GetList 获取表单列表
func (m *FormModel) GetList(page, pageSize int) ([]*Form, int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "diyform")
	qb.Select("*")

	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("查询表单总数失败", "error", err)
		return nil, 0, err
	}

	// 设置分页
	offset := (page - 1) * pageSize
	qb.OrderBy("id DESC")
	qb.Limit(pageSize, offset)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("查询表单列表失败", "error", err)
		return nil, 0, err
	}

	// 转换为表单对象
	forms := make([]*Form, 0, len(results))
	for _, result := range results {
		form := &Form{}
		form.ID, _ = result["id"].(int64)
		form.Title, _ = result["title"].(string)
		form.Description, _ = result["description"].(string)
		form.Template, _ = result["template"].(string)

		// 处理整数字段
		if status, ok := result["status"].(int64); ok {
			form.Status = int(status)
		}

		// 处理日期
		if addtime, ok := result["addtime"].(time.Time); ok {
			form.AddTime = addtime
		}

		forms = append(forms, form)
	}

	return forms, total, nil
}

// GetFormFields 获取表单字段
func (m *FormModel) GetFormFields(formID int64) ([]*FormField, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "diyform_field")
	qb.Select("*")
	qb.Where("formid = ?", formID)
	qb.OrderBy("sortrank ASC")

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("查询表单字段失败", "formid", formID, "error", err)
		return nil, err
	}

	// 转换为表单字段对象
	fields := make([]*FormField, 0, len(results))
	for _, result := range results {
		field := &FormField{}
		field.ID, _ = result["id"].(int64)
		field.FormID, _ = result["formid"].(int64)
		field.FieldName, _ = result["fieldname"].(string)
		field.FieldType, _ = result["fieldtype"].(string)
		field.FieldTitle, _ = result["fieldtitle"].(string)
		field.FieldValue, _ = result["fieldvalue"].(string)

		// 处理整数字段
		if isrequired, ok := result["isrequired"].(int64); ok {
			field.IsRequired = int(isrequired)
		}
		if sortrank, ok := result["sortrank"].(int64); ok {
			field.SortRank = int(sortrank)
		}

		fields = append(fields, field)
	}

	return fields, nil
}

// GetFormData 获取表单数据
func (m *FormModel) GetFormData(formID int64, page, pageSize int) ([]*FormData, int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "diyform_data")
	qb.Select("*")
	qb.Where("formid = ?", formID)

	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("查询表单数据总数失败", "formid", formID, "error", err)
		return nil, 0, err
	}

	// 设置分页
	offset := (page - 1) * pageSize
	qb.OrderBy("id DESC")
	qb.Limit(pageSize, offset)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("查询表单数据失败", "formid", formID, "error", err)
		return nil, 0, err
	}

	// 获取表单字段
	fields, err := m.GetFormFields(formID)
	if err != nil {
		logger.Error("获取表单字段失败", "formid", formID, "error", err)
		return nil, 0, err
	}

	// 转换为表单数据对象
	formDataList := make([]*FormData, 0, len(results))
	for _, result := range results {
		formData := &FormData{}
		formData.ID, _ = result["id"].(int64)
		formData.FormID, _ = result["formid"].(int64)
		formData.IP, _ = result["ip"].(string)

		// 处理日期
		if addtime, ok := result["addtime"].(time.Time); ok {
			formData.AddTime = addtime
		}

		// 处理整数字段
		if status, ok := result["status"].(int64); ok {
			formData.Status = int(status)
		}

		// 获取表单数据
		formData.Data = make(map[string]string)
		for _, field := range fields {
			if value, ok := result[field.FieldName].(string); ok {
				formData.Data[field.FieldName] = value
			}
		}

		formDataList = append(formDataList, formData)
	}

	return formDataList, total, nil
}

// Create 创建表单
func (m *FormModel) Create(form *Form) (int64, error) {
	// 开始事务
	tx, err := m.db.Begin()
	if err != nil {
		logger.Error("开始事务失败", "error", err)
		return 0, err
	}
	defer tx.Rollback()

	// 设置默认值
	if form.AddTime.IsZero() {
		form.AddTime = time.Now()
	}

	// 准备插入数据

	// 执行插入
	result, err := tx.Exec(
		"INSERT INTO "+m.db.TableName("diyform")+" (title, description, template, status, addtime) VALUES (?, ?, ?, ?, ?)",
		form.Title, form.Description, form.Template, form.Status, form.AddTime,
	)
	if err != nil {
		logger.Error("创建表单失败", "error", err)
		return 0, err
	}

	// 获取插入ID
	formID, err := result.LastInsertId()
	if err != nil {
		logger.Error("获取插入ID失败", "error", err)
		return 0, err
	}

	// 插入表单字段
	for _, field := range form.Fields {
		field.FormID = formID

		_, err = tx.Exec(
			"INSERT INTO "+m.db.TableName("diyform_field")+" (formid, fieldname, fieldtype, fieldtitle, fieldvalue, isrequired, sortrank) VALUES (?, ?, ?, ?, ?, ?, ?)",
			field.FormID, field.FieldName, field.FieldType, field.FieldTitle, field.FieldValue, field.IsRequired, field.SortRank,
		)
		if err != nil {
			logger.Error("创建表单字段失败", "error", err)
			return 0, err
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		logger.Error("提交事务失败", "error", err)
		return 0, err
	}

	return formID, nil
}

// Update 更新表单
func (m *FormModel) Update(form *Form) error {
	// 开始事务
	tx, err := m.db.Begin()
	if err != nil {
		logger.Error("开始事务失败", "error", err)
		return err
	}
	defer tx.Rollback()

	// 更新表单
	_, err = tx.Exec(
		"UPDATE "+m.db.TableName("diyform")+" SET title = ?, description = ?, template = ?, status = ? WHERE id = ?",
		form.Title, form.Description, form.Template, form.Status, form.ID,
	)
	if err != nil {
		logger.Error("更新表单失败", "id", form.ID, "error", err)
		return err
	}

	// 删除旧的表单字段
	_, err = tx.Exec(
		"DELETE FROM "+m.db.TableName("diyform_field")+" WHERE formid = ?",
		form.ID,
	)
	if err != nil {
		logger.Error("删除表单字段失败", "formid", form.ID, "error", err)
		return err
	}

	// 插入新的表单字段
	for _, field := range form.Fields {
		field.FormID = form.ID

		_, err = tx.Exec(
			"INSERT INTO "+m.db.TableName("diyform_field")+" (formid, fieldname, fieldtype, fieldtitle, fieldvalue, isrequired, sortrank) VALUES (?, ?, ?, ?, ?, ?, ?)",
			field.FormID, field.FieldName, field.FieldType, field.FieldTitle, field.FieldValue, field.IsRequired, field.SortRank,
		)
		if err != nil {
			logger.Error("创建表单字段失败", "error", err)
			return err
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		logger.Error("提交事务失败", "error", err)
		return err
	}

	return nil
}

// Delete 删除表单
func (m *FormModel) Delete(id int64) error {
	// 开始事务
	tx, err := m.db.Begin()
	if err != nil {
		logger.Error("开始事务失败", "error", err)
		return err
	}
	defer tx.Rollback()

	// 删除表单
	_, err = tx.Exec(
		"DELETE FROM "+m.db.TableName("diyform")+" WHERE id = ?",
		id,
	)
	if err != nil {
		logger.Error("删除表单失败", "id", id, "error", err)
		return err
	}

	// 删除表单字段
	_, err = tx.Exec(
		"DELETE FROM "+m.db.TableName("diyform_field")+" WHERE formid = ?",
		id,
	)
	if err != nil {
		logger.Error("删除表单字段失败", "formid", id, "error", err)
		return err
	}

	// 删除表单数据
	_, err = tx.Exec(
		"DELETE FROM "+m.db.TableName("diyform_data")+" WHERE formid = ?",
		id,
	)
	if err != nil {
		logger.Error("删除表单数据失败", "formid", id, "error", err)
		return err
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		logger.Error("提交事务失败", "error", err)
		return err
	}

	return nil
}

// SubmitForm 提交表单
func (m *FormModel) SubmitForm(formID int64, data map[string]string, ip string) (int64, error) {
	// 获取表单
	form, err := m.GetByID(formID)
	if err != nil {
		logger.Error("获取表单失败", "formid", formID, "error", err)
		return 0, err
	}

	// 检查表单状态
	if form.Status != 1 {
		return 0, fmt.Errorf("表单已关闭")
	}

	// 开始事务
	tx, err := m.db.Begin()
	if err != nil {
		logger.Error("开始事务失败", "error", err)
		return 0, err
	}
	defer tx.Rollback()

	// 构建SQL
	sql := "INSERT INTO " + m.db.TableName("diyform_data") + " (formid, ip, addtime, status"
	values := []interface{}{formID, ip, time.Now(), 0}

	// 添加字段
	for _, field := range form.Fields {
		sql += ", " + field.FieldName

		// 获取字段值
		value, ok := data[field.FieldName]
		if !ok {
			value = ""
		}

		// 清理字段值，防止XSS攻击
		value = security.CleanHTML(value)

		values = append(values, value)
	}

	// 完成SQL
	sql += ") VALUES (?, ?, ?, ?"
	for range form.Fields {
		sql += ", ?"
	}
	sql += ")"

	// 执行插入
	result, err := tx.Exec(sql, values...)
	if err != nil {
		logger.Error("提交表单失败", "formid", formID, "error", err)
		return 0, err
	}

	// 获取插入ID
	dataID, err := result.LastInsertId()
	if err != nil {
		logger.Error("获取插入ID失败", "error", err)
		return 0, err
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		logger.Error("提交事务失败", "error", err)
		return 0, err
	}

	return dataID, nil
}
