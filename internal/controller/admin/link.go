package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"aq3cms/config"
	"aq3cms/internal/middleware"
	"aq3cms/internal/model"
	"aq3cms/internal/service"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// LinkController 友情链接控制器
type LinkController struct {
	db             *database.DB
	cache          cache.Cache
	config         *config.Config
	linkModel      *model.LinkModel
	linkTypeModel  *model.LinkTypeModel
	templateService *service.TemplateService
}

// NewLinkController 创建友情链接控制器
func NewLinkController(db *database.DB, cache cache.Cache, config *config.Config) *LinkController {
	return &LinkController{
		db:             db,
		cache:          cache,
		config:         config,
		linkModel:      model.NewLinkModel(db),
		linkTypeModel:  model.NewLinkTypeModel(db),
		templateService: service.NewTemplateService(db, cache, config),
	}
}

// List 友情链接列表
func (c *LinkController) List(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取查询参数
	typeIDStr := r.URL.Query().Get("typeid")

	// 解析参数
	typeID := int64(0)
	if typeIDStr != "" {
		typeID, _ = strconv.ParseInt(typeIDStr, 10, 64)
	}

	// 获取友情链接列表
	links, err := c.linkModel.GetAll(typeID, -1)
	if err != nil {
		logger.Error("获取友情链接列表失败", "error", err)
		http.Error(w, "Failed to get links", http.StatusInternalServerError)
		return
	}

	// 获取友情链接分类列表
	linkTypes, err := c.linkTypeModel.GetAll()
	if err != nil {
		logger.Error("获取友情链接分类列表失败", "error", err)
		http.Error(w, "Failed to get link types", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Links":       links,
		"LinkTypes":   linkTypes,
		"TypeID":      typeID,
		"CurrentMenu": "link",
		"PageTitle":   "友情链接管理",
	}

	// 渲染模板
	tplFile := "admin/link_list.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染友情链接列表模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Add 添加友情链接页面
func (c *LinkController) Add(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取友情链接分类列表
	linkTypes, err := c.linkTypeModel.GetAll()
	if err != nil {
		logger.Error("获取友情链接分类列表失败", "error", err)
		http.Error(w, "Failed to get link types", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"LinkTypes":   linkTypes,
		"CurrentMenu": "link",
		"PageTitle":   "添加友情链接",
	}

	// 渲染模板
	tplFile := "admin/link_add.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染添加友情链接模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoAdd 处理添加友情链接
func (c *LinkController) DoAdd(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	typeIDStr := r.FormValue("typeid")
	title := r.FormValue("title")
	url := r.FormValue("url")
	logo := r.FormValue("logo")
	description := r.FormValue("description")
	email := r.FormValue("email")
	orderIDStr := r.FormValue("orderid")
	statusStr := r.FormValue("status")
	isLogoStr := r.FormValue("islogo")

	// 验证必填字段
	if title == "" || url == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 解析数值
	typeID, _ := strconv.ParseInt(typeIDStr, 10, 64)
	orderID, _ := strconv.Atoi(orderIDStr)
	status, _ := strconv.Atoi(statusStr)
	isLogo, _ := strconv.Atoi(isLogoStr)

	// 创建友情链接
	link := &model.Link{
		TypeID:      typeID,
		Title:       title,
		URL:         url,
		Logo:        logo,
		Description: description,
		Email:       email,
		OrderID:     orderID,
		IsCheck:     status,
		IsLogo:      isLogo,
	}

	// 保存友情链接
	id, err := c.linkModel.Create(link)
	if err != nil {
		logger.Error("创建友情链接失败", "error", err)
		http.Error(w, "Failed to create link", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "友情链接创建成功",
			"id":      id,
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/link_list", http.StatusFound)
	}
}

// Edit 编辑友情链接页面
func (c *LinkController) Edit(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取友情链接ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid link ID", http.StatusBadRequest)
		return
	}

	// 获取友情链接
	link, err := c.linkModel.GetByID(id)
	if err != nil {
		logger.Error("获取友情链接失败", "id", id, "error", err)
		http.Error(w, "Link not found", http.StatusNotFound)
		return
	}

	// 获取友情链接分类列表
	linkTypes, err := c.linkTypeModel.GetAll()
	if err != nil {
		logger.Error("获取友情链接分类列表失败", "error", err)
		http.Error(w, "Failed to get link types", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Link":        link,
		"LinkTypes":   linkTypes,
		"CurrentMenu": "link",
		"PageTitle":   "编辑友情链接",
	}

	// 渲染模板
	tplFile := "admin/link_edit.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染编辑友情链接模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoEdit 处理编辑友情链接
func (c *LinkController) DoEdit(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取友情链接ID
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid link ID", http.StatusBadRequest)
		return
	}

	// 获取原友情链接
	link, err := c.linkModel.GetByID(id)
	if err != nil {
		logger.Error("获取友情链接失败", "id", id, "error", err)
		http.Error(w, "Link not found", http.StatusNotFound)
		return
	}

	// 获取表单数据
	typeIDStr := r.FormValue("typeid")
	title := r.FormValue("title")
	url := r.FormValue("url")
	logo := r.FormValue("logo")
	description := r.FormValue("description")
	email := r.FormValue("email")
	orderIDStr := r.FormValue("orderid")
	statusStr := r.FormValue("status")
	isLogoStr := r.FormValue("islogo")

	// 验证必填字段
	if title == "" || url == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 解析数值
	typeID, _ := strconv.ParseInt(typeIDStr, 10, 64)
	orderID, _ := strconv.Atoi(orderIDStr)
	status, _ := strconv.Atoi(statusStr)
	isLogo, _ := strconv.Atoi(isLogoStr)

	// 更新友情链接
	link.TypeID = typeID
	link.Title = title
	link.URL = url
	link.Logo = logo
	link.Description = description
	link.Email = email
	link.OrderID = orderID
	link.IsCheck = status
	link.IsLogo = isLogo

	// 保存友情链接
	err = c.linkModel.Update(link)
	if err != nil {
		logger.Error("更新友情链接失败", "error", err)
		http.Error(w, "Failed to update link", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "友情链接更新成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/link_list", http.StatusFound)
	}
}

// Delete 删除友情链接
func (c *LinkController) Delete(w http.ResponseWriter, r *http.Request) {
	// 获取友情链接ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid link ID", http.StatusBadRequest)
		return
	}

	// 删除友情链接
	err = c.linkModel.Delete(id)
	if err != nil {
		logger.Error("删除友情链接失败", "id", id, "error", err)
		http.Error(w, "Failed to delete link", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "友情链接删除成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/link_list", http.StatusFound)
	}
}

// TypeList 友情链接分类列表
func (c *LinkController) TypeList(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取友情链接分类列表
	linkTypes, err := c.linkTypeModel.GetAll()
	if err != nil {
		logger.Error("获取友情链接分类列表失败", "error", err)
		http.Error(w, "Failed to get link types", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"LinkTypes":   linkTypes,
		"CurrentMenu": "link",
		"PageTitle":   "友情链接分类管理",
	}

	// 渲染模板
	tplFile := "admin/link_type_list.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染友情链接分类列表模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// TypeAdd 添加友情链接分类页面
func (c *LinkController) TypeAdd(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"CurrentMenu": "link",
		"PageTitle":   "添加友情链接分类",
	}

	// 渲染模板
	tplFile := "admin/link_type_add.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染添加友情链接分类模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// TypeDoAdd 处理添加友情链接分类
func (c *LinkController) TypeDoAdd(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	name := r.FormValue("name")
	orderIDStr := r.FormValue("orderid")

	// 验证必填字段
	if name == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 解析数值
	orderID, _ := strconv.Atoi(orderIDStr)

	// 创建友情链接分类
	linkType := &model.LinkType{
		Name:    name,
		OrderID: orderID,
	}

	// 保存友情链接分类
	id, err := c.linkTypeModel.Create(linkType)
	if err != nil {
		logger.Error("创建友情链接分类失败", "error", err)
		http.Error(w, "Failed to create link type", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "友情链接分类创建成功",
			"id":      id,
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/link_type_list", http.StatusFound)
	}
}

// TypeEdit 编辑友情链接分类页面
func (c *LinkController) TypeEdit(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取友情链接分类ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid link type ID", http.StatusBadRequest)
		return
	}

	// 获取友情链接分类
	linkType, err := c.linkTypeModel.GetByID(id)
	if err != nil {
		logger.Error("获取友情链接分类失败", "id", id, "error", err)
		http.Error(w, "Link type not found", http.StatusNotFound)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"LinkType":    linkType,
		"CurrentMenu": "link",
		"PageTitle":   "编辑友情链接分类",
	}

	// 渲染模板
	tplFile := "admin/link_type_edit.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染编辑友情链接分类模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// TypeDoEdit 处理编辑友情链接分类
func (c *LinkController) TypeDoEdit(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取友情链接分类ID
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid link type ID", http.StatusBadRequest)
		return
	}

	// 获取原友情链接分类
	linkType, err := c.linkTypeModel.GetByID(id)
	if err != nil {
		logger.Error("获取友情链接分类失败", "id", id, "error", err)
		http.Error(w, "Link type not found", http.StatusNotFound)
		return
	}

	// 获取表单数据
	name := r.FormValue("name")
	orderIDStr := r.FormValue("orderid")

	// 验证必填字段
	if name == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 解析数值
	orderID, _ := strconv.Atoi(orderIDStr)

	// 更新友情链接分类
	linkType.Name = name
	linkType.OrderID = orderID

	// 保存友情链接分类
	err = c.linkTypeModel.Update(linkType)
	if err != nil {
		logger.Error("更新友情链接分类失败", "error", err)
		http.Error(w, "Failed to update link type", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "友情链接分类更新成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/link_type_list", http.StatusFound)
	}
}

// TypeDelete 删除友情链接分类
func (c *LinkController) TypeDelete(w http.ResponseWriter, r *http.Request) {
	// 获取友情链接分类ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid link type ID", http.StatusBadRequest)
		return
	}

	// 检查分类是否有友情链接
	hasLinks, err := c.linkTypeModel.HasLinks(id)
	if err != nil {
		logger.Error("检查分类是否有友情链接失败", "id", id, "error", err)
		http.Error(w, "Failed to check if link type has links", http.StatusInternalServerError)
		return
	}
	if hasLinks {
		http.Error(w, "Link type has links, cannot delete", http.StatusBadRequest)
		return
	}

	// 删除友情链接分类
	err = c.linkTypeModel.Delete(id)
	if err != nil {
		logger.Error("删除友情链接分类失败", "id", id, "error", err)
		http.Error(w, "Failed to delete link type", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "友情链接分类删除成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/link_type_list", http.StatusFound)
	}
}
