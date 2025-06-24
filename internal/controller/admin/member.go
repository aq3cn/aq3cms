package admin

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"aq3cms/config"
	"aq3cms/internal/middleware"
	"aq3cms/internal/model"
	"aq3cms/internal/service"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
	"aq3cms/pkg/security"
	"github.com/gorilla/mux"
)

// MemberController 会员控制器
type MemberController struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	memberModel     *model.MemberModel
	memberTypeModel *model.MemberTypeModel
	templateService *service.TemplateService
}

// NewMemberController 创建会员控制器
func NewMemberController(db *database.DB, cache cache.Cache, config *config.Config) *MemberController {
	return &MemberController{
		db:              db,
		cache:           cache,
		config:          config,
		memberModel:     model.NewMemberModel(db),
		memberTypeModel: model.NewMemberTypeModel(db),
		templateService: service.NewTemplateService(db, cache, config),
	}
}

// Index 会员管理首页
func (c *MemberController) Index(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取会员统计信息
	totalMembers, err := c.memberModel.GetTotalCount()
	if err != nil {
		logger.Error("获取会员总数失败", "error", err)
		totalMembers = 0
	}

	// 获取今日新增会员数
	todayMembers, err := c.memberModel.GetTodayCount()
	if err != nil {
		logger.Error("获取今日新增会员数失败", "error", err)
		todayMembers = 0
	}

	// 获取活跃会员数（最近30天登录）
	activeMembers, err := c.memberModel.GetActiveCount(30)
	if err != nil {
		logger.Error("获取活跃会员数失败", "error", err)
		activeMembers = 0
	}

	// 获取禁用会员数
	disabledMembers, err := c.memberModel.GetDisabledCount()
	if err != nil {
		logger.Error("获取禁用会员数失败", "error", err)
		disabledMembers = 0
	}

	// 获取最新注册的会员
	latestMembers, err := c.memberModel.GetLatest(10)
	if err != nil {
		logger.Error("获取最新会员失败", "error", err)
		latestMembers = []*model.Member{}
	}

	// 获取会员类型统计
	memberTypeStats, err := c.memberModel.GetTypeStats()
	if err != nil {
		logger.Error("获取会员类型统计失败", "error", err)
		memberTypeStats = []map[string]interface{}{}
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":         adminID,
		"AdminName":       adminName,
		"TotalMembers":    totalMembers,
		"TodayMembers":    todayMembers,
		"ActiveMembers":   activeMembers,
		"DisabledMembers": disabledMembers,
		"LatestMembers":   latestMembers,
		"MemberTypeStats": memberTypeStats,
		"CurrentMenu":     "member",
		"PageTitle":       "会员管理",
	}

	// 渲染模板
	tplFile := "admin/member_index.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染会员管理首页模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// List 会员列表
func (c *MemberController) List(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取查询参数
	mtypeStr := r.URL.Query().Get("mtype")
	keyword := r.URL.Query().Get("keyword")
	sort := r.URL.Query().Get("sort")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pagesize")

	// 解析参数
	mtype := 0
	if mtypeStr != "" {
		mtype, _ = strconv.Atoi(mtypeStr)
	}

	page := 1
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
		if page < 1 {
			page = 1
		}
	}

	pageSize := 20
	if pageSizeStr != "" {
		pageSize, _ = strconv.Atoi(pageSizeStr)
		if pageSize < 1 {
			pageSize = 20
		}
	}

	// 获取会员列表
	var members []*model.Member
	var total int
	var err error

	if keyword != "" {
		// 搜索会员
		members, total, err = c.memberModel.Search(keyword, page, pageSize)
	} else {
		// 获取会员列表
		members, total, err = c.memberModel.GetList(page, pageSize)
	}

	if err != nil {
		logger.Error("获取会员列表失败", "error", err)
		http.Error(w, "Failed to get members", http.StatusInternalServerError)
		return
	}

	// 获取会员类型列表
	memberTypes, err := c.memberTypeModel.GetAll()
	if err != nil {
		logger.Error("获取会员类型列表失败", "error", err)
	}

	// 计算分页信息
	totalPages := (total + pageSize - 1) / pageSize

	// 生成页码数组
	pageNumbers := make([]int, 0)
	start := page - 2
	if start < 1 {
		start = 1
	}
	end := start + 4
	if end > totalPages {
		end = totalPages
		start = end - 4
		if start < 1 {
			start = 1
		}
	}
	for i := start; i <= end; i++ {
		pageNumbers = append(pageNumbers, i)
	}

	pagination := map[string]interface{}{
		"Page":        page,
		"TotalPages":  totalPages,
		"Total":       total,
		"HasPrev":     page > 1,
		"HasNext":     page < totalPages,
		"PrevPage":    page - 1,
		"NextPage":    page + 1,
		"PageNumbers": pageNumbers,
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Members":     members,
		"MemberTypes": memberTypes,
		"Pagination":  pagination,
		"MType":       mtype,
		"Keyword":     keyword,
		"Sort":        sort,
		"CurrentMenu": "member",
		"PageTitle":   "会员管理",
	}

	// 渲染模板
	tplFile := "admin/member_list.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染会员列表模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Add 添加会员页面
func (c *MemberController) Add(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取会员类型列表
	memberTypes, err := c.memberTypeModel.GetAll()
	if err != nil {
		logger.Error("获取会员类型列表失败", "error", err)
		http.Error(w, "Failed to get member types", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"MemberTypes": memberTypes,
		"CurrentMenu": "member",
		"PageTitle":   "添加会员",
	}

	// 渲染模板
	tplFile := "admin/member_add.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染添加会员模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoAdd 处理添加会员
func (c *MemberController) DoAdd(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	mtypeStr := r.FormValue("mtype")
	sex := r.FormValue("sex")
	mobile := r.FormValue("mobile")
	qq := r.FormValue("qq")
	statusStr := r.FormValue("status")
	moneyStr := r.FormValue("money")
	scoreStr := r.FormValue("score")

	// 验证必填字段
	if username == "" || password == "" || email == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 检查用户名是否已存在
	existingMember, err := c.memberModel.GetByUsername(username)
	if err == nil && existingMember != nil {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}

	// 检查邮箱是否已存在
	existingMember, err = c.memberModel.GetByEmail(email)
	if err == nil && existingMember != nil {
		http.Error(w, "Email already exists", http.StatusBadRequest)
		return
	}

	// 解析字段
	mtype := 1
	if mtypeStr != "" {
		mtype, _ = strconv.Atoi(mtypeStr)
	}

	status := 1
	if statusStr != "" {
		status, _ = strconv.Atoi(statusStr)
	}

	money := 0.0
	if moneyStr != "" {
		money, _ = strconv.ParseFloat(moneyStr, 64)
	}

	score := 0
	if scoreStr != "" {
		score, _ = strconv.Atoi(scoreStr)
	}

	// 处理性别字段，如果为空则设置为默认值
	if sex == "" {
		sex = "保密"
	}

	// 创建会员
	member := &model.Member{
		Username:  username,
		Password:  security.HashPassword(password),
		Email:     email,
		MType:     mtype,
		Sex:       sex,
		Mobile:    mobile,
		QQ:        qq,
		Status:    status,
		RegTime:   time.Now(),
		RegIP:     r.RemoteAddr,
		LastLogin: time.Time{},
		LastIP:    "",
		Money:     money,
		Score:     score,
		Face:      "", // 默认头像为空
	}

	// 保存会员
	id, err := c.memberModel.Create(member)
	if err != nil {
		logger.Error("创建会员失败", "error", err)
		http.Error(w, "Failed to create member", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "会员创建成功",
			"id":      id,
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/aq3cms/member_list", http.StatusFound)
	}
}

// Edit 编辑会员页面
func (c *MemberController) Edit(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取会员ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid member ID", http.StatusBadRequest)
		return
	}

	// 获取会员
	member, err := c.memberModel.GetByID(id)
	if err != nil {
		logger.Error("获取会员失败", "id", id, "error", err)
		http.Error(w, "Member not found", http.StatusNotFound)
		return
	}

	// 获取会员类型列表
	memberTypes, err := c.memberTypeModel.GetAll()
	if err != nil {
		logger.Error("获取会员类型列表失败", "error", err)
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Member":      member,
		"MemberTypes": memberTypes,
		"CurrentMenu": "member",
		"PageTitle":   "编辑会员",
	}

	// 渲染模板
	tplFile := "admin/member_edit.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染编辑会员模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoEdit 处理编辑会员
func (c *MemberController) DoEdit(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取会员ID
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid member ID", http.StatusBadRequest)
		return
	}

	// 获取原会员
	member, err := c.memberModel.GetByID(id)
	if err != nil {
		logger.Error("获取会员失败", "id", id, "error", err)
		http.Error(w, "Member not found", http.StatusNotFound)
		return
	}

	// 获取表单数据
	email := r.FormValue("email")
	mtypeStr := r.FormValue("mtype")
	sex := r.FormValue("sex")
	mobile := r.FormValue("mobile")
	qq := r.FormValue("qq")
	statusStr := r.FormValue("status")
	password := r.FormValue("password")
	moneyStr := r.FormValue("money")
	scoreStr := r.FormValue("score")

	// 验证必填字段
	if email == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 检查邮箱是否已存在
	if email != member.Email {
		existingMember, err := c.memberModel.GetByEmail(email)
		if err == nil && existingMember != nil && existingMember.ID != id {
			http.Error(w, "Email already exists", http.StatusBadRequest)
			return
		}
	}

	// 解析字段
	mtype := 0
	if mtypeStr != "" {
		mtype, _ = strconv.Atoi(mtypeStr)
	}

	status := 1
	if statusStr != "" {
		status, _ = strconv.Atoi(statusStr)
	}

	money := member.Money
	if moneyStr != "" {
		moneyFloat, err := strconv.ParseFloat(moneyStr, 64)
		if err == nil {
			money = moneyFloat
		}
	}

	score := member.Score
	if scoreStr != "" {
		scoreInt, err := strconv.Atoi(scoreStr)
		if err == nil {
			score = scoreInt
		}
	}

	// 处理性别字段，如果为空则设置为默认值
	if sex == "" {
		sex = "保密"
	}

	// 更新会员
	member.Email = email
	member.MType = mtype
	member.Sex = sex
	member.Mobile = mobile
	member.QQ = qq
	member.Status = status
	member.Money = money
	member.Score = score

	// 如果提供了新密码，则更新密码
	if password != "" {
		member.Password = security.HashPassword(password)
	}

	// 保存会员
	err = c.memberModel.Update(member)
	if err != nil {
		logger.Error("更新会员失败", "error", err)
		http.Error(w, "Failed to update member", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "会员更新成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/aq3cms/member_list", http.StatusFound)
	}
}

// Delete 删除会员
func (c *MemberController) Delete(w http.ResponseWriter, r *http.Request) {
	// 获取会员ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid member ID", http.StatusBadRequest)
		return
	}

	// 删除会员
	err = c.memberModel.Delete(id)
	if err != nil {
		logger.Error("删除会员失败", "id", id, "error", err)
		http.Error(w, "Failed to delete member", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "会员删除成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/member/list", http.StatusFound)
	}
}

// BatchDelete 批量删除会员
func (c *MemberController) BatchDelete(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取会员ID列表
	idsStr := r.FormValue("ids")
	if idsStr == "" {
		http.Error(w, "No members selected", http.StatusBadRequest)
		return
	}

	// 分割ID列表
	idStrs := strings.Split(idsStr, ",")
	ids := make([]int64, 0, len(idStrs))
	for _, idStr := range idStrs {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}

	// 批量删除会员
	for _, id := range ids {
		err := c.memberModel.Delete(id)
		if err != nil {
			logger.Error("删除会员失败", "id", id, "error", err)
		}
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "会员批量删除成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/aq3cms/member_list", http.StatusFound)
	}
}
