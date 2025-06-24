package model

import (
	"strconv"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// MemberType 会员类型模型
type MemberType struct {
	ID          int64   `json:"id"`
	TypeName    string  `json:"typename"`
	Description string  `json:"description"`
	Rank        int     `json:"rank"`
	Money       float64 `json:"money"`
	Scores      int     `json:"scores"`
	Purviews    string  `json:"purviews"`
}

// MemberTypeModel 会员类型模型操作
type MemberTypeModel struct {
	db *database.DB
}

// NewMemberTypeModel 创建会员类型模型
func NewMemberTypeModel(db *database.DB) *MemberTypeModel {
	return &MemberTypeModel{
		db: db,
	}
}

// GetByID 根据ID获取会员类型
func (m *MemberTypeModel) GetByID(id int64) (*MemberType, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "member_type")
	qb.Select("*")
	qb.Where("id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("查询会员类型失败", "id", id, "error", err)
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	// 转换为会员类型对象
	memberType := &MemberType{}

	// 处理ID字段
	if mid, ok := result["id"].(int64); ok {
		memberType.ID = mid
	} else if mid, ok := result["id"].([]byte); ok {
		if idStr := string(mid); idStr != "" {
			if idInt, err := strconv.ParseInt(idStr, 10, 64); err == nil {
				memberType.ID = idInt
			}
		}
	}

	memberType.TypeName, _ = result["typename"].(string)
	memberType.Description, _ = result["description"].(string)
	memberType.Rank = convertToInt(result["rank"])
	memberType.Money = convertToFloat64(result["money"])
	memberType.Scores = convertToInt(result["scores"])
	memberType.Purviews, _ = result["purviews"].(string)

	return memberType, nil
}

// GetAll 获取所有会员类型
func (m *MemberTypeModel) GetAll() ([]*MemberType, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "member_type")
	qb.Select("*")
	qb.OrderBy("`rank` ASC")

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("查询所有会员类型失败", "error", err)
		return nil, err
	}

	// 转换为会员类型对象
	memberTypes := make([]*MemberType, 0, len(results))
	for _, result := range results {
		memberType := &MemberType{}

		// 处理ID字段
		if mid, ok := result["id"].(int64); ok {
			memberType.ID = mid
		} else if mid, ok := result["id"].([]byte); ok {
			if idStr := string(mid); idStr != "" {
				if idInt, err := strconv.ParseInt(idStr, 10, 64); err == nil {
					memberType.ID = idInt
				}
			}
		}

		memberType.TypeName, _ = result["typename"].(string)
		memberType.Description, _ = result["description"].(string)
		memberType.Rank = convertToInt(result["rank"])
		memberType.Money = convertToFloat64(result["money"])
		memberType.Scores = convertToInt(result["scores"])
		memberType.Purviews, _ = result["purviews"].(string)

		memberTypes = append(memberTypes, memberType)
	}

	return memberTypes, nil
}
