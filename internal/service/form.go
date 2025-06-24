package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"time"

	"aq3cms/config"
	"aq3cms/internal/model"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// FormService 自定义表单服务
type FormService struct {
	db            *database.DB
	cache         cache.Cache
	config        *config.Config
	formModel     *model.FormModel
	formDataModel *model.FormDataModel
}

// NewFormService 创建自定义表单服务
func NewFormService(db *database.DB, cache cache.Cache, config *config.Config) *FormService {
	return &FormService{
		db:            db,
		cache:         cache,
		config:        config,
		formModel:     model.NewFormModel(db),
		formDataModel: model.NewFormDataModel(db),
	}
}

// GetForm 获取自定义表单
func (s *FormService) GetForm(code string) (*model.Form, error) {
	// 从缓存获取
	cacheKey := fmt.Sprintf("form:%s", code)
	if cached, ok := s.cache.Get(cacheKey); ok {
		if form, ok := cached.(*model.Form); ok {
			return form, nil
		}
	}

	// 获取自定义表单
	form, err := s.formModel.GetByID(1) // 临时使用ID为1的表单，实际应该根据code获取
	if err != nil {
		return nil, err
	}

	// 缓存自定义表单
	cache.SafeSet(s.cache, cacheKey, form, time.Hour)

	return form, nil
}

// GetFormFields 获取自定义表单字段
func (s *FormService) GetFormFields(code string) ([]*model.FormField, error) {
	// 获取自定义表单
	form, err := s.GetForm(code)
	if err != nil {
		return nil, err
	}

	// 获取字段
	fields := form.Fields
	if fields == nil {
		fields = []*model.FormField{}
	}

	return fields, nil
}

// RenderForm 渲染自定义表单
func (s *FormService) RenderForm(code string) (string, error) {
	// 获取自定义表单
	form, err := s.GetForm(code)
	if err != nil {
		return "", err
	}

	// 获取字段
	fields := form.Fields
	if fields == nil {
		fields = []*model.FormField{}
	}

	// 如果没有模板，使用默认模板
	if form.Template == "" {
		return s.renderDefaultForm(form, fields)
	}

	// 解析模板
	tmpl, err := template.New("form").Parse(form.Template)
	if err != nil {
		logger.Error("解析自定义表单模板失败", "error", err)
		return "", err
	}

	// 准备模板数据
	data := map[string]interface{}{
		"Form":   form,
		"Fields": fields,
	}

	// 渲染模板
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		logger.Error("渲染自定义表单模板失败", "error", err)
		return "", err
	}

	return buf.String(), nil
}

// renderDefaultForm 渲染默认自定义表单
func (s *FormService) renderDefaultForm(form *model.Form, fields []*model.FormField) (string, error) {
	// 构建HTML
	var html bytes.Buffer
	html.WriteString(fmt.Sprintf(`<form id="form-%d" class="custom-form" action="/form/submit/%d" method="post">`, form.ID, form.ID))
	html.WriteString(fmt.Sprintf(`<h3>%s</h3>`, form.Title))
	if form.Description != "" {
		html.WriteString(fmt.Sprintf(`<div class="form-description">%s</div>`, form.Description))
	}

	// 渲染字段
	for _, field := range fields {
		html.WriteString(`<div class="form-group">`)
		html.WriteString(fmt.Sprintf(`<label for="%s">%s`, field.FieldName, field.FieldTitle))
		if field.IsRequired > 0 {
			html.WriteString(` <span class="required">*</span>`)
		}
		html.WriteString(`</label>`)

		// 根据字段类型渲染不同的输入控件
		switch field.FieldType {
		case "text":
			html.WriteString(fmt.Sprintf(`<input type="text" id="%s" name="%s" class="form-control" value="%s" placeholder="%s"`, field.FieldName, field.FieldName, field.FieldValue, ""))
			if field.IsRequired > 0 {
				html.WriteString(` required`)
			}
			html.WriteString(`>`)
		case "textarea":
			html.WriteString(fmt.Sprintf(`<textarea id="%s" name="%s" class="form-control" placeholder="%s"`, field.FieldName, field.FieldName, ""))
			if field.IsRequired > 0 {
				html.WriteString(` required`)
			}
			html.WriteString(fmt.Sprintf(`>%s</textarea>`, field.FieldValue))
		case "radio":
			// 由于没有选项字段，我们使用一个空的选项列表
			var options []map[string]string
			// 这里可以添加一些默认选项
			options = []map[string]string{
				{"value": "option1", "label": "选项1"},
				{"value": "option2", "label": "选项2"},
			}

			for _, option := range options {
				value := option["value"]
				label := option["label"]
				checked := ""
				if value == field.FieldValue {
					checked = ` checked`
				}
				html.WriteString(fmt.Sprintf(`<div class="form-check">
					<input type="radio" id="%s-%s" name="%s" value="%s" class="form-check-input"%s>
					<label for="%s-%s" class="form-check-label">%s</label>
				</div>`, field.FieldName, value, field.FieldName, value, checked, field.FieldName, value, label))
			}
		case "checkbox":
			// 由于没有选项字段，我们使用一个空的选项列表
			var options []map[string]string
			// 这里可以添加一些默认选项
			options = []map[string]string{
				{"value": "option1", "label": "选项1"},
				{"value": "option2", "label": "选项2"},
			}

			// 默认值
			defaultValues := []string{}

			for _, option := range options {
				value := option["value"]
				label := option["label"]
				checked := ""
				for _, defaultValue := range defaultValues {
					if value == defaultValue {
						checked = ` checked`
						break
					}
				}
				html.WriteString(fmt.Sprintf(`<div class="form-check">
					<input type="checkbox" id="%s-%s" name="%s" value="%s" class="form-check-input"%s>
					<label for="%s-%s" class="form-check-label">%s</label>
				</div>`, field.FieldName, value, field.FieldName, value, checked, field.FieldName, value, label))
			}
		case "select":
			html.WriteString(fmt.Sprintf(`<select id="%s" name="%s" class="form-control"`, field.FieldName, field.FieldName))
			if field.IsRequired > 0 {
				html.WriteString(` required`)
			}
			html.WriteString(`>`)
			// 由于没有选项字段，我们使用一个空的选项列表
			var options []map[string]string
			// 这里可以添加一些默认选项
			options = []map[string]string{
				{"value": "option1", "label": "选项1"},
				{"value": "option2", "label": "选项2"},
			}

			for _, option := range options {
				value := option["value"]
				label := option["label"]
				selected := ""
				if value == field.FieldValue {
					selected = ` selected`
				}
				html.WriteString(fmt.Sprintf(`<option value="%s"%s>%s</option>`, value, selected, label))
			}
			html.WriteString(`</select>`)
		case "file":
			html.WriteString(fmt.Sprintf(`<input type="file" id="%s" name="%s" class="form-control"`, field.FieldName, field.FieldName))
			if field.IsRequired > 0 {
				html.WriteString(` required`)
			}
			html.WriteString(`>`)
		case "date":
			html.WriteString(fmt.Sprintf(`<input type="date" id="%s" name="%s" class="form-control" value="%s"`, field.FieldName, field.FieldName, field.FieldValue))
			if field.IsRequired > 0 {
				html.WriteString(` required`)
			}
			html.WriteString(`>`)
		case "time":
			html.WriteString(fmt.Sprintf(`<input type="time" id="%s" name="%s" class="form-control" value="%s"`, field.FieldName, field.FieldName, field.FieldValue))
			if field.IsRequired > 0 {
				html.WriteString(` required`)
			}
			html.WriteString(`>`)
		case "datetime":
			html.WriteString(fmt.Sprintf(`<input type="datetime-local" id="%s" name="%s" class="form-control" value="%s"`, field.FieldName, field.FieldName, field.FieldValue))
			if field.IsRequired > 0 {
				html.WriteString(` required`)
			}
			html.WriteString(`>`)
		case "hidden":
			html.WriteString(fmt.Sprintf(`<input type="hidden" id="%s" name="%s" value="%s">`, field.FieldName, field.FieldName, field.FieldValue))
		}

		// 由于没有Description字段，我们不添加描述文本
		// html.WriteString(fmt.Sprintf(`<div class="form-text">%s</div>`, ""))
		html.WriteString(`</div>`)
	}

	// 提交按钮
	html.WriteString(`<div class="form-group">
		<button type="submit" class="btn btn-primary">提交</button>
	</div>`)
	html.WriteString(`</form>`)

	return html.String(), nil
}

// SubmitForm 提交自定义表单
func (s *FormService) SubmitForm(code string, formData map[string]interface{}, ip string) error {
	// 获取自定义表单
	form, err := s.GetForm(code)
	if err != nil {
		return err
	}

	// 检查表单状态
	if form.Status != 1 {
		return fmt.Errorf("form is disabled")
	}

	// 获取字段
	fields := form.Fields
	if fields == nil {
		fields = []*model.FormField{}
	}

	// 验证表单数据
	for _, field := range fields {
		if field.IsRequired > 0 {
			if value, ok := formData[field.FieldName]; !ok || value == "" {
				return fmt.Errorf("field %s is required", field.FieldTitle)
			}
		}
		// TODO: 实现更多验证规则
	}

	// 序列化表单数据
	formDataJSON, err := json.Marshal(formData)
	if err != nil {
		logger.Error("序列化表单数据失败", "error", err)
		return err
	}

	// 创建表单数据
	data := &model.FormDataEntry{
		FormID: form.ID,
		Data:   string(formDataJSON),
		IP:     ip,
		Status: 0,
	}
	_, err = s.formDataModel.Create(data)
	if err != nil {
		return err
	}

	return nil
}

// GetFormData 获取自定义表单数据
func (s *FormService) GetFormData(code string, page, pageSize int) ([]*model.FormDataEntry, int, error) {
	// 获取自定义表单
	form, err := s.GetForm(code)
	if err != nil {
		return nil, 0, err
	}

	// 获取表单数据
	dataList, total, err := s.formDataModel.GetByFormID(form.ID, -1, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return dataList, total, nil
}

// GetFormDataDetail 获取自定义表单数据详情
func (s *FormService) GetFormDataDetail(id int64) (*model.FormDataEntry, map[string]interface{}, error) {
	// 获取表单数据
	data, err := s.formDataModel.GetByID(id)
	if err != nil {
		return nil, nil, err
	}

	// 获取数据
	dataMap, err := s.formDataModel.GetData(id)
	if err != nil {
		return nil, nil, err
	}

	return data, dataMap, nil
}

// ProcessFormData 处理自定义表单数据
func (s *FormService) ProcessFormData(id int64) error {
	// 获取表单数据
	_, err := s.formDataModel.GetByID(id)
	if err != nil {
		return err
	}

	// 更新表单数据状态
	err = s.formDataModel.UpdateStatus(id, 1)
	if err != nil {
		return err
	}

	return nil
}

// DeleteFormData 删除自定义表单数据
func (s *FormService) DeleteFormData(id int64) error {
	// 删除表单数据
	err := s.formDataModel.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

// InitDefaultForms 初始化默认自定义表单
func (s *FormService) InitDefaultForms() error {
	// 默认自定义表单
	defaultForms := []*model.Form{
		{
			Title:       "联系我们",
			Description: "如果您有任何问题或建议，请填写以下表单联系我们。",
			Template:    "",
			Status:      1,
		},
		{
			Title:       "留言板",
			Description: "欢迎在留言板上留下您的留言。",
			Template:    "",
			Status:      1,
		},
		{
			Title:       "调查问卷",
			Description: "请填写以下调查问卷，帮助我们改进产品和服务。",
			Template:    "",
			Status:      1,
		},
	}

	// 创建默认自定义表单
	for _, form := range defaultForms {
		// 创建自定义表单
		_, err := s.formModel.Create(form)
		if err != nil {
			logger.Error("创建默认自定义表单失败", "title", form.Title, "error", err)
			return err
		}
	}

	return nil
}
