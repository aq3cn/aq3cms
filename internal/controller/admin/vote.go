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

// VoteController 投票控制器
type VoteController struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	voteService     *service.VoteService
	voteModel       *model.VoteModel
	voteOptionModel *model.VoteOptionModel
	voteLogModel    *model.VoteLogModel
	templateService *service.TemplateService
}

// NewVoteController 创建投票控制器
func NewVoteController(db *database.DB, cache cache.Cache, config *config.Config) *VoteController {
	return &VoteController{
		db:              db,
		cache:           cache,
		config:          config,
		voteService:     service.NewVoteService(db, cache, config),
		voteModel:       model.NewVoteModel(db),
		voteOptionModel: model.NewVoteOptionModel(db),
		voteLogModel:    model.NewVoteLogModel(db),
		templateService: service.NewTemplateService(db, cache, config),
	}
}

// List 投票列表
func (c *VoteController) List(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取查询参数
	pageStr := r.URL.Query().Get("page")

	// 解析参数
	page := 1
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
		if page < 1 {
			page = 1
		}
	}

	// 获取投票列表
	votes, total, err := c.voteModel.GetAll(page, 20)
	if err != nil {
		logger.Error("获取投票列表失败", "error", err)
		http.Error(w, "Failed to get votes", http.StatusInternalServerError)
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
		"Votes":       votes,
		"Pagination":  pagination,
		"CurrentMenu": "vote",
		"PageTitle":   "投票管理",
	}

	// 渲染模板
	tplFile := "admin/vote_list.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染投票列表模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Add 添加投票页面
func (c *VoteController) Add(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"CurrentMenu": "vote",
		"PageTitle":   "添加投票",
	}

	// 渲染模板
	tplFile := "admin/vote_add.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染添加投票模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoAdd 处理添加投票
func (c *VoteController) DoAdd(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	title := r.FormValue("title")
	description := r.FormValue("description")
	startTimeStr := r.FormValue("starttime")
	endTimeStr := r.FormValue("endtime")
	isMultiStr := r.FormValue("ismulti")
	maxChoicesStr := r.FormValue("maxchoices")
	statusStr := r.FormValue("status")
	optionsStr := r.FormValue("options")

	// 验证必填字段
	if title == "" || startTimeStr == "" || endTimeStr == "" || optionsStr == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 解析时间
	startTime, err := time.Parse("2006-01-02 15:04:05", startTimeStr)
	if err != nil {
		http.Error(w, "Invalid start time", http.StatusBadRequest)
		return
	}
	endTime, err := time.Parse("2006-01-02 15:04:05", endTimeStr)
	if err != nil {
		http.Error(w, "Invalid end time", http.StatusBadRequest)
		return
	}

	// 解析数值
	isMulti, _ := strconv.Atoi(isMultiStr)
	maxChoices, _ := strconv.Atoi(maxChoicesStr)
	status, _ := strconv.Atoi(statusStr)

	// 解析选项
	var optionTitles []string
	err = json.Unmarshal([]byte(optionsStr), &optionTitles)
	if err != nil {
		http.Error(w, "Invalid options format", http.StatusBadRequest)
		return
	}

	// 创建投票
	vote := &model.Vote{
		Title:       title,
		Description: description,
		StartTime:   startTime,
		EndTime:     endTime,
		IsMulti:     isMulti,
		MaxChoices:  maxChoices,
		Status:      status,
		TotalCount:  0,
	}

	// 创建选项
	options := make([]*model.VoteOption, 0, len(optionTitles))
	for _, optionTitle := range optionTitles {
		option := &model.VoteOption{
			Title: optionTitle,
			Count: 0,
		}
		options = append(options, option)
	}

	// 保存投票
	id, err := c.voteService.CreateVote(vote, options)
	if err != nil {
		logger.Error("创建投票失败", "error", err)
		http.Error(w, "Failed to create vote", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "投票创建成功",
			"id":      id,
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/vote_list", http.StatusFound)
	}
}

// Edit 编辑投票页面
func (c *VoteController) Edit(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取投票ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid vote ID", http.StatusBadRequest)
		return
	}

	// 获取投票
	vote, options, err := c.voteService.GetVote(id)
	if err != nil {
		logger.Error("获取投票失败", "id", id, "error", err)
		http.Error(w, "Vote not found", http.StatusNotFound)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Vote":        vote,
		"Options":     options,
		"CurrentMenu": "vote",
		"PageTitle":   "编辑投票",
	}

	// 渲染模板
	tplFile := "admin/vote_edit.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染编辑投票模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoEdit 处理编辑投票
func (c *VoteController) DoEdit(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取投票ID
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid vote ID", http.StatusBadRequest)
		return
	}

	// 获取原投票
	vote, _, err := c.voteService.GetVote(id)
	if err != nil {
		logger.Error("获取投票失败", "id", id, "error", err)
		http.Error(w, "Vote not found", http.StatusNotFound)
		return
	}

	// 获取表单数据
	title := r.FormValue("title")
	description := r.FormValue("description")
	startTimeStr := r.FormValue("starttime")
	endTimeStr := r.FormValue("endtime")
	isMultiStr := r.FormValue("ismulti")
	maxChoicesStr := r.FormValue("maxchoices")
	statusStr := r.FormValue("status")
	optionsStr := r.FormValue("options")

	// 验证必填字段
	if title == "" || startTimeStr == "" || endTimeStr == "" || optionsStr == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 解析时间
	startTime, err := time.Parse("2006-01-02 15:04:05", startTimeStr)
	if err != nil {
		http.Error(w, "Invalid start time", http.StatusBadRequest)
		return
	}
	endTime, err := time.Parse("2006-01-02 15:04:05", endTimeStr)
	if err != nil {
		http.Error(w, "Invalid end time", http.StatusBadRequest)
		return
	}

	// 解析数值
	isMulti, _ := strconv.Atoi(isMultiStr)
	maxChoices, _ := strconv.Atoi(maxChoicesStr)
	status, _ := strconv.Atoi(statusStr)

	// 解析选项
	var optionTitles []string
	err = json.Unmarshal([]byte(optionsStr), &optionTitles)
	if err != nil {
		http.Error(w, "Invalid options format", http.StatusBadRequest)
		return
	}

	// 更新投票
	vote.Title = title
	vote.Description = description
	vote.StartTime = startTime
	vote.EndTime = endTime
	vote.IsMulti = isMulti
	vote.MaxChoices = maxChoices
	vote.Status = status

	// 创建选项
	options := make([]*model.VoteOption, 0, len(optionTitles))
	for _, optionTitle := range optionTitles {
		option := &model.VoteOption{
			Title: optionTitle,
			Count: 0,
		}
		options = append(options, option)
	}

	// 保存投票
	err = c.voteService.UpdateVote(vote, options)
	if err != nil {
		logger.Error("更新投票失败", "error", err)
		http.Error(w, "Failed to update vote", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "投票更新成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/vote_list", http.StatusFound)
	}
}

// Delete 删除投票
func (c *VoteController) Delete(w http.ResponseWriter, r *http.Request) {
	// 获取投票ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid vote ID", http.StatusBadRequest)
		return
	}

	// 删除投票
	err = c.voteService.DeleteVote(id)
	if err != nil {
		logger.Error("删除投票失败", "id", id, "error", err)
		http.Error(w, "Failed to delete vote", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "投票删除成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/vote_list", http.StatusFound)
	}
}

// Result 投票结果
func (c *VoteController) Result(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取投票ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid vote ID", http.StatusBadRequest)
		return
	}

	// 获取投票结果
	vote, options, err := c.voteService.GetVoteResult(id)
	if err != nil {
		logger.Error("获取投票结果失败", "id", id, "error", err)
		http.Error(w, "Failed to get vote result", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Vote":        vote,
		"Options":     options,
		"CurrentMenu": "vote",
		"PageTitle":   "投票结果",
	}

	// 渲染模板
	tplFile := "admin/vote_result.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染投票结果模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Logs 投票日志
func (c *VoteController) Logs(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取投票ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid vote ID", http.StatusBadRequest)
		return
	}

	// 获取查询参数
	pageStr := r.URL.Query().Get("page")

	// 解析参数
	page := 1
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
		if page < 1 {
			page = 1
		}
	}

	// 获取投票
	vote, _, err := c.voteService.GetVote(id)
	if err != nil {
		logger.Error("获取投票失败", "id", id, "error", err)
		http.Error(w, "Vote not found", http.StatusNotFound)
		return
	}

	// 获取投票日志
	logs, total, err := c.voteService.GetVoteLogs(id, page, 20)
	if err != nil {
		logger.Error("获取投票日志失败", "error", err)
		http.Error(w, "Failed to get vote logs", http.StatusInternalServerError)
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
		"Vote":        vote,
		"Logs":        logs,
		"Pagination":  pagination,
		"CurrentMenu": "vote",
		"PageTitle":   "投票日志",
	}

	// 渲染模板
	tplFile := "admin/vote_logs.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染投票日志模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
