package admin

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"aq3cms/config"
	"aq3cms/internal/middleware"
	"aq3cms/internal/model"
	"aq3cms/internal/service"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// CollectController 采集控制器
type CollectController struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	collectService  *service.CollectService
	collectRuleModel *model.CollectRuleModel
	collectItemModel *model.CollectItemModel
	categoryModel   *model.CategoryModel
	modelModel      *model.ContentModelModel
	templateService *service.TemplateService
}

// NewCollectController 创建采集控制器
func NewCollectController(db *database.DB, cache cache.Cache, config *config.Config) *CollectController {
	return &CollectController{
		db:              db,
		cache:           cache,
		config:          config,
		collectService:  service.NewCollectService(db, cache, config),
		collectRuleModel: model.NewCollectRuleModel(db),
		collectItemModel: model.NewCollectItemModel(db),
		categoryModel:   model.NewCategoryModel(db),
		modelModel:      model.NewContentModelModel(db),
		templateService: service.NewTemplateService(db, cache, config),
	}
}

// RuleList 规则列表
func (c *CollectController) RuleList(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取规则列表
	rules, err := c.collectRuleModel.GetAll()
	if err != nil {
		logger.Error("获取采集规则列表失败", "error", err)
		http.Error(w, "Failed to get collect rules", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Rules":       rules,
		"CurrentMenu": "collect",
		"PageTitle":   "采集规则管理",
	}

	// 渲染模板
	tplFile := "admin/collect_rule_list.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染采集规则列表模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// RuleAdd 添加规则页面
func (c *CollectController) RuleAdd(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取栏目列表
	categories, err := c.categoryModel.GetAll()
	if err != nil {
		logger.Error("获取栏目列表失败", "error", err)
		http.Error(w, "Failed to get categories", http.StatusInternalServerError)
		return
	}

	// 获取模型列表
	models, err := c.modelModel.GetAll()
	if err != nil {
		logger.Error("获取模型列表失败", "error", err)
		http.Error(w, "Failed to get models", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Categories":  categories,
		"Models":      models,
		"CurrentMenu": "collect",
		"PageTitle":   "添加采集规则",
	}

	// 渲染模板
	tplFile := "admin/collect_rule_add.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染添加采集规则模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// RuleDoAdd 处理添加规则
func (c *CollectController) RuleDoAdd(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	name := r.FormValue("name")
	sourceTypeStr := r.FormValue("sourcetype")
	sourceURL := r.FormValue("sourceurl")
	startPageStr := r.FormValue("startpage")
	endPageStr := r.FormValue("endpage")
	pageRule := r.FormValue("pagerule")
	listRule := r.FormValue("listrule")
	titleRule := r.FormValue("titlerule")
	contentRule := r.FormValue("contentrule")
	typeIDStr := r.FormValue("typeid")
	modelIDStr := r.FormValue("modelid")
	statusStr := r.FormValue("status")
	fieldRules := r.FormValue("fieldrules")

	// 验证必填字段
	if name == "" || sourceURL == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 解析数值
	sourceType, _ := strconv.Atoi(sourceTypeStr)
	startPage, _ := strconv.Atoi(startPageStr)
	endPage, _ := strconv.Atoi(endPageStr)
	typeID, _ := strconv.ParseInt(typeIDStr, 10, 64)
	modelID, _ := strconv.ParseInt(modelIDStr, 10, 64)
	status, _ := strconv.Atoi(statusStr)

	// 创建规则
	rule := &model.CollectRule{
		Name:        name,
		SourceType:  sourceType,
		SourceURL:   sourceURL,
		StartPage:   startPage,
		EndPage:     endPage,
		PageRule:    pageRule,
		ListRule:    listRule,
		TitleRule:   titleRule,
		ContentRule: contentRule,
		TypeID:      typeID,
		ModelID:     modelID,
		Status:      status,
		LastTime:    time.Time{},
		FieldRules:  fieldRules,
	}

	// 保存规则
	id, err := c.collectRuleModel.Create(rule)
	if err != nil {
		logger.Error("创建采集规则失败", "error", err)
		http.Error(w, "Failed to create collect rule", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "采集规则创建成功",
			"id":      id,
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/collect_rule_list", http.StatusFound)
	}
}

// RuleEdit 编辑规则页面
func (c *CollectController) RuleEdit(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取规则ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid rule ID", http.StatusBadRequest)
		return
	}

	// 获取规则
	rule, err := c.collectRuleModel.GetByID(id)
	if err != nil {
		logger.Error("获取采集规则失败", "id", id, "error", err)
		http.Error(w, "Rule not found", http.StatusNotFound)
		return
	}

	// 获取栏目列表
	categories, err := c.categoryModel.GetAll()
	if err != nil {
		logger.Error("获取栏目列表失败", "error", err)
		http.Error(w, "Failed to get categories", http.StatusInternalServerError)
		return
	}

	// 获取模型列表
	models, err := c.modelModel.GetAll()
	if err != nil {
		logger.Error("获取模型列表失败", "error", err)
		http.Error(w, "Failed to get models", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Rule":        rule,
		"Categories":  categories,
		"Models":      models,
		"CurrentMenu": "collect",
		"PageTitle":   "编辑采集规则",
	}

	// 渲染模板
	tplFile := "admin/collect_rule_edit.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染编辑采集规则模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// RuleDoEdit 处理编辑规则
func (c *CollectController) RuleDoEdit(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取规则ID
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid rule ID", http.StatusBadRequest)
		return
	}

	// 获取原规则
	rule, err := c.collectRuleModel.GetByID(id)
	if err != nil {
		logger.Error("获取采集规则失败", "id", id, "error", err)
		http.Error(w, "Rule not found", http.StatusNotFound)
		return
	}

	// 获取表单数据
	name := r.FormValue("name")
	sourceTypeStr := r.FormValue("sourcetype")
	sourceURL := r.FormValue("sourceurl")
	startPageStr := r.FormValue("startpage")
	endPageStr := r.FormValue("endpage")
	pageRule := r.FormValue("pagerule")
	listRule := r.FormValue("listrule")
	titleRule := r.FormValue("titlerule")
	contentRule := r.FormValue("contentrule")
	typeIDStr := r.FormValue("typeid")
	modelIDStr := r.FormValue("modelid")
	statusStr := r.FormValue("status")
	fieldRules := r.FormValue("fieldrules")

	// 验证必填字段
	if name == "" || sourceURL == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 解析数值
	sourceType, _ := strconv.Atoi(sourceTypeStr)
	startPage, _ := strconv.Atoi(startPageStr)
	endPage, _ := strconv.Atoi(endPageStr)
	typeID, _ := strconv.ParseInt(typeIDStr, 10, 64)
	modelID, _ := strconv.ParseInt(modelIDStr, 10, 64)
	status, _ := strconv.Atoi(statusStr)

	// 更新规则
	rule.Name = name
	rule.SourceType = sourceType
	rule.SourceURL = sourceURL
	rule.StartPage = startPage
	rule.EndPage = endPage
	rule.PageRule = pageRule
	rule.ListRule = listRule
	rule.TitleRule = titleRule
	rule.ContentRule = contentRule
	rule.TypeID = typeID
	rule.ModelID = modelID
	rule.Status = status
	rule.FieldRules = fieldRules

	// 保存规则
	err = c.collectRuleModel.Update(rule)
	if err != nil {
		logger.Error("更新采集规则失败", "error", err)
		http.Error(w, "Failed to update collect rule", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "采集规则更新成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/collect_rule_list", http.StatusFound)
	}
}

// RuleDelete 删除规则
func (c *CollectController) RuleDelete(w http.ResponseWriter, r *http.Request) {
	// 获取规则ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid rule ID", http.StatusBadRequest)
		return
	}

	// 删除规则
	err = c.collectService.DeleteRule(id)
	if err != nil {
		logger.Error("删除采集规则失败", "id", id, "error", err)
		http.Error(w, "Failed to delete collect rule", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "采集规则删除成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/collect_rule_list", http.StatusFound)
	}
}

// RuleCollect 采集规则
func (c *CollectController) RuleCollect(w http.ResponseWriter, r *http.Request) {
	// 获取规则ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid rule ID", http.StatusBadRequest)
		return
	}

	// 采集规则
	count, err := c.collectService.CollectRule(id)
	if err != nil {
		logger.Error("采集规则失败", "id", id, "error", err)
		http.Error(w, "Failed to collect rule", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "采集成功，共采集" + strconv.Itoa(count) + "条数据",
			"count":   count,
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/collect_item_list/"+idStr, http.StatusFound)
	}
}

// ItemList 项目列表
func (c *CollectController) ItemList(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取规则ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid rule ID", http.StatusBadRequest)
		return
	}

	// 获取规则
	rule, err := c.collectRuleModel.GetByID(id)
	if err != nil {
		logger.Error("获取采集规则失败", "id", id, "error", err)
		http.Error(w, "Rule not found", http.StatusNotFound)
		return
	}

	// 获取查询参数
	pageStr := r.URL.Query().Get("page")
	statusStr := r.URL.Query().Get("status")

	// 解析参数
	page := 1
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
		if page < 1 {
			page = 1
		}
	}

	status := -1
	if statusStr != "" {
		status, _ = strconv.Atoi(statusStr)
	}

	// 获取项目列表
	items, total, err := c.collectItemModel.GetByRuleID(id, status, page, 20)
	if err != nil {
		logger.Error("获取采集项目列表失败", "error", err)
		http.Error(w, "Failed to get collect items", http.StatusInternalServerError)
		return
	}

	// 计算分页信息
	totalPages := (total + 20 - 1) / 20
	pagination := map[string]interface{}{
		"CurrentPage": page,
		"TotalPages":  totalPages,
		"TotalItems":  total,
		"HasPrev":     page > 1,
		"HasNext":     page < totalPages,
		"PrevPage":    page - 1,
		"NextPage":    page + 1,
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Rule":        rule,
		"Items":       items,
		"Status":      status,
		"Pagination":  pagination,
		"CurrentMenu": "collect",
		"PageTitle":   "采集项目管理",
	}

	// 渲染模板
	tplFile := "admin/collect_item_list.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染采集项目列表模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// ItemDetail 项目详情
func (c *CollectController) ItemDetail(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取项目ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	// 获取项目
	item, err := c.collectItemModel.GetByID(id)
	if err != nil {
		logger.Error("获取采集项目失败", "id", id, "error", err)
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	// 获取规则
	rule, err := c.collectRuleModel.GetByID(item.RuleID)
	if err != nil {
		logger.Error("获取采集规则失败", "id", item.RuleID, "error", err)
		http.Error(w, "Rule not found", http.StatusNotFound)
		return
	}

	// 获取字段数据
	fieldData, err := c.collectItemModel.GetFieldData(id)
	if err != nil {
		logger.Error("获取字段数据失败", "error", err)
		http.Error(w, "Failed to get field data", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Rule":        rule,
		"Item":        item,
		"FieldData":   fieldData,
		"CurrentMenu": "collect",
		"PageTitle":   "采集项目详情",
	}

	// 渲染模板
	tplFile := "admin/collect_item_detail.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染采集项目详情模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// ItemPublish 发布项目
func (c *CollectController) ItemPublish(w http.ResponseWriter, r *http.Request) {
	// 获取项目ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	// 获取项目
	item, err := c.collectItemModel.GetByID(id)
	if err != nil {
		logger.Error("获取采集项目失败", "id", id, "error", err)
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	// 发布项目
	err = c.collectService.PublishItem(id)
	if err != nil {
		logger.Error("发布采集项目失败", "id", id, "error", err)
		http.Error(w, "Failed to publish collect item", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "采集项目发布成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/collect_item_list/"+strconv.FormatInt(item.RuleID, 10), http.StatusFound)
	}
}

// ItemDelete 删除项目
func (c *CollectController) ItemDelete(w http.ResponseWriter, r *http.Request) {
	// 获取项目ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	// 获取项目
	item, err := c.collectItemModel.GetByID(id)
	if err != nil {
		logger.Error("获取采集项目失败", "id", id, "error", err)
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	// 删除项目
	err = c.collectItemModel.Delete(id)
	if err != nil {
		logger.Error("删除采集项目失败", "id", id, "error", err)
		http.Error(w, "Failed to delete collect item", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "采集项目删除成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/collect_item_list/"+strconv.FormatInt(item.RuleID, 10), http.StatusFound)
	}
}

// BatchPublish 批量发布
func (c *CollectController) BatchPublish(w http.ResponseWriter, r *http.Request) {
	// 获取规则ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid rule ID", http.StatusBadRequest)
		return
	}

	// 批量发布
	count, err := c.collectService.BatchPublish(id)
	if err != nil {
		logger.Error("批量发布采集项目失败", "id", id, "error", err)
		http.Error(w, "Failed to batch publish collect items", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "批量发布成功，共发布" + strconv.Itoa(count) + "条数据",
			"count":   count,
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/collect_item_list/"+idStr, http.StatusFound)
	}
}
