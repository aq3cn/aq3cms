package frontend

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
	"aq3cms/config"
	"aq3cms/internal/model"
	"aq3cms/internal/service"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
	"aq3cms/pkg/security"
)

// MemberController 会员控制器
type MemberController struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	memberModel     *model.MemberModel
	articleModel    *model.ArticleModel
	templateService *service.TemplateService
	sessionStore    *sessions.CookieStore
}

// NewMemberController 创建会员控制器
func NewMemberController(db *database.DB, cache cache.Cache, config *config.Config) *MemberController {
	return &MemberController{
		db:              db,
		cache:           cache,
		config:          config,
		memberModel:     model.NewMemberModel(db),
		articleModel:    model.NewArticleModel(db),
		templateService: service.NewTemplateService(db, cache, config),
		sessionStore:    sessions.NewCookieStore([]byte(config.Site.SessionSecret)),
	}
}

// Login 登录页面
func (c *MemberController) Login(w http.ResponseWriter, r *http.Request) {
	// 检查是否已登录
	session, _ := c.sessionStore.Get(r, "member-session")
	if session.Values["member_id"] != nil {
		http.Redirect(w, r, "/member/", http.StatusFound)
		return
	}

	// 获取全局变量
	globals := c.templateService.GetGlobals()

	// 准备模板数据
	data := map[string]interface{}{
		"Globals":     globals,
		"PageTitle":   "会员登录 - " + c.config.Site.Name,
		"Keywords":    c.config.Site.Keywords,
		"Description": c.config.Site.Description,
	}

	// 渲染模板
	tplFile := c.config.Template.DefaultTpl + "/member/login.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染登录模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoLogin 处理登录请求
func (c *MemberController) DoLogin(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	username := r.FormValue("username")
	password := r.FormValue("password")

	// 验证用户名和密码
	member, err := c.memberModel.CheckLogin(username, password)
	if err != nil {
		logger.Error("登录失败", "username", username, "error", err)
		http.Redirect(w, r, "/member/login?error=1", http.StatusFound)
		return
	}

	// 创建会话
	session, _ := c.sessionStore.Get(r, "member-session")
	session.Values["member_id"] = member.ID
	session.Values["member_name"] = member.Username
	session.Values["member_rank"] = 1 // 默认等级
	session.Save(r, w)

	// 更新登录信息
	go c.memberModel.UpdateLoginInfo(member.ID, r.RemoteAddr)

	// 重定向到会员中心
	http.Redirect(w, r, "/member/", http.StatusFound)
}

// Logout 退出登录
func (c *MemberController) Logout(w http.ResponseWriter, r *http.Request) {
	// 清除会话
	session, _ := c.sessionStore.Get(r, "member-session")
	session.Values = make(map[interface{}]interface{})
	session.Save(r, w)

	// 重定向到首页
	http.Redirect(w, r, "/", http.StatusFound)
}

// Register 注册页面
func (c *MemberController) Register(w http.ResponseWriter, r *http.Request) {
	// 检查是否已登录
	session, _ := c.sessionStore.Get(r, "member-session")
	if session.Values["member_id"] != nil {
		http.Redirect(w, r, "/member/", http.StatusFound)
		return
	}

	// 获取全局变量
	globals := c.templateService.GetGlobals()

	// 准备模板数据
	data := map[string]interface{}{
		"Globals":     globals,
		"PageTitle":   "会员注册 - " + c.config.Site.Name,
		"Keywords":    c.config.Site.Keywords,
		"Description": c.config.Site.Description,
	}

	// 渲染模板
	tplFile := c.config.Template.DefaultTpl + "/member/register.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染注册模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoRegister 处理注册请求
func (c *MemberController) DoRegister(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	// 检查用户名是否已存在
	member, err := c.memberModel.GetByUsername(username)
	exists := err == nil && member != nil
	if err != nil {
		logger.Error("检查用户名失败", "username", username, "error", err)
		http.Redirect(w, r, "/member/register?error=1", http.StatusFound)
		return
	}
	if exists {
		http.Redirect(w, r, "/member/register?error=2", http.StatusFound)
		return
	}

	// 检查邮箱是否已存在
	memberByEmail, err := c.memberModel.GetByEmail(email)
	exists = err == nil && memberByEmail != nil
	if err != nil {
		logger.Error("检查邮箱失败", "email", email, "error", err)
		http.Redirect(w, r, "/member/register?error=1", http.StatusFound)
		return
	}
	if exists {
		http.Redirect(w, r, "/member/register?error=3", http.StatusFound)
		return
	}

	// 创建会员
	newMember := &model.Member{
		Username:  username,
		Password:  security.HashPassword(password),
		Email:     email,
		RegTime:   time.Now(),
		RegIP:     r.RemoteAddr,
		LastLogin: time.Now(),
		LastIP:    r.RemoteAddr,
		Status:    1, // 正常状态
	}

	id, err := c.memberModel.Create(newMember)
	if err != nil {
		logger.Error("创建会员失败", "username", username, "error", err)
		http.Redirect(w, r, "/member/register?error=1", http.StatusFound)
		return
	}

	// 创建会话
	session, _ := c.sessionStore.Get(r, "member-session")
	session.Values["member_id"] = id
	session.Values["member_name"] = username
	session.Values["member_rank"] = 1 // 默认等级
	session.Save(r, w)

	// 重定向到会员中心
	http.Redirect(w, r, "/member/", http.StatusFound)
}

// Index 会员中心首页
func (c *MemberController) Index(w http.ResponseWriter, r *http.Request) {
	// 检查是否已登录
	session, _ := c.sessionStore.Get(r, "member-session")
	memberID, ok := session.Values["member_id"].(int64)
	if !ok {
		http.Redirect(w, r, "/member/login", http.StatusFound)
		return
	}

	// 获取会员信息
	member, err := c.memberModel.GetByID(memberID)
	if err != nil {
		logger.Error("获取会员信息失败", "id", memberID, "error", err)
		http.Redirect(w, r, "/member/login", http.StatusFound)
		return
	}

	// 获取会员文章
	articles, _, err := c.articleModel.GetMemberArticles(memberID, 1, 10)
	if err != nil {
		logger.Error("获取会员文章失败", "id", memberID, "error", err)
	}

	// 获取全局变量
	globals := c.templateService.GetGlobals()

	// 准备模板数据
	data := map[string]interface{}{
		"Globals":     globals,
		"Member":      member,
		"Articles":    articles,
		"PageTitle":   "会员中心 - " + c.config.Site.Name,
		"Keywords":    c.config.Site.Keywords,
		"Description": c.config.Site.Description,
	}

	// 渲染模板
	tplFile := c.config.Template.DefaultTpl + "/member/index.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染会员中心模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Profile 会员资料
func (c *MemberController) Profile(w http.ResponseWriter, r *http.Request) {
	// 检查是否已登录
	session, _ := c.sessionStore.Get(r, "member-session")
	memberID, ok := session.Values["member_id"].(int64)
	if !ok {
		http.Redirect(w, r, "/member/login", http.StatusFound)
		return
	}

	// 获取会员信息
	member, err := c.memberModel.GetByID(memberID)
	if err != nil {
		logger.Error("获取会员信息失败", "id", memberID, "error", err)
		http.Redirect(w, r, "/member/login", http.StatusFound)
		return
	}

	// 获取全局变量
	globals := c.templateService.GetGlobals()

	// 准备模板数据
	data := map[string]interface{}{
		"Globals":     globals,
		"Member":      member,
		"PageTitle":   "会员资料 - " + c.config.Site.Name,
		"Keywords":    c.config.Site.Keywords,
		"Description": c.config.Site.Description,
	}

	// 渲染模板
	tplFile := c.config.Template.DefaultTpl + "/member/profile.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染会员资料模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// UpdateProfile 更新会员资料
func (c *MemberController) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// 检查是否已登录
	session, _ := c.sessionStore.Get(r, "member-session")
	memberID, ok := session.Values["member_id"].(int64)
	if !ok {
		http.Redirect(w, r, "/member/login", http.StatusFound)
		return
	}

	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	email := r.FormValue("email")
	// 不使用的字段
	// nickname := r.FormValue("nickname")
	// sex := r.FormValue("sex") // 暂时不使用，避免数据库兼容性问题
	// birthday := r.FormValue("birthday")
	qq := r.FormValue("qq")
	// tel := r.FormValue("tel")

	// 获取会员信息
	member, err := c.memberModel.GetByID(memberID)
	if err != nil {
		logger.Error("获取会员信息失败", "id", memberID, "error", err)
		http.Redirect(w, r, "/member/profile?error=1", http.StatusFound)
		return
	}

	// 更新会员信息
	member.Email = email
	// member.NickName = nickname // 不存在的字段
	// member.Sex = sex // 暂时跳过，避免数据库兼容性问题
	// member.Birthday = birthday // 不存在的字段
	member.QQ = qq
	// member.Tel = tel // 不存在的字段

	if err := c.memberModel.Update(member); err != nil {
		logger.Error("更新会员信息失败", "id", memberID, "error", err)
		http.Redirect(w, r, "/member/profile?error=1", http.StatusFound)
		return
	}

	// 重定向到会员资料页
	http.Redirect(w, r, "/member/profile?success=1", http.StatusFound)
}

// ChangePassword 修改密码
func (c *MemberController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// 检查是否已登录
	session, _ := c.sessionStore.Get(r, "member-session")
	memberIDValue := session.Values["member_id"]
	_, ok := memberIDValue.(int64)
	if !ok {
		http.Redirect(w, r, "/member/login", http.StatusFound)
		return
	}

	// 获取全局变量
	globals := c.templateService.GetGlobals()

	// 准备模板数据
	data := map[string]interface{}{
		"Globals":     globals,
		"PageTitle":   "修改密码 - " + c.config.Site.Name,
		"Keywords":    c.config.Site.Keywords,
		"Description": c.config.Site.Description,
	}

	// 渲染模板
	tplFile := c.config.Template.DefaultTpl + "/member/password.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染修改密码模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoChangePassword 处理修改密码请求
func (c *MemberController) DoChangePassword(w http.ResponseWriter, r *http.Request) {
	// 检查是否已登录
	session, _ := c.sessionStore.Get(r, "member-session")
	memberID, ok := session.Values["member_id"].(int64)
	if !ok {
		http.Redirect(w, r, "/member/login", http.StatusFound)
		return
	}

	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	oldPassword := r.FormValue("old_password")
	newPassword := r.FormValue("new_password")
	confirmPassword := r.FormValue("confirm_password")

	// 检查新密码和确认密码是否一致
	if newPassword != confirmPassword {
		http.Redirect(w, r, "/member/password?error=1", http.StatusFound)
		return
	}

	// 验证旧密码
	member, err := c.memberModel.GetByID(memberID)
	if err != nil {
		logger.Error("获取会员信息失败", "id", memberID, "error", err)
		http.Redirect(w, r, "/member/password?error=2", http.StatusFound)
		return
	}

	if !security.CheckPassword(oldPassword, member.Password) {
		http.Redirect(w, r, "/member/password?error=3", http.StatusFound)
		return
	}

	// 更新密码
	member.Password = security.HashPassword(newPassword)
	if err := c.memberModel.Update(member); err != nil {
		logger.Error("更新密码失败", "id", memberID, "error", err)
		http.Redirect(w, r, "/member/password?error=2", http.StatusFound)
		return
	}

	// 重定向到会员中心
	http.Redirect(w, r, "/member/password?success=1", http.StatusFound)
}

// Articles 会员文章列表
func (c *MemberController) Articles(w http.ResponseWriter, r *http.Request) {
	// 检查是否已登录
	session, _ := c.sessionStore.Get(r, "member-session")
	memberID, ok := session.Values["member_id"].(int64)
	if !ok {
		http.Redirect(w, r, "/member/login", http.StatusFound)
		return
	}

	// 获取页码
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// 获取会员文章
	pageSize := 10
	articles, total, err := c.articleModel.GetMemberArticles(memberID, page, pageSize)
	if err != nil {
		logger.Error("获取会员文章失败", "id", memberID, "error", err)
		http.Error(w, "Failed to get member articles", http.StatusInternalServerError)
		return
	}

	// 计算分页信息
	totalPages := (total + pageSize - 1) / pageSize
	pagination := map[string]interface{}{
		"CurrentPage": page,
		"TotalPages":  totalPages,
		"TotalItems":  total,
		"HasPrev":     page > 1,
		"HasNext":     page < totalPages,
		"PrevPage":    page - 1,
		"NextPage":    page + 1,
	}

	// 获取全局变量
	globals := c.templateService.GetGlobals()

	// 准备模板数据
	data := map[string]interface{}{
		"Globals":     globals,
		"Articles":    articles,
		"Pagination":  pagination,
		"PageTitle":   "我的文章 - " + c.config.Site.Name,
		"Keywords":    c.config.Site.Keywords,
		"Description": c.config.Site.Description,
	}

	// 渲染模板
	tplFile := c.config.Template.DefaultTpl + "/member/articles.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染会员文章模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
