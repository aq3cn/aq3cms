package model

import (
	"fmt"
	"strconv"
	"time"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
	"aq3cms/pkg/security"
)

// Member 会员
type Member struct {
	ID         int64     `json:"id"`
	Username   string    `json:"username"`   // 用户名
	Password   string    `json:"password"`   // 密码
	Email      string    `json:"email"`      // 邮箱
	Mobile     string    `json:"mobile"`     // 手机
	MType      int       `json:"mtype"`      // 会员类型
	Sex        string    `json:"sex"`        // 性别
	Avatar     string    `json:"avatar"`     // 头像
	Face       string    `json:"face"`       // 头像（兼容旧版）
	QQ         string    `json:"qq"`         // QQ
	Score      int       `json:"score"`      // 积分
	Money      float64   `json:"money"`      // 余额
	Status     int       `json:"status"`     // 状态
	RegTime    time.Time `json:"regtime"`    // 注册时间
	RegIP      string    `json:"regip"`      // 注册IP
	LastLogin  time.Time `json:"lastlogin"`  // 最后登录时间
	LastIP     string    `json:"lastip"`     // 最后登录IP
	LoginCount int       `json:"logincount"` // 登录次数
}

// MemberModel 会员模型
type MemberModel struct {
	db *database.DB
}

// NewMemberModel 创建会员模型
func NewMemberModel(db *database.DB) *MemberModel {
	return &MemberModel{
		db: db,
	}
}

// GetByID 根据ID获取会员
func (m *MemberModel) GetByID(id int64) (*Member, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "member")
	qb.Where("mid = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取会员失败", "id", id, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("member not found: %d", id)
	}

	// 转换为会员
	member := &Member{}

	// 处理ID字段
	if mid, ok := result["mid"].(int64); ok {
		member.ID = mid
	} else if mid, ok := result["mid"].([]byte); ok {
		if idStr := string(mid); idStr != "" {
			if idInt, err := strconv.ParseInt(idStr, 10, 64); err == nil {
				member.ID = idInt
			}
		}
	}

	member.Username, _ = result["userid"].(string)
	member.Password, _ = result["pwd"].(string)
	member.Email, _ = result["email"].(string)
	member.MType = convertToInt(result["mtype"])
	member.Sex, _ = result["sex"].(string)
	member.Face, _ = result["face"].(string)
	member.Avatar = member.Face // 兼容性
	member.Score = convertToInt(result["scores"])
	member.Money = convertToFloat64(result["money"])

	// 处理时间戳字段
	if jointime, ok := result["jointime"].(int64); ok {
		member.RegTime = time.Unix(jointime, 0)
	} else if jointime, ok := result["jointime"].([]byte); ok {
		if timeStr := string(jointime); timeStr != "" {
			if timestamp, err := strconv.ParseInt(timeStr, 10, 64); err == nil {
				member.RegTime = time.Unix(timestamp, 0)
			}
		}
	}

	if logintime, ok := result["logintime"].(int64); ok {
		member.LastLogin = time.Unix(logintime, 0)
	} else if logintime, ok := result["logintime"].([]byte); ok {
		if timeStr := string(logintime); timeStr != "" {
			if timestamp, err := strconv.ParseInt(timeStr, 10, 64); err == nil {
				member.LastLogin = time.Unix(timestamp, 0)
			}
		}
	}

	member.RegIP, _ = result["joinip"].(string)
	member.LastIP, _ = result["loginip"].(string)
	member.Status = 1 // 默认状态为正常

	return member, nil
}

// GetCount 获取会员总数
func (m *MemberModel) GetCount() (int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "member")

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("获取会员总数失败", "error", err)
		return 0, err
	}

	return count, nil
}

// GetList 获取会员列表
func (m *MemberModel) GetList(page, pageSize int) ([]*Member, int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "member")
	qb.Select("*")

	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("获取会员总数失败", "error", err)
		return nil, 0, err
	}

	// 设置分页
	offset := (page - 1) * pageSize
	qb.OrderBy("mid DESC")
	qb.Limit(pageSize, offset)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取会员列表失败", "error", err)
		return nil, 0, err
	}

	// 转换为会员对象
	members := make([]*Member, 0, len(results))
	for _, result := range results {
		member := &Member{}

		// 处理ID字段
		if mid, ok := result["mid"].(int64); ok {
			member.ID = mid
		} else if mid, ok := result["mid"].([]byte); ok {
			if idStr := string(mid); idStr != "" {
				if idInt, err := strconv.ParseInt(idStr, 10, 64); err == nil {
					member.ID = idInt
				}
			}
		}

		member.Username, _ = result["userid"].(string)
		member.Password, _ = result["pwd"].(string)
		member.Email, _ = result["email"].(string)
		member.MType = convertToInt(result["mtype"])
		member.Sex, _ = result["sex"].(string)
		member.Face, _ = result["face"].(string)
		member.Avatar = member.Face // 兼容性
		member.Score = convertToInt(result["scores"])
		member.Money = convertToFloat64(result["money"])

		// 处理时间戳字段
		if jointime, ok := result["jointime"].(int64); ok {
			member.RegTime = time.Unix(jointime, 0)
		} else if jointime, ok := result["jointime"].([]byte); ok {
			if timeStr := string(jointime); timeStr != "" {
				if timestamp, err := strconv.ParseInt(timeStr, 10, 64); err == nil {
					member.RegTime = time.Unix(timestamp, 0)
				}
			}
		}

		if logintime, ok := result["logintime"].(int64); ok {
			member.LastLogin = time.Unix(logintime, 0)
		} else if logintime, ok := result["logintime"].([]byte); ok {
			if timeStr := string(logintime); timeStr != "" {
				if timestamp, err := strconv.ParseInt(timeStr, 10, 64); err == nil {
					member.LastLogin = time.Unix(timestamp, 0)
				}
			}
		}

		member.RegIP, _ = result["joinip"].(string)
		member.LastIP, _ = result["loginip"].(string)
		member.Status = 1 // 默认状态为正常

		members = append(members, member)
	}

	return members, total, nil
}

// Search 搜索会员
func (m *MemberModel) Search(keyword string, page, pageSize int) ([]*Member, int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "member")
	qb.Select("*")
	qb.Where("userid LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%")

	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("获取搜索会员总数失败", "keyword", keyword, "error", err)
		return nil, 0, err
	}

	// 设置分页
	offset := (page - 1) * pageSize
	qb.OrderBy("mid DESC")
	qb.Limit(pageSize, offset)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("搜索会员失败", "keyword", keyword, "error", err)
		return nil, 0, err
	}

	// 转换为会员对象
	members := make([]*Member, 0, len(results))
	for _, result := range results {
		member := &Member{}

		// 处理ID字段
		if mid, ok := result["mid"].(int64); ok {
			member.ID = mid
		} else if mid, ok := result["mid"].([]byte); ok {
			if idStr := string(mid); idStr != "" {
				if idInt, err := strconv.ParseInt(idStr, 10, 64); err == nil {
					member.ID = idInt
				}
			}
		}

		member.Username, _ = result["userid"].(string)
		member.Password, _ = result["pwd"].(string)
		member.Email, _ = result["email"].(string)
		member.MType = convertToInt(result["mtype"])
		member.Sex, _ = result["sex"].(string)
		member.Face, _ = result["face"].(string)
		member.Avatar = member.Face // 兼容性
		member.Score = convertToInt(result["scores"])
		member.Money = convertToFloat64(result["money"])

		// 处理时间戳字段
		if jointime, ok := result["jointime"].(int64); ok {
			member.RegTime = time.Unix(jointime, 0)
		} else if jointime, ok := result["jointime"].([]byte); ok {
			if timeStr := string(jointime); timeStr != "" {
				if timestamp, err := strconv.ParseInt(timeStr, 10, 64); err == nil {
					member.RegTime = time.Unix(timestamp, 0)
				}
			}
		}

		if logintime, ok := result["logintime"].(int64); ok {
			member.LastLogin = time.Unix(logintime, 0)
		} else if logintime, ok := result["logintime"].([]byte); ok {
			if timeStr := string(logintime); timeStr != "" {
				if timestamp, err := strconv.ParseInt(timeStr, 10, 64); err == nil {
					member.LastLogin = time.Unix(timestamp, 0)
				}
			}
		}

		member.RegIP, _ = result["joinip"].(string)
		member.LastIP, _ = result["loginip"].(string)
		member.Status = 1 // 默认状态为正常

		members = append(members, member)
	}

	return members, total, nil
}

// GetByUsername 根据用户名获取会员
func (m *MemberModel) GetByUsername(username string) (*Member, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "member")
	qb.Where("userid = ?", username)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取会员失败", "username", username, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("member not found: %s", username)
	}

	// 转换为会员
	member := &Member{}

	// 处理ID字段
	if mid, ok := result["mid"].(int64); ok {
		member.ID = mid
	} else if mid, ok := result["mid"].([]byte); ok {
		if idStr := string(mid); idStr != "" {
			if idInt, err := strconv.ParseInt(idStr, 10, 64); err == nil {
				member.ID = idInt
			}
		}
	}

	member.Username, _ = result["userid"].(string)
	member.Password, _ = result["pwd"].(string)
	member.Email, _ = result["email"].(string)
	member.MType = convertToInt(result["mtype"])
	member.Sex, _ = result["sex"].(string)
	member.Face, _ = result["face"].(string)
	member.Avatar = member.Face // 兼容性
	member.Score = convertToInt(result["scores"])
	member.Money = convertToFloat64(result["money"])

	// 处理时间戳字段
	if jointime, ok := result["jointime"].(int64); ok {
		member.RegTime = time.Unix(jointime, 0)
	} else if jointime, ok := result["jointime"].([]byte); ok {
		if timeStr := string(jointime); timeStr != "" {
			if timestamp, err := strconv.ParseInt(timeStr, 10, 64); err == nil {
				member.RegTime = time.Unix(timestamp, 0)
			}
		}
	}

	if logintime, ok := result["logintime"].(int64); ok {
		member.LastLogin = time.Unix(logintime, 0)
	} else if logintime, ok := result["logintime"].([]byte); ok {
		if timeStr := string(logintime); timeStr != "" {
			if timestamp, err := strconv.ParseInt(timeStr, 10, 64); err == nil {
				member.LastLogin = time.Unix(timestamp, 0)
			}
		}
	}

	member.RegIP, _ = result["joinip"].(string)
	member.LastIP, _ = result["loginip"].(string)
	member.Status = 1 // 默认状态为正常

	return member, nil
}

// GetByEmail 根据邮箱获取会员
func (m *MemberModel) GetByEmail(email string) (*Member, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "member")
	qb.Where("email = ?", email)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取会员失败", "email", email, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("member not found: %s", email)
	}

	// 转换为会员
	member := &Member{}

	// 处理ID字段
	if mid, ok := result["mid"].(int64); ok {
		member.ID = mid
	} else if mid, ok := result["mid"].([]byte); ok {
		if idStr := string(mid); idStr != "" {
			if idInt, err := strconv.ParseInt(idStr, 10, 64); err == nil {
				member.ID = idInt
			}
		}
	}

	member.Username, _ = result["userid"].(string)
	member.Password, _ = result["pwd"].(string)
	member.Email, _ = result["email"].(string)
	member.MType = convertToInt(result["mtype"])
	member.Sex, _ = result["sex"].(string)
	member.Face, _ = result["face"].(string)
	member.Avatar = member.Face // 兼容性
	member.Score = convertToInt(result["scores"])
	member.Money = convertToFloat64(result["money"])

	// 处理时间戳字段
	if jointime, ok := result["jointime"].(int64); ok {
		member.RegTime = time.Unix(jointime, 0)
	} else if jointime, ok := result["jointime"].([]byte); ok {
		if timeStr := string(jointime); timeStr != "" {
			if timestamp, err := strconv.ParseInt(timeStr, 10, 64); err == nil {
				member.RegTime = time.Unix(timestamp, 0)
			}
		}
	}

	if logintime, ok := result["logintime"].(int64); ok {
		member.LastLogin = time.Unix(logintime, 0)
	} else if logintime, ok := result["logintime"].([]byte); ok {
		if timeStr := string(logintime); timeStr != "" {
			if timestamp, err := strconv.ParseInt(timeStr, 10, 64); err == nil {
				member.LastLogin = time.Unix(timestamp, 0)
			}
		}
	}

	member.RegIP, _ = result["joinip"].(string)
	member.LastIP, _ = result["loginip"].(string)
	member.Status = 1 // 默认状态为正常

	return member, nil
}

// GetAll 获取所有会员
func (m *MemberModel) GetAll(page, pageSize int) ([]*Member, int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "member")
	qb.OrderBy("mid DESC")
	qb.Limit(pageSize)
	qb.Offset((page - 1) * pageSize)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取所有会员失败", "error", err)
		return nil, 0, err
	}

	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("获取会员总数失败", "error", err)
		return nil, 0, err
	}

	// 转换为会员列表
	members := make([]*Member, 0, len(results))
	for _, result := range results {
		member := &Member{}

		// 处理ID字段
		if mid, ok := result["mid"].(int64); ok {
			member.ID = mid
		} else if mid, ok := result["mid"].([]byte); ok {
			if idStr := string(mid); idStr != "" {
				if idInt, err := strconv.ParseInt(idStr, 10, 64); err == nil {
					member.ID = idInt
				}
			}
		}

		member.Username, _ = result["userid"].(string)
		member.Password, _ = result["pwd"].(string)
		member.Email, _ = result["email"].(string)
		member.MType = convertToInt(result["mtype"])
		member.Sex, _ = result["sex"].(string)
		member.Face, _ = result["face"].(string)
		member.Avatar = member.Face // 兼容性
		member.Score = convertToInt(result["scores"])
		member.Money = convertToFloat64(result["money"])

		// 处理时间戳字段
		if jointime, ok := result["jointime"].(int64); ok {
			member.RegTime = time.Unix(jointime, 0)
		} else if jointime, ok := result["jointime"].([]byte); ok {
			if timeStr := string(jointime); timeStr != "" {
				if timestamp, err := strconv.ParseInt(timeStr, 10, 64); err == nil {
					member.RegTime = time.Unix(timestamp, 0)
				}
			}
		}

		if logintime, ok := result["logintime"].(int64); ok {
			member.LastLogin = time.Unix(logintime, 0)
		} else if logintime, ok := result["logintime"].([]byte); ok {
			if timeStr := string(logintime); timeStr != "" {
				if timestamp, err := strconv.ParseInt(timeStr, 10, 64); err == nil {
					member.LastLogin = time.Unix(timestamp, 0)
				}
			}
		}

		member.RegIP, _ = result["joinip"].(string)
		member.LastIP, _ = result["loginip"].(string)
		member.Status = 1 // 默认状态为正常

		members = append(members, member)
	}

	return members, total, nil
}

// Create 创建会员
func (m *MemberModel) Create(member *Member) (int64, error) {
	// 设置创建时间和更新时间
	now := time.Now()
	member.RegTime = now
	member.LastLogin = now

	// 执行插入 - 暂时跳过sex字段以避免数据截断问题
	result, err := m.db.Exec(
		"INSERT INTO "+m.db.TableName("member")+" (userid, pwd, email, mtype, scores, money, jointime, joinip, logintime, loginip) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		member.Username, member.Password, member.Email, member.MType, member.Score, member.Money, time.Now().Unix(), member.RegIP, time.Now().Unix(), member.LastIP,
	)
	if err != nil {
		logger.Error("创建会员失败", "error", err)
		return 0, err
	}

	// 获取插入ID
	id, err := result.LastInsertId()
	if err != nil {
		logger.Error("获取插入ID失败", "error", err)
		return 0, err
	}

	return id, nil
}

// Update 更新会员
func (m *MemberModel) Update(member *Member) error {
	// 执行更新 - 暂时跳过sex字段以避免数据截断问题
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("member")+" SET userid = ?, pwd = ?, email = ?, mtype = ?, scores = ?, money = ?, jointime = ?, joinip = ?, logintime = ?, loginip = ? WHERE mid = ?",
		member.Username, member.Password, member.Email, member.MType, member.Score, member.Money, member.RegTime.Unix(), member.RegIP, member.LastLogin.Unix(), member.LastIP, member.ID,
	)
	if err != nil {
		logger.Error("更新会员失败", "error", err)
		return err
	}

	return nil
}

// Delete 删除会员
func (m *MemberModel) Delete(id int64) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("member")+" WHERE mid = ?", id)
	if err != nil {
		logger.Error("删除会员失败", "error", err)
		return err
	}

	return nil
}

// UpdateScore 更新会员积分
func (m *MemberModel) UpdateScore(id int64, score int) error {
	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("member")+" SET scores = scores + ? WHERE mid = ?",
		score, id,
	)
	if err != nil {
		logger.Error("更新会员积分失败", "error", err)
		return err
	}

	return nil
}

// UpdateMoney 更新会员余额
func (m *MemberModel) UpdateMoney(id int64, money float64) error {
	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("member")+" SET money = money + ? WHERE mid = ?",
		money, id,
	)
	if err != nil {
		logger.Error("更新会员余额失败", "error", err)
		return err
	}

	return nil
}

// UpdateLoginInfo 更新会员登录信息
func (m *MemberModel) UpdateLoginInfo(id int64, ip string) error {
	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("member")+" SET logintime = ?, loginip = ? WHERE mid = ?",
		time.Now().Unix(), ip, id,
	)
	if err != nil {
		logger.Error("更新会员登录信息失败", "error", err)
		return err
	}

	return nil
}

// CheckLogin 检查登录
func (m *MemberModel) CheckLogin(username, password string) (*Member, error) {
	// 根据用户名获取会员
	member, err := m.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	// 验证密码
	if !security.CheckPassword(password, member.Password) {
		return nil, fmt.Errorf("invalid password")
	}

	return member, nil
}

// GetTotalCount 获取会员总数
func (m *MemberModel) GetTotalCount() (int, error) {
	return m.GetCount()
}

// GetTodayCount 获取今日新增会员数
func (m *MemberModel) GetTodayCount() (int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "member")
	qb.Where("DATE(FROM_UNIXTIME(jointime)) = CURDATE()")

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("获取今日新增会员数失败", "error", err)
		return 0, err
	}

	return count, nil
}

// GetActiveCount 获取活跃会员数（最近N天登录）
func (m *MemberModel) GetActiveCount(days int) (int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "member")
	qb.Where("logintime >= UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL ? DAY))", days)

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("获取活跃会员数失败", "days", days, "error", err)
		return 0, err
	}

	return count, nil
}

// GetDisabledCount 获取禁用会员数
func (m *MemberModel) GetDisabledCount() (int, error) {
	// aq3cmsCMS会员表可能没有status字段，暂时返回0
	// 如果需要实现禁用功能，需要添加相应字段
	return 0, nil
}

// GetLatest 获取最新注册的会员
func (m *MemberModel) GetLatest(limit int) ([]*Member, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "member")
	qb.Select("*")
	qb.OrderBy("jointime DESC")
	qb.Limit(limit)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取最新会员失败", "limit", limit, "error", err)
		return nil, err
	}

	// 转换为会员对象
	members := make([]*Member, 0, len(results))
	for _, result := range results {
		member := &Member{}

		// 处理ID字段
		if mid, ok := result["mid"].(int64); ok {
			member.ID = mid
		} else if mid, ok := result["mid"].([]byte); ok {
			if idStr := string(mid); idStr != "" {
				if idInt, err := strconv.ParseInt(idStr, 10, 64); err == nil {
					member.ID = idInt
				}
			}
		}

		member.Username, _ = result["userid"].(string)
		member.Email, _ = result["email"].(string)
		member.MType = convertToInt(result["mtype"])
		member.Sex, _ = result["sex"].(string)
		member.Face, _ = result["face"].(string)
		member.Avatar = member.Face // 兼容性
		member.Score = convertToInt(result["scores"])
		member.Money = convertToFloat64(result["money"])

		// 处理时间戳字段
		if jointime, ok := result["jointime"].(int64); ok {
			member.RegTime = time.Unix(jointime, 0)
		} else if jointime, ok := result["jointime"].([]byte); ok {
			if timeStr := string(jointime); timeStr != "" {
				if timestamp, err := strconv.ParseInt(timeStr, 10, 64); err == nil {
					member.RegTime = time.Unix(timestamp, 0)
				}
			}
		}

		if logintime, ok := result["logintime"].(int64); ok {
			member.LastLogin = time.Unix(logintime, 0)
		} else if logintime, ok := result["logintime"].([]byte); ok {
			if timeStr := string(logintime); timeStr != "" {
				if timestamp, err := strconv.ParseInt(timeStr, 10, 64); err == nil {
					member.LastLogin = time.Unix(timestamp, 0)
				}
			}
		}

		member.RegIP, _ = result["joinip"].(string)
		member.LastIP, _ = result["loginip"].(string)
		member.Status = 1 // 默认状态为正常

		members = append(members, member)
	}

	return members, nil
}

// GetTypeStats 获取会员类型统计
func (m *MemberModel) GetTypeStats() ([]map[string]interface{}, error) {
	// 构建查询
	query := `
		SELECT
			mt.typename as type_name,
			COUNT(m.mid) as count,
			mt.id as type_id
		FROM ` + m.db.TableName("member") + ` m
		LEFT JOIN ` + m.db.TableName("member_type") + ` mt ON m.mtype = mt.id
		GROUP BY m.mtype, mt.typename, mt.id
		ORDER BY count DESC
	`

	rows, err := m.db.DB.Query(query)
	if err != nil {
		logger.Error("获取会员类型统计失败", "error", err)
		return nil, err
	}
	defer rows.Close()

	stats := make([]map[string]interface{}, 0)
	for rows.Next() {
		var typeName string
		var count int
		var typeID int

		err := rows.Scan(&typeName, &count, &typeID)
		if err != nil {
			logger.Error("扫描会员类型统计失败", "error", err)
			continue
		}

		if typeName == "" {
			typeName = "未分类"
		}

		stats = append(stats, map[string]interface{}{
			"type_name": typeName,
			"count":     count,
			"type_id":   typeID,
		})
	}

	return stats, nil
}

// convertToFloat64 将interface{}转换为float64，支持多种类型
func convertToFloat64(value interface{}) float64 {
	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case int32:
		return float64(v)
	case []byte:
		if str := string(v); str != "" {
			if f, err := strconv.ParseFloat(str, 64); err == nil {
				return f
			}
		}
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return 0.0
}
