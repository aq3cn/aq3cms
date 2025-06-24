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

// ScoreController 积分控制器
type ScoreController struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	scoreService    *service.ScoreService
	scoreRuleModel  *model.ScoreRuleModel
	scoreLogModel   *model.ScoreLogModel
	memberModel     *model.MemberModel
	templateService *service.TemplateService
}

// NewScoreController 创建积分控制器
func NewScoreController(db *database.DB, cache cache.Cache, config *config.Config) *ScoreController {
	return &ScoreController{
		db:              db,
		cache:           cache,
		config:          config,
		scoreService:    service.NewScoreService(db, cache, config),
		scoreRuleModel:  model.NewScoreRuleModel(db),
		scoreLogModel:   model.NewScoreLogModel(db),
		memberModel:     model.NewMemberModel(db),
		templateService: service.NewTemplateService(db, cache, config),
	}
}

// RuleList 积分规则列表
func (c *ScoreController) RuleList(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取积分规则列表
	rules, err := c.scoreRuleModel.GetAll(-1)
	if err != nil {
		logger.Error("获取积分规则列表失败", "error", err)
		http.Error(w, "Failed to get score rules", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Rules":       rules,
		"CurrentMenu": "score",
		"PageTitle":   "积分规则管理",
	}

	// 渲染模板
	tplFile := "admin/score_rule_list.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染积分规则列表模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// RuleAdd 添加积分规则页面
func (c *ScoreController) RuleAdd(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"CurrentMenu": "score",
		"PageTitle":   "添加积分规则",
	}

	// 渲染模板
	tplFile := "admin/score_rule_add.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染添加积分规则模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// RuleDoAdd 处理添加积分规则
func (c *ScoreController) RuleDoAdd(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	name := r.FormValue("name")
	code := r.FormValue("code")
	scoreStr := r.FormValue("score")
	maxTimesStr := r.FormValue("maxtimes")
	cycleTypeStr := r.FormValue("cycletype")
	statusStr := r.FormValue("status")
	description := r.FormValue("description")

	// 验证必填字段
	if name == "" || code == "" || scoreStr == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 解析数值
	score, _ := strconv.Atoi(scoreStr)
	maxTimes, _ := strconv.Atoi(maxTimesStr)
	cycleType, _ := strconv.Atoi(cycleTypeStr)
	status, _ := strconv.Atoi(statusStr)

	// 创建积分规则
	rule := &model.ScoreRule{
		Name:        name,
		Code:        code,
		Score:       score,
		MaxTimes:    maxTimes,
		CycleType:   cycleType,
		Status:      status,
		Description: description,
	}

	// 保存积分规则
	id, err := c.scoreRuleModel.Create(rule)
	if err != nil {
		logger.Error("创建积分规则失败", "error", err)
		http.Error(w, "Failed to create score rule", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "积分规则创建成功",
			"id":      id,
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/score_rule_list", http.StatusFound)
	}
}

// RuleEdit 编辑积分规则页面
func (c *ScoreController) RuleEdit(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取积分规则ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid score rule ID", http.StatusBadRequest)
		return
	}

	// 获取积分规则
	rule, err := c.scoreRuleModel.GetByID(id)
	if err != nil {
		logger.Error("获取积分规则失败", "id", id, "error", err)
		http.Error(w, "Score rule not found", http.StatusNotFound)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Rule":        rule,
		"CurrentMenu": "score",
		"PageTitle":   "编辑积分规则",
	}

	// 渲染模板
	tplFile := "admin/score_rule_edit.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染编辑积分规则模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// RuleDoEdit 处理编辑积分规则
func (c *ScoreController) RuleDoEdit(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取积分规则ID
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid score rule ID", http.StatusBadRequest)
		return
	}

	// 获取原积分规则
	rule, err := c.scoreRuleModel.GetByID(id)
	if err != nil {
		logger.Error("获取积分规则失败", "id", id, "error", err)
		http.Error(w, "Score rule not found", http.StatusNotFound)
		return
	}

	// 获取表单数据
	name := r.FormValue("name")
	code := r.FormValue("code")
	scoreStr := r.FormValue("score")
	maxTimesStr := r.FormValue("maxtimes")
	cycleTypeStr := r.FormValue("cycletype")
	statusStr := r.FormValue("status")
	description := r.FormValue("description")

	// 验证必填字段
	if name == "" || code == "" || scoreStr == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 解析数值
	score, _ := strconv.Atoi(scoreStr)
	maxTimes, _ := strconv.Atoi(maxTimesStr)
	cycleType, _ := strconv.Atoi(cycleTypeStr)
	status, _ := strconv.Atoi(statusStr)

	// 更新积分规则
	rule.Name = name
	rule.Code = code
	rule.Score = score
	rule.MaxTimes = maxTimes
	rule.CycleType = cycleType
	rule.Status = status
	rule.Description = description

	// 保存积分规则
	err = c.scoreRuleModel.Update(rule)
	if err != nil {
		logger.Error("更新积分规则失败", "error", err)
		http.Error(w, "Failed to update score rule", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "积分规则更新成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/score_rule_list", http.StatusFound)
	}
}

// RuleDelete 删除积分规则
func (c *ScoreController) RuleDelete(w http.ResponseWriter, r *http.Request) {
	// 获取积分规则ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid score rule ID", http.StatusBadRequest)
		return
	}

	// 删除积分规则
	err = c.scoreRuleModel.Delete(id)
	if err != nil {
		logger.Error("删除积分规则失败", "id", id, "error", err)
		http.Error(w, "Failed to delete score rule", http.StatusInternalServerError)
		return
	}

	// 删除积分日志
	err = c.scoreLogModel.DeleteByRuleID(id)
	if err != nil {
		logger.Error("删除积分日志失败", "ruleid", id, "error", err)
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "积分规则删除成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/score_rule_list", http.StatusFound)
	}
}

// LogList 积分日志列表
func (c *ScoreController) LogList(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取查询参数
	memberIDStr := r.URL.Query().Get("memberid")
	ruleIDStr := r.URL.Query().Get("ruleid")
	pageStr := r.URL.Query().Get("page")

	// 解析参数
	memberID := int64(0)
	if memberIDStr != "" {
		memberID, _ = strconv.ParseInt(memberIDStr, 10, 64)
	}
	ruleID := int64(0)
	if ruleIDStr != "" {
		ruleID, _ = strconv.ParseInt(ruleIDStr, 10, 64)
	}
	page := 1
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
		if page < 1 {
			page = 1
		}
	}

	// 获取积分日志
	var logs []*model.ScoreLog
	var total int
	var err error
	if memberID > 0 {
		logs, total, err = c.scoreLogModel.GetByMemberID(memberID, page, 20)
	} else if ruleID > 0 {
		logs, total, err = c.scoreLogModel.GetByRuleID(ruleID, page, 20)
	} else {
		// 获取所有积分日志
		// 这里需要实现一个获取所有积分日志的方法
		// 暂时使用获取会员ID为0的日志代替
		logs, total, err = c.scoreLogModel.GetByMemberID(0, page, 20)
	}
	if err != nil {
		logger.Error("获取积分日志失败", "error", err)
		http.Error(w, "Failed to get score logs", http.StatusInternalServerError)
		return
	}

	// 获取积分规则列表
	rules, err := c.scoreRuleModel.GetAll(-1)
	if err != nil {
		logger.Error("获取积分规则列表失败", "error", err)
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
		"Logs":        logs,
		"Rules":       rules,
		"MemberID":    memberID,
		"RuleID":      ruleID,
		"Pagination":  pagination,
		"CurrentMenu": "score",
		"PageTitle":   "积分日志管理",
	}

	// 渲染模板
	tplFile := "admin/score_log_list.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染积分日志列表模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// LogDelete 删除积分日志
func (c *ScoreController) LogDelete(w http.ResponseWriter, r *http.Request) {
	// 获取积分日志ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid score log ID", http.StatusBadRequest)
		return
	}

	// 获取积分日志
	log, err := c.scoreLogModel.GetByID(id)
	if err != nil {
		logger.Error("获取积分日志失败", "id", id, "error", err)
		http.Error(w, "Score log not found", http.StatusNotFound)
		return
	}

	// 更新会员积分
	err = c.memberModel.UpdateScore(log.MemberID, -log.Score)
	if err != nil {
		logger.Error("更新会员积分失败", "memberid", log.MemberID, "error", err)
	}

	// 删除积分日志
	err = c.scoreLogModel.Delete(id)
	if err != nil {
		logger.Error("删除积分日志失败", "id", id, "error", err)
		http.Error(w, "Failed to delete score log", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "积分日志删除成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/score_log_list", http.StatusFound)
	}
}

// AddScore 添加积分页面
func (c *ScoreController) AddScore(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取会员ID
	vars := mux.Vars(r)
	memberIDStr := vars["memberid"]
	memberID, err := strconv.ParseInt(memberIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid member ID", http.StatusBadRequest)
		return
	}

	// 获取会员
	member, err := c.memberModel.GetByID(memberID)
	if err != nil {
		logger.Error("获取会员失败", "id", memberID, "error", err)
		http.Error(w, "Member not found", http.StatusNotFound)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Member":      member,
		"CurrentMenu": "member",
		"PageTitle":   "添加积分",
	}

	// 渲染模板
	tplFile := "admin/score_add.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染添加积分模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoAddScore 处理添加积分
func (c *ScoreController) DoAddScore(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	memberIDStr := r.FormValue("memberid")
	scoreStr := r.FormValue("score")
	remark := r.FormValue("remark")

	// 验证必填字段
	if memberIDStr == "" || scoreStr == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 解析数值
	memberID, _ := strconv.ParseInt(memberIDStr, 10, 64)
	score, _ := strconv.Atoi(scoreStr)

	// 检查会员是否存在
	_, err := c.memberModel.GetByID(memberID)
	if err != nil {
		logger.Error("获取会员失败", "id", memberID, "error", err)
		http.Error(w, "Member not found", http.StatusNotFound)
		return
	}

	// 创建积分日志
	log := &model.ScoreLog{
		MemberID: memberID,
		RuleID:   0, // 管理员手动添加
		Score:    score,
		Remark:   remark,
		IP:       r.RemoteAddr,
	}
	_, err = c.scoreLogModel.Create(log)
	if err != nil {
		logger.Error("创建积分日志失败", "error", err)
		http.Error(w, "Failed to create score log", http.StatusInternalServerError)
		return
	}

	// 更新会员积分
	err = c.memberModel.UpdateScore(memberID, score)
	if err != nil {
		logger.Error("更新会员积分失败", "error", err)
		http.Error(w, "Failed to update member score", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "积分添加成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/member_edit/"+memberIDStr, http.StatusFound)
	}
}

// InitRules 初始化积分规则
func (c *ScoreController) InitRules(w http.ResponseWriter, r *http.Request) {
	// 初始化默认规则
	err := c.scoreService.InitDefaultRules()
	if err != nil {
		logger.Error("初始化积分规则失败", "error", err)
		http.Error(w, "Failed to initialize score rules", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "积分规则初始化成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/score_rule_list", http.StatusFound)
	}
}
