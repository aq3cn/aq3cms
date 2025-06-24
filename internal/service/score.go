package service

import (
	"fmt"
	"time"

	"aq3cms/config"
	"aq3cms/internal/model"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// ScoreService 积分服务
type ScoreService struct {
	db             *database.DB
	cache          cache.Cache
	config         *config.Config
	scoreRuleModel *model.ScoreRuleModel
	scoreLogModel  *model.ScoreLogModel
	memberModel    *model.MemberModel
}

// NewScoreService 创建积分服务
func NewScoreService(db *database.DB, cache cache.Cache, config *config.Config) *ScoreService {
	return &ScoreService{
		db:             db,
		cache:          cache,
		config:         config,
		scoreRuleModel: model.NewScoreRuleModel(db),
		scoreLogModel:  model.NewScoreLogModel(db),
		memberModel:    model.NewMemberModel(db),
	}
}

// AddScore 添加积分
func (s *ScoreService) AddScore(memberID int64, ruleCode string, ip string, remark string) (int, error) {
	// 获取积分规则
	rule, err := s.scoreRuleModel.GetByCode(ruleCode)
	if err != nil {
		return 0, err
	}

	// 检查规则状态
	if rule.Status != 1 {
		return 0, fmt.Errorf("score rule is disabled")
	}

	// 检查次数限制
	if rule.MaxTimes > 0 {
		// 根据周期类型检查次数
		var count int
		switch rule.CycleType {
		case 1: // 每天
			count, err = s.scoreLogModel.GetTodayCount(memberID, rule.ID)
		case 2: // 每周
			count, err = s.scoreLogModel.GetWeekCount(memberID, rule.ID)
		case 3: // 每月
			count, err = s.scoreLogModel.GetMonthCount(memberID, rule.ID)
		case 4: // 每年
			count, err = s.scoreLogModel.GetYearCount(memberID, rule.ID)
		default: // 不限
			count, err = s.scoreLogModel.GetTotalCount(memberID, rule.ID)
		}
		if err != nil {
			return 0, err
		}

		// 检查是否超过最大次数
		if count >= rule.MaxTimes {
			return 0, fmt.Errorf("score rule max times exceeded")
		}
	}

	// 创建积分日志
	log := &model.ScoreLog{
		MemberID: memberID,
		RuleID:   rule.ID,
		Score:    rule.Score,
		Remark:   remark,
		IP:       ip,
	}
	_, err = s.scoreLogModel.Create(log)
	if err != nil {
		return 0, err
	}

	// 更新会员积分
	err = s.memberModel.UpdateScore(memberID, rule.Score)
	if err != nil {
		return 0, err
	}

	return rule.Score, nil
}

// GetScoreLogs 获取积分日志
func (s *ScoreService) GetScoreLogs(memberID int64, page, pageSize int) ([]*model.ScoreLog, int, error) {
	// 获取积分日志
	logs, total, err := s.scoreLogModel.GetByMemberID(memberID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetScoreRules 获取积分规则
func (s *ScoreService) GetScoreRules() ([]*model.ScoreRule, error) {
	// 获取积分规则
	rules, err := s.scoreRuleModel.GetAll(1)
	if err != nil {
		return nil, err
	}

	return rules, nil
}

// GetMemberScore 获取会员积分
func (s *ScoreService) GetMemberScore(memberID int64) (int, error) {
	// 获取会员
	member, err := s.memberModel.GetByID(memberID)
	if err != nil {
		return 0, err
	}

	return member.Score, nil
}

// GetScoreRanking 获取积分排行
func (s *ScoreService) GetScoreRanking(limit int) ([]*model.Member, error) {
	// 构建查询
	qb := database.NewQueryBuilder(s.db, "member")
	qb.OrderBy("score DESC")
	qb.Limit(limit)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取积分排行失败", "error", err)
		return nil, err
	}

	// 转换为会员列表
	members := make([]*model.Member, 0, len(results))
	for _, result := range results {
		member := &model.Member{}
		member.ID, _ = result["id"].(int64)
		member.Username, _ = result["username"].(string)
		member.Email, _ = result["email"].(string)
		member.Mobile, _ = result["mobile"].(string)
		member.MType, _ = result["mtype"].(int)
		member.Sex, _ = result["sex"].(string)
		member.Avatar, _ = result["avatar"].(string)
		member.QQ, _ = result["qq"].(string)
		member.Score, _ = result["score"].(int)
		member.Money, _ = result["money"].(float64)
		member.Status, _ = result["status"].(int)
		member.RegTime, _ = result["regtime"].(time.Time)
		member.RegIP, _ = result["regip"].(string)
		member.LastLogin, _ = result["lastlogin"].(time.Time)
		member.LastIP, _ = result["lastip"].(string)
		member.LoginCount, _ = result["logincount"].(int)
		members = append(members, member)
	}

	return members, nil
}

// InitDefaultRules 初始化默认规则
func (s *ScoreService) InitDefaultRules() error {
	// 默认规则
	defaultRules := []*model.ScoreRule{
		{
			Name:        "注册",
			Code:        "register",
			Score:       50,
			MaxTimes:    1,
			CycleType:   0,
			Status:      1,
			Description: "注册成为会员",
		},
		{
			Name:        "登录",
			Code:        "login",
			Score:       2,
			MaxTimes:    1,
			CycleType:   1,
			Status:      1,
			Description: "每天登录",
		},
		{
			Name:        "发表文章",
			Code:        "post_article",
			Score:       5,
			MaxTimes:    10,
			CycleType:   1,
			Status:      1,
			Description: "发表文章",
		},
		{
			Name:        "评论",
			Code:        "post_comment",
			Score:       2,
			MaxTimes:    10,
			CycleType:   1,
			Status:      1,
			Description: "发表评论",
		},
		{
			Name:        "完善资料",
			Code:        "complete_profile",
			Score:       20,
			MaxTimes:    1,
			CycleType:   0,
			Status:      1,
			Description: "完善个人资料",
		},
	}

	// 创建默认规则
	for _, rule := range defaultRules {
		// 检查规则是否已存在
		existingRule, err := s.scoreRuleModel.GetByCode(rule.Code)
		if err == nil && existingRule != nil {
			// 规则已存在，跳过
			continue
		}

		// 创建规则
		_, err = s.scoreRuleModel.Create(rule)
		if err != nil {
			logger.Error("创建默认积分规则失败", "code", rule.Code, "error", err)
			return err
		}
	}

	return nil
}
