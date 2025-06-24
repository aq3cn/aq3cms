package admin

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strconv"

	"aq3cms/config"
	"aq3cms/internal/middleware"
	"aq3cms/internal/model"
	"aq3cms/internal/service"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"

	"github.com/gorilla/mux"
)

// CategoryController 栏目控制器
type CategoryController struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	categoryModel   *model.CategoryModel
	htmlService     *service.HtmlService
	templateService *service.TemplateService
}

// NewCategoryController 创建栏目控制器
func NewCategoryController(db *database.DB, cache cache.Cache, config *config.Config) *CategoryController {
	return &CategoryController{
		db:              db,
		cache:           cache,
		config:          config,
		categoryModel:   model.NewCategoryModel(db),
		htmlService:     service.NewHtmlService(db, cache, config),
		templateService: service.NewTemplateService(db, cache, config),
	}
}

// Index 栏目管理主页
func (c *CategoryController) Index(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取栏目列表
	categories, err := c.categoryModel.GetAll()
	if err != nil {
		logger.Error("获取栏目列表失败", "error", err)
		categories = []*model.Category{}
	}

	// 获取栏目统计信息
	totalCategories := len(categories)

	// 计算顶级栏目数
	topCategories := 0
	for _, category := range categories {
		if category.ParentID == 0 {
			topCategories++
		}
	}

	// 构建栏目树
	categoryTree := buildCategoryTree(categories)

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":         adminID,
		"AdminName":       adminName,
		"TotalCategories": totalCategories,
		"TopCategories":   topCategories,
		"Categories":      categories,
		"CategoryTree":    categoryTree,
		"CurrentMenu":     "category",
		"PageTitle":       "栏目管理",
	}

	// 渲染模板
	tplFile := "admin/category_index.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染栏目管理主页模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// List 栏目列表
func (c *CategoryController) List(w http.ResponseWriter, r *http.Request) {
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

	// 为每个栏目获取文章数量
	for _, category := range categories {
		count, err := c.categoryModel.GetArticleCount(category.ID)
		if err != nil {
			logger.Error("获取栏目文章数量失败", "categoryID", category.ID, "error", err)
			count = 0
		}
		category.ArticleCount = count
	}

	// 构建栏目树
	categoryTree := buildCategoryTree(categories)

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":      adminID,
		"AdminName":    adminName,
		"Categories":   categories,
		"CategoryTree": categoryTree,
		"CurrentMenu":  "category",
		"PageTitle":    "栏目列表",
	}

	// 渲染模板
	tplFile := "admin/category_list.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染栏目列表模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Add 添加栏目页面
func (c *CategoryController) Add(w http.ResponseWriter, r *http.Request) {
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

	// 获取模板列表
	templates, err := getTemplateList(c.config.Template.Dir)
	if err != nil {
		logger.Error("获取模板列表失败", "error", err)
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Categories":  categories,
		"Templates":   templates,
		"CurrentMenu": "category",
		"PageTitle":   "添加栏目",
	}

	// 渲染模板
	tplFile := "admin/category_add.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染添加栏目模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoAdd 处理添加栏目
func (c *CategoryController) DoAdd(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	typeName := r.FormValue("typename")
	parentIDStr := r.FormValue("reid")
	typeDir := r.FormValue("typedir")
	channelTypeStr := r.FormValue("channeltype")
	isHtmlStr := r.FormValue("ishtml")
	_ = r.FormValue("defaultname") // 未使用
	_ = r.FormValue("isdefault")   // 未使用
	sortRankStr := r.FormValue("sortrank")
	keywords := r.FormValue("keywords")
	description := r.FormValue("description")
	listTemplate := r.FormValue("list_template")
	articleTemplate := r.FormValue("article_template")

	// 验证必填字段
	if typeName == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 解析字段
	parentID := int64(0)
	if parentIDStr != "" {
		parentID, _ = strconv.ParseInt(parentIDStr, 10, 64)
	}

	channelType := 1
	if channelTypeStr != "" {
		channelType, _ = strconv.Atoi(channelTypeStr)
	}

	isHtml := 0
	if isHtmlStr == "1" {
		isHtml = 1
	}

	// 这些变量在当前版本的模型中不使用

	sortRank := 0
	if sortRankStr != "" {
		sortRank, _ = strconv.Atoi(sortRankStr)
	}

	// 创建栏目
	category := &model.Category{
		TypeName:    typeName,
		ParentID:    parentID,
		TypeDir:     typeDir,
		ChannelType: channelType,
		IsHidden:    isHtml,
		SortRank:    sortRank,
		Keywords:    keywords,
		Description: description,
		ListTpl:     listTemplate,
		ArticleTpl:  articleTemplate,
	}

	// 保存栏目
	id, err := c.categoryModel.Create(category)
	if err != nil {
		logger.Error("创建栏目失败", "error", err)
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		return
	}

	// 生成静态页面
	go c.htmlService.GenerateList(id)

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "栏目创建成功",
			"id":      id,
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/aq3cms/category_list", http.StatusFound)
	}
}

// Edit 编辑栏目页面
func (c *CategoryController) Edit(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取栏目ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	// 获取栏目
	category, err := c.categoryModel.GetByID(id)
	if err != nil {
		logger.Error("获取栏目失败", "id", id, "error", err)
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	// 获取栏目列表
	categories, err := c.categoryModel.GetAll()
	if err != nil {
		logger.Error("获取栏目列表失败", "error", err)
	}

	// 获取模板列表
	templates, err := getTemplateList(c.config.Template.Dir)
	if err != nil {
		logger.Error("获取模板列表失败", "error", err)
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Category":    category,
		"Categories":  categories,
		"Templates":   templates,
		"CurrentMenu": "category",
		"PageTitle":   "编辑栏目",
	}

	// 渲染模板
	tplFile := "admin/category_edit.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染编辑栏目模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoEdit 处理编辑栏目
func (c *CategoryController) DoEdit(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取栏目ID
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	// 获取原栏目
	category, err := c.categoryModel.GetByID(id)
	if err != nil {
		logger.Error("获取栏目失败", "id", id, "error", err)
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	// 获取表单数据
	typeName := r.FormValue("typename")
	parentIDStr := r.FormValue("reid")
	typeDir := r.FormValue("typedir")
	channelTypeStr := r.FormValue("channeltype")
	isHtmlStr := r.FormValue("ishtml")
	_ = r.FormValue("defaultname") // 未使用
	_ = r.FormValue("isdefault")   // 未使用
	sortRankStr := r.FormValue("sortrank")
	keywords := r.FormValue("keywords")
	description := r.FormValue("description")
	listTemplate := r.FormValue("list_template")
	articleTemplate := r.FormValue("article_template")

	// 验证必填字段
	if typeName == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 解析字段
	parentID := int64(0)
	if parentIDStr != "" {
		parentID, _ = strconv.ParseInt(parentIDStr, 10, 64)
	}

	// 检查是否形成循环引用
	if parentID == id {
		http.Error(w, "Cannot set parent to self", http.StatusBadRequest)
		return
	}

	// 检查是否将栏目设置为其子栏目的父栏目
	children, err := c.categoryModel.GetChildCategories(id)
	if err == nil {
		for _, child := range children {
			if child.ID == parentID {
				http.Error(w, "Cannot set parent to child category", http.StatusBadRequest)
				return
			}
		}
	}

	channelType := 1
	if channelTypeStr != "" {
		channelType, _ = strconv.Atoi(channelTypeStr)
	}

	isHtml := 0
	if isHtmlStr == "1" {
		isHtml = 1
	}

	// 这些变量在当前版本的模型中不使用

	sortRank := 0
	if sortRankStr != "" {
		sortRank, _ = strconv.Atoi(sortRankStr)
	}

	// 更新栏目
	category.TypeName = typeName
	category.ParentID = parentID
	category.TypeDir = typeDir
	category.ChannelType = channelType
	category.IsHidden = isHtml
	category.SortRank = sortRank
	category.Keywords = keywords
	category.Description = description
	category.ListTpl = listTemplate
	category.ArticleTpl = articleTemplate

	// 保存栏目
	err = c.categoryModel.Update(category)
	if err != nil {
		logger.Error("更新栏目失败", "error", err)
		http.Error(w, "Failed to update category", http.StatusInternalServerError)
		return
	}

	// 生成静态页面
	go c.htmlService.GenerateList(id)

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "栏目更新成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/aq3cms/category_list", http.StatusFound)
	}
}

// Delete 删除栏目
func (c *CategoryController) Delete(w http.ResponseWriter, r *http.Request) {
	// 获取栏目ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	// 检查是否有子栏目
	children, err := c.categoryModel.GetChildCategories(id)
	if err == nil && len(children) > 0 {
		http.Error(w, "Cannot delete category with child categories", http.StatusBadRequest)
		return
	}

	// 检查是否有文章
	count, err := c.categoryModel.GetArticleCount(id)
	if err == nil && count > 0 {
		http.Error(w, "Cannot delete category with articles", http.StatusBadRequest)
		return
	}

	// 删除栏目
	err = c.categoryModel.Delete(id)
	if err != nil {
		logger.Error("删除栏目失败", "id", id, "error", err)
		http.Error(w, "Failed to delete category", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "栏目删除成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/aq3cms/category_list", http.StatusFound)
	}
}

// 构建栏目树
func buildCategoryTree(categories []*model.Category) []*model.Category {
	// 创建ID到栏目的映射
	categoryMap := make(map[int64]*model.Category)
	for _, category := range categories {
		categoryMap[category.ID] = category
		category.Children = make([]*model.Category, 0)
	}

	// 构建树
	rootCategories := make([]*model.Category, 0)
	for _, category := range categories {
		if category.ParentID == 0 {
			rootCategories = append(rootCategories, category)
		} else {
			if parent, ok := categoryMap[category.ParentID]; ok {
				parent.Children = append(parent.Children, category)
			}
		}
	}

	return rootCategories
}

// 获取模板列表
func getTemplateList(templateDir string) ([]string, error) {
	// 获取模板文件列表
	files, err := filepath.Glob(filepath.Join(templateDir, "*.htm"))
	if err != nil {
		return nil, err
	}

	// 提取模板名称
	templates := make([]string, 0, len(files))
	for _, file := range files {
		templates = append(templates, filepath.Base(file))
	}

	return templates, nil
}
