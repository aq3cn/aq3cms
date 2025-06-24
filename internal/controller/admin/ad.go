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

// AdController 广告控制器
type AdController struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	adService       *service.AdService
	adModel         *model.AdModel
	adPositionModel *model.AdPositionModel
	templateService *service.TemplateService
}

// NewAdController 创建广告控制器
func NewAdController(db *database.DB, cache cache.Cache, config *config.Config) *AdController {
	return &AdController{
		db:              db,
		cache:           cache,
		config:          config,
		adService:       service.NewAdService(db, cache, config),
		adModel:         model.NewAdModel(db),
		adPositionModel: model.NewAdPositionModel(db),
		templateService: service.NewTemplateService(db, cache, config),
	}
}

// List 广告列表
func (c *AdController) List(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取查询参数
	positionIDStr := r.URL.Query().Get("positionid")

	// 解析参数
	positionID := int64(0)
	if positionIDStr != "" {
		positionID, _ = strconv.ParseInt(positionIDStr, 10, 64)
	}

	// 获取广告列表
	ads, err := c.adModel.GetAll(positionID, -1)
	if err != nil {
		logger.Error("获取广告列表失败", "error", err)
		http.Error(w, "Failed to get ads", http.StatusInternalServerError)
		return
	}

	// 获取广告位列表
	positions, err := c.adPositionModel.GetAll(-1)
	if err != nil {
		logger.Error("获取广告位列表失败", "error", err)
		http.Error(w, "Failed to get ad positions", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Ads":         ads,
		"Positions":   positions,
		"PositionID":  positionID,
		"CurrentMenu": "ad",
		"PageTitle":   "广告管理",
	}

	// 渲染模板
	tplFile := "admin/ad_list.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染广告列表模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Add 添加广告页面
func (c *AdController) Add(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取广告位列表
	positions, err := c.adPositionModel.GetAll(1)
	if err != nil {
		logger.Error("获取广告位列表失败", "error", err)
		http.Error(w, "Failed to get ad positions", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Positions":   positions,
		"CurrentMenu": "ad",
		"PageTitle":   "添加广告",
	}

	// 渲染模板
	tplFile := "admin/ad_add.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染添加广告模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoAdd 处理添加广告
func (c *AdController) DoAdd(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	positionIDStr := r.FormValue("positionid")
	title := r.FormValue("title")
	typeStr := r.FormValue("type")
	image := r.FormValue("image")
	flash := r.FormValue("flash")
	code := r.FormValue("code")
	text := r.FormValue("text")
	url := r.FormValue("url")
	startTimeStr := r.FormValue("starttime")
	endTimeStr := r.FormValue("endtime")
	orderIDStr := r.FormValue("orderid")
	statusStr := r.FormValue("status")
	target := r.FormValue("target")
	widthStr := r.FormValue("width")
	heightStr := r.FormValue("height")

	// 验证必填字段
	if positionIDStr == "" || title == "" || typeStr == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 解析数值
	positionID, _ := strconv.ParseInt(positionIDStr, 10, 64)
	adType, _ := strconv.Atoi(typeStr)
	orderID, _ := strconv.Atoi(orderIDStr)
	status, _ := strconv.Atoi(statusStr)
	width, _ := strconv.Atoi(widthStr)
	height, _ := strconv.Atoi(heightStr)

	// 解析时间
	startTime, err := time.Parse("2006-01-02 15:04:05", startTimeStr)
	if err != nil {
		startTime = time.Now()
	}
	endTime, err := time.Parse("2006-01-02 15:04:05", endTimeStr)
	if err != nil {
		endTime = time.Now().AddDate(1, 0, 0)
	}

	// 创建广告
	ad := &model.Ad{
		PositionID: positionID,
		Title:      title,
		Type:       adType,
		Image:      image,
		Flash:      flash,
		Code:       code,
		Text:       text,
		URL:        url,
		StartTime:  startTime,
		EndTime:    endTime,
		OrderID:    orderID,
		Status:     status,
		Target:     target,
		Width:      width,
		Height:     height,
		Click:      0,
	}

	// 保存广告
	id, err := c.adModel.Create(ad)
	if err != nil {
		logger.Error("创建广告失败", "error", err)
		http.Error(w, "Failed to create ad", http.StatusInternalServerError)
		return
	}

	// 清除缓存
	c.adService.ClearCache(id)

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "广告创建成功",
			"id":      id,
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/ad_list", http.StatusFound)
	}
}

// Edit 编辑广告页面
func (c *AdController) Edit(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取广告ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ad ID", http.StatusBadRequest)
		return
	}

	// 获取广告
	ad, err := c.adModel.GetByID(id)
	if err != nil {
		logger.Error("获取广告失败", "id", id, "error", err)
		http.Error(w, "Ad not found", http.StatusNotFound)
		return
	}

	// 获取广告位列表
	positions, err := c.adPositionModel.GetAll(1)
	if err != nil {
		logger.Error("获取广告位列表失败", "error", err)
		http.Error(w, "Failed to get ad positions", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Ad":          ad,
		"Positions":   positions,
		"CurrentMenu": "ad",
		"PageTitle":   "编辑广告",
	}

	// 渲染模板
	tplFile := "admin/ad_edit.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染编辑广告模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoEdit 处理编辑广告
func (c *AdController) DoEdit(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取广告ID
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ad ID", http.StatusBadRequest)
		return
	}

	// 获取原广告
	ad, err := c.adModel.GetByID(id)
	if err != nil {
		logger.Error("获取广告失败", "id", id, "error", err)
		http.Error(w, "Ad not found", http.StatusNotFound)
		return
	}

	// 获取表单数据
	positionIDStr := r.FormValue("positionid")
	title := r.FormValue("title")
	typeStr := r.FormValue("type")
	image := r.FormValue("image")
	flash := r.FormValue("flash")
	code := r.FormValue("code")
	text := r.FormValue("text")
	url := r.FormValue("url")
	startTimeStr := r.FormValue("starttime")
	endTimeStr := r.FormValue("endtime")
	orderIDStr := r.FormValue("orderid")
	statusStr := r.FormValue("status")
	target := r.FormValue("target")
	widthStr := r.FormValue("width")
	heightStr := r.FormValue("height")

	// 验证必填字段
	if positionIDStr == "" || title == "" || typeStr == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 解析数值
	positionID, _ := strconv.ParseInt(positionIDStr, 10, 64)
	adType, _ := strconv.Atoi(typeStr)
	orderID, _ := strconv.Atoi(orderIDStr)
	status, _ := strconv.Atoi(statusStr)
	width, _ := strconv.Atoi(widthStr)
	height, _ := strconv.Atoi(heightStr)

	// 解析时间
	startTime, err := time.Parse("2006-01-02 15:04:05", startTimeStr)
	if err != nil {
		startTime = time.Now()
	}
	endTime, err := time.Parse("2006-01-02 15:04:05", endTimeStr)
	if err != nil {
		endTime = time.Now().AddDate(1, 0, 0)
	}

	// 更新广告
	ad.PositionID = positionID
	ad.Title = title
	ad.Type = adType
	ad.Image = image
	ad.Flash = flash
	ad.Code = code
	ad.Text = text
	ad.URL = url
	ad.StartTime = startTime
	ad.EndTime = endTime
	ad.OrderID = orderID
	ad.Status = status
	ad.Target = target
	ad.Width = width
	ad.Height = height

	// 保存广告
	err = c.adModel.Update(ad)
	if err != nil {
		logger.Error("更新广告失败", "error", err)
		http.Error(w, "Failed to update ad", http.StatusInternalServerError)
		return
	}

	// 清除缓存
	c.adService.ClearCache(id)

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "广告更新成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/ad_list", http.StatusFound)
	}
}

// Delete 删除广告
func (c *AdController) Delete(w http.ResponseWriter, r *http.Request) {
	// 获取广告ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ad ID", http.StatusBadRequest)
		return
	}

	// 获取广告
	_, err = c.adModel.GetByID(id)
	if err != nil {
		logger.Error("获取广告失败", "id", id, "error", err)
		http.Error(w, "Ad not found", http.StatusNotFound)
		return
	}

	// 清除缓存
	c.adService.ClearCache(id)

	// 删除广告
	err = c.adModel.Delete(id)
	if err != nil {
		logger.Error("删除广告失败", "id", id, "error", err)
		http.Error(w, "Failed to delete ad", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "广告删除成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/ad_list", http.StatusFound)
	}
}

// PositionList 广告位列表
func (c *AdController) PositionList(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取广告位列表
	positions, err := c.adPositionModel.GetAll(-1)
	if err != nil {
		logger.Error("获取广告位列表失败", "error", err)
		http.Error(w, "Failed to get ad positions", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Positions":   positions,
		"CurrentMenu": "ad",
		"PageTitle":   "广告位管理",
	}

	// 渲染模板
	tplFile := "admin/ad_position_list.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染广告位列表模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// PositionAdd 添加广告位页面
func (c *AdController) PositionAdd(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"CurrentMenu": "ad",
		"PageTitle":   "添加广告位",
	}

	// 渲染模板
	tplFile := "admin/ad_position_add.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染添加广告位模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// PositionDoAdd 处理添加广告位
func (c *AdController) PositionDoAdd(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	name := r.FormValue("name")
	code := r.FormValue("code")
	widthStr := r.FormValue("width")
	heightStr := r.FormValue("height")
	template := r.FormValue("template")
	description := r.FormValue("description")
	statusStr := r.FormValue("status")

	// 验证必填字段
	if name == "" || code == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 解析数值
	width, _ := strconv.Atoi(widthStr)
	height, _ := strconv.Atoi(heightStr)
	status, _ := strconv.Atoi(statusStr)

	// 创建广告位
	position := &model.AdPosition{
		Name:        name,
		Code:        code,
		Width:       width,
		Height:      height,
		Template:    template,
		Description: description,
		Status:      status,
	}

	// 保存广告位
	id, err := c.adPositionModel.Create(position)
	if err != nil {
		logger.Error("创建广告位失败", "error", err)
		http.Error(w, "Failed to create ad position", http.StatusInternalServerError)
		return
	}

	// 清除缓存
	c.adService.ClearAllCache()

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "广告位创建成功",
			"id":      id,
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/ad_position_list", http.StatusFound)
	}
}

// PositionEdit 编辑广告位页面
func (c *AdController) PositionEdit(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取广告位ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ad position ID", http.StatusBadRequest)
		return
	}

	// 获取广告位
	position, err := c.adPositionModel.GetByID(id)
	if err != nil {
		logger.Error("获取广告位失败", "id", id, "error", err)
		http.Error(w, "Ad position not found", http.StatusNotFound)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Position":    position,
		"CurrentMenu": "ad",
		"PageTitle":   "编辑广告位",
	}

	// 渲染模板
	tplFile := "admin/ad_position_edit.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染编辑广告位模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// PositionDoEdit 处理编辑广告位
func (c *AdController) PositionDoEdit(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取广告位ID
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ad position ID", http.StatusBadRequest)
		return
	}

	// 获取原广告位
	position, err := c.adPositionModel.GetByID(id)
	if err != nil {
		logger.Error("获取广告位失败", "id", id, "error", err)
		http.Error(w, "Ad position not found", http.StatusNotFound)
		return
	}

	// 获取表单数据
	name := r.FormValue("name")
	code := r.FormValue("code")
	widthStr := r.FormValue("width")
	heightStr := r.FormValue("height")
	template := r.FormValue("template")
	description := r.FormValue("description")
	statusStr := r.FormValue("status")

	// 验证必填字段
	if name == "" || code == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 解析数值
	width, _ := strconv.Atoi(widthStr)
	height, _ := strconv.Atoi(heightStr)
	status, _ := strconv.Atoi(statusStr)

	// 更新广告位
	position.Name = name
	position.Code = code
	position.Width = width
	position.Height = height
	position.Template = template
	position.Description = description
	position.Status = status

	// 保存广告位
	err = c.adPositionModel.Update(position)
	if err != nil {
		logger.Error("更新广告位失败", "error", err)
		http.Error(w, "Failed to update ad position", http.StatusInternalServerError)
		return
	}

	// 清除缓存
	c.adService.ClearAllCache()

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "广告位更新成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/ad_position_list", http.StatusFound)
	}
}

// PositionDelete 删除广告位
func (c *AdController) PositionDelete(w http.ResponseWriter, r *http.Request) {
	// 获取广告位ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ad position ID", http.StatusBadRequest)
		return
	}

	// 检查广告位是否有广告
	hasAds, err := c.adPositionModel.HasAds(id)
	if err != nil {
		logger.Error("检查广告位是否有广告失败", "id", id, "error", err)
		http.Error(w, "Failed to check if ad position has ads", http.StatusInternalServerError)
		return
	}
	if hasAds {
		http.Error(w, "Ad position has ads, cannot delete", http.StatusBadRequest)
		return
	}

	// 删除广告位
	err = c.adPositionModel.Delete(id)
	if err != nil {
		logger.Error("删除广告位失败", "id", id, "error", err)
		http.Error(w, "Failed to delete ad position", http.StatusInternalServerError)
		return
	}

	// 清除缓存
	c.adService.ClearAllCache()

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "广告位删除成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/ad_position_list", http.StatusFound)
	}
}
