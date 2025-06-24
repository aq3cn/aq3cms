package model

import (
	"fmt"
	"strconv"
	"time"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// Category 栏目
type Category struct {
	ID              int64       `json:"id"`
	ParentID        int64       `json:"reid"`            // 父栏目ID
	TypeName        string      `json:"typename"`        // 栏目名称
	EnName          string      `json:"enname"`          // 英文名称
	TypeDir         string      `json:"typedir"`         // 栏目目录
	IsHidden        int         `json:"ishidden"`        // 是否隐藏
	ChannelType     int         `json:"channeltype"`     // 栏目类型
	CrossID         int64       `json:"crossid"`         // 交叉栏目ID
	Description     string      `json:"description"`     // 栏目描述
	Keywords        string      `json:"keywords"`        // 关键词
	SortRank        int         `json:"sortrank"`        // 排序
	ListTpl         string      `json:"listtpl"`         // 列表模板
	ArticleTpl      string      `json:"articletpl"`      // 文章模板
	ArticleTemplate string      `json:"articletemplate"` // 文章模板（兼容旧版）
	Status          int         `json:"status"`          // 状态
	CreateTime      time.Time   `json:"createtime"`      // 创建时间
	UpdateTime      time.Time   `json:"updatetime"`      // 更新时间
	Children        []*Category `json:"-"`               // 子栏目
	ArticleCount    int         `json:"article_count"`   // 文章数量（运行时计算）
}

// CategoryModel 栏目模型
type CategoryModel struct {
	db *database.DB
}

// NewCategoryModel 创建栏目模型
func NewCategoryModel(db *database.DB) *CategoryModel {
	return &CategoryModel{
		db: db,
	}
}

// GetByID 根据ID获取栏目
func (m *CategoryModel) GetByID(id int64) (*Category, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "arctype")
	qb.Where("id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取栏目失败", "id", id, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("category not found: %d", id)
	}

	// 转换为栏目
	category := &Category{}

	// 处理ID字段 - 支持多种类型
	if id, ok := result["id"].(int64); ok {
		category.ID = id
	} else if id, ok := result["id"].(int); ok {
		category.ID = int64(id)
	} else if id, ok := result["id"].(string); ok {
		if idInt, err := strconv.ParseInt(id, 10, 64); err == nil {
			category.ID = idInt
		}
	}

	// 处理ParentID字段
	if reid, ok := result["reid"].(int64); ok {
		category.ParentID = reid
	} else if reid, ok := result["reid"].(int); ok {
		category.ParentID = int64(reid)
	} else if reid, ok := result["reid"].(string); ok {
		if reidInt, err := strconv.ParseInt(reid, 10, 64); err == nil {
			category.ParentID = reidInt
		}
	}

	category.TypeName, _ = result["typename"].(string)
	category.TypeDir, _ = result["typedir"].(string)

	// 处理整数字段
	if ishidden, ok := result["ishidden"].(int64); ok {
		category.IsHidden = int(ishidden)
	} else if ishidden, ok := result["ishidden"].(int); ok {
		category.IsHidden = ishidden
	}

	if channeltype, ok := result["channeltype"].(int64); ok {
		category.ChannelType = int(channeltype)
	} else if channeltype, ok := result["channeltype"].(int); ok {
		category.ChannelType = channeltype
	}

	if sortrank, ok := result["sortrank"].(int64); ok {
		category.SortRank = int(sortrank)
	} else if sortrank, ok := result["sortrank"].(int); ok {
		category.SortRank = sortrank
	}

	// 处理CrossID字段
	if crossid, ok := result["crossid"].(int64); ok {
		category.CrossID = crossid
	} else if crossid, ok := result["crossid"].(int); ok {
		category.CrossID = int64(crossid)
	}

	category.Description, _ = result["description"].(string)
	category.Keywords, _ = result["keywords"].(string)
	category.ListTpl, _ = result["templist"].(string)
	category.ArticleTpl, _ = result["temparticle"].(string)

	// 设置默认状态为启用
	category.Status = 1

	return category, nil
}

// GetAll 获取所有栏目
func (m *CategoryModel) GetAll() ([]*Category, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "arctype")
	qb.OrderBy("sortrank ASC")

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取所有栏目失败", "error", err)
		return nil, err
	}

	// 转换为栏目列表
	categories := make([]*Category, 0, len(results))
	for _, result := range results {
		category := &Category{}

		// 处理ID字段 - 支持多种类型
		if id, ok := result["id"].(int64); ok {
			category.ID = id
		} else if id, ok := result["id"].(int); ok {
			category.ID = int64(id)
		} else if id, ok := result["id"].(string); ok {
			if idInt, err := strconv.ParseInt(id, 10, 64); err == nil {
				category.ID = idInt
			}
		}

		// 处理ParentID字段
		if reid, ok := result["reid"].(int64); ok {
			category.ParentID = reid
		} else if reid, ok := result["reid"].(int); ok {
			category.ParentID = int64(reid)
		} else if reid, ok := result["reid"].(string); ok {
			if reidInt, err := strconv.ParseInt(reid, 10, 64); err == nil {
				category.ParentID = reidInt
			}
		}

		category.TypeName, _ = result["typename"].(string)
		category.TypeDir, _ = result["typedir"].(string)

		// 处理整数字段
		if ishidden, ok := result["ishidden"].(int64); ok {
			category.IsHidden = int(ishidden)
		} else if ishidden, ok := result["ishidden"].(int); ok {
			category.IsHidden = ishidden
		}

		if channeltype, ok := result["channeltype"].(int64); ok {
			category.ChannelType = int(channeltype)
		} else if channeltype, ok := result["channeltype"].(int); ok {
			category.ChannelType = channeltype
		}

		if sortrank, ok := result["sortrank"].(int64); ok {
			category.SortRank = int(sortrank)
		} else if sortrank, ok := result["sortrank"].(int); ok {
			category.SortRank = sortrank
		}

		// 处理CrossID字段
		if crossid, ok := result["crossid"].(int64); ok {
			category.CrossID = crossid
		} else if crossid, ok := result["crossid"].(int); ok {
			category.CrossID = int64(crossid)
		}

		category.Description, _ = result["description"].(string)
		category.Keywords, _ = result["keywords"].(string)
		category.ListTpl, _ = result["templist"].(string)
		category.ArticleTpl, _ = result["temparticle"].(string)
		categories = append(categories, category)
	}

	return categories, nil
}

// GetTopCategories 获取顶级栏目
func (m *CategoryModel) GetTopCategories() ([]*Category, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "arctype")
	qb.Where("reid = ?", 0)
	qb.OrderBy("sortrank ASC")

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取顶级栏目失败", "error", err)
		return nil, err
	}

	// 转换为栏目列表
	categories := make([]*Category, 0, len(results))
	for _, result := range results {
		category := &Category{}

		// 处理ID字段 - 支持多种类型
		if id, ok := result["id"].(int64); ok {
			category.ID = id
		} else if id, ok := result["id"].(int); ok {
			category.ID = int64(id)
		} else if id, ok := result["id"].(string); ok {
			if idInt, err := strconv.ParseInt(id, 10, 64); err == nil {
				category.ID = idInt
			}
		}

		// 处理ParentID字段
		if reid, ok := result["reid"].(int64); ok {
			category.ParentID = reid
		} else if reid, ok := result["reid"].(int); ok {
			category.ParentID = int64(reid)
		} else if reid, ok := result["reid"].(string); ok {
			if reidInt, err := strconv.ParseInt(reid, 10, 64); err == nil {
				category.ParentID = reidInt
			}
		}

		category.TypeName, _ = result["typename"].(string)
		category.TypeDir, _ = result["typedir"].(string)

		// 处理整数字段
		if ishidden, ok := result["ishidden"].(int64); ok {
			category.IsHidden = int(ishidden)
		} else if ishidden, ok := result["ishidden"].(int); ok {
			category.IsHidden = ishidden
		}

		if channeltype, ok := result["channeltype"].(int64); ok {
			category.ChannelType = int(channeltype)
		} else if channeltype, ok := result["channeltype"].(int); ok {
			category.ChannelType = channeltype
		}

		if sortrank, ok := result["sortrank"].(int64); ok {
			category.SortRank = int(sortrank)
		} else if sortrank, ok := result["sortrank"].(int); ok {
			category.SortRank = sortrank
		}

		// 处理CrossID字段
		if crossid, ok := result["crossid"].(int64); ok {
			category.CrossID = crossid
		} else if crossid, ok := result["crossid"].(int); ok {
			category.CrossID = int64(crossid)
		}

		category.Description, _ = result["description"].(string)
		category.Keywords, _ = result["keywords"].(string)
		category.ListTpl, _ = result["templist"].(string)
		category.ArticleTpl, _ = result["temparticle"].(string)
		categories = append(categories, category)
	}

	return categories, nil
}

// GetChildCategories 获取子栏目
func (m *CategoryModel) GetChildCategories(parentID int64) ([]*Category, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "arctype")
	qb.Where("reid = ?", parentID)
	qb.OrderBy("sortrank ASC")

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取子栏目失败", "parentid", parentID, "error", err)
		return nil, err
	}

	// 转换为栏目列表
	categories := make([]*Category, 0, len(results))
	for _, result := range results {
		category := &Category{}

		// 处理ID字段 - 支持多种类型
		if id, ok := result["id"].(int64); ok {
			category.ID = id
		} else if id, ok := result["id"].(int); ok {
			category.ID = int64(id)
		} else if id, ok := result["id"].(string); ok {
			if idInt, err := strconv.ParseInt(id, 10, 64); err == nil {
				category.ID = idInt
			}
		}

		// 处理ParentID字段
		if reid, ok := result["reid"].(int64); ok {
			category.ParentID = reid
		} else if reid, ok := result["reid"].(int); ok {
			category.ParentID = int64(reid)
		} else if reid, ok := result["reid"].(string); ok {
			if reidInt, err := strconv.ParseInt(reid, 10, 64); err == nil {
				category.ParentID = reidInt
			}
		}

		category.TypeName, _ = result["typename"].(string)
		category.TypeDir, _ = result["typedir"].(string)

		// 处理整数字段
		if ishidden, ok := result["ishidden"].(int64); ok {
			category.IsHidden = int(ishidden)
		} else if ishidden, ok := result["ishidden"].(int); ok {
			category.IsHidden = ishidden
		}

		if channeltype, ok := result["channeltype"].(int64); ok {
			category.ChannelType = int(channeltype)
		} else if channeltype, ok := result["channeltype"].(int); ok {
			category.ChannelType = channeltype
		}

		if sortrank, ok := result["sortrank"].(int64); ok {
			category.SortRank = int(sortrank)
		} else if sortrank, ok := result["sortrank"].(int); ok {
			category.SortRank = sortrank
		}

		// 处理CrossID字段
		if crossid, ok := result["crossid"].(int64); ok {
			category.CrossID = crossid
		} else if crossid, ok := result["crossid"].(int); ok {
			category.CrossID = int64(crossid)
		}

		category.Description, _ = result["description"].(string)
		category.Keywords, _ = result["keywords"].(string)
		category.ListTpl, _ = result["templist"].(string)
		category.ArticleTpl, _ = result["temparticle"].(string)
		categories = append(categories, category)
	}

	return categories, nil
}

// Create 创建栏目
func (m *CategoryModel) Create(category *Category) (int64, error) {
	// 设置创建时间和更新时间
	now := time.Now()
	category.CreateTime = now
	category.UpdateTime = now

	// 执行插入
	result, err := m.db.Exec(
		"INSERT INTO "+m.db.TableName("arctype")+" (reid, typename, typedir, ishidden, channeltype, crossid, description, keywords, sortrank, templist, temparticle) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		category.ParentID, category.TypeName, category.TypeDir, category.IsHidden, category.ChannelType, category.CrossID, category.Description, category.Keywords, category.SortRank, category.ListTpl, category.ArticleTpl,
	)
	if err != nil {
		logger.Error("创建栏目失败", "error", err)
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

// Update 更新栏目
func (m *CategoryModel) Update(category *Category) error {
	// 设置更新时间
	category.UpdateTime = time.Now()

	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("arctype")+" SET reid = ?, typename = ?, typedir = ?, ishidden = ?, channeltype = ?, crossid = ?, description = ?, keywords = ?, sortrank = ?, templist = ?, temparticle = ? WHERE id = ?",
		category.ParentID, category.TypeName, category.TypeDir, category.IsHidden, category.ChannelType, category.CrossID, category.Description, category.Keywords, category.SortRank, category.ListTpl, category.ArticleTpl, category.ID,
	)
	if err != nil {
		logger.Error("更新栏目失败", "error", err)
		return err
	}

	return nil
}

// Delete 删除栏目
func (m *CategoryModel) Delete(id int64) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("arctype")+" WHERE id = ?", id)
	if err != nil {
		logger.Error("删除栏目失败", "error", err)
		return err
	}

	return nil
}

// GetByDir 根据目录获取栏目
func (m *CategoryModel) GetByDir(dir string) (*Category, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "arctype")
	qb.Select("*")
	qb.Where("typedir = ?", dir)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("根据目录获取栏目失败", "dir", dir, "error", err)
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	// 转换为栏目对象
	category := &Category{}

	// 处理ID字段 - 支持多种类型
	if id, ok := result["id"].(int64); ok {
		category.ID = id
	} else if id, ok := result["id"].(int); ok {
		category.ID = int64(id)
	} else if id, ok := result["id"].(string); ok {
		if idInt, err := strconv.ParseInt(id, 10, 64); err == nil {
			category.ID = idInt
		}
	}

	// 处理ParentID字段
	if reid, ok := result["reid"].(int64); ok {
		category.ParentID = reid
	} else if reid, ok := result["reid"].(int); ok {
		category.ParentID = int64(reid)
	} else if reid, ok := result["reid"].(string); ok {
		if reidInt, err := strconv.ParseInt(reid, 10, 64); err == nil {
			category.ParentID = reidInt
		}
	}

	category.TypeName, _ = result["typename"].(string)
	category.TypeDir, _ = result["typedir"].(string)

	// 处理整数字段
	if ishidden, ok := result["ishidden"].(int64); ok {
		category.IsHidden = int(ishidden)
	} else if ishidden, ok := result["ishidden"].(int); ok {
		category.IsHidden = ishidden
	}

	if channeltype, ok := result["channeltype"].(int64); ok {
		category.ChannelType = int(channeltype)
	} else if channeltype, ok := result["channeltype"].(int); ok {
		category.ChannelType = channeltype
	}

	// 处理CrossID字段
	if crossid, ok := result["crossid"].(int64); ok {
		category.CrossID = crossid
	} else if crossid, ok := result["crossid"].(int); ok {
		category.CrossID = int64(crossid)
	}

	category.Description, _ = result["description"].(string)
	category.Keywords, _ = result["keywords"].(string)

	// 处理整数字段
	if sortrank, ok := result["sortrank"].(int64); ok {
		category.SortRank = int(sortrank)
	} else if sortrank, ok := result["sortrank"].(int); ok {
		category.SortRank = sortrank
	}

	category.ListTpl, _ = result["templist"].(string)
	category.ArticleTpl, _ = result["temparticle"].(string)

	// 处理整数字段
	if status, ok := result["status"].(int64); ok {
		category.Status = int(status)
	} else if status, ok := result["status"].(int); ok {
		category.Status = status
	} else {
		// 如果数据库中没有status字段，设置默认值为启用
		category.Status = 1
	}

	// 处理日期
	if createtime, ok := result["createtime"].(time.Time); ok {
		category.CreateTime = createtime
	}
	if updatetime, ok := result["updatetime"].(time.Time); ok {
		category.UpdateTime = updatetime
	}

	return category, nil
}

// GetByPath 根据路径获取栏目（支持多级路径）
func (m *CategoryModel) GetByPath(parentDir, subDir string) (*Category, error) {
	// 首先获取父栏目
	parentCategory, err := m.GetByDir(parentDir)
	if err != nil {
		logger.Error("获取父栏目失败", "dir", parentDir, "error", err)
		return nil, err
	}

	if parentCategory == nil {
		return nil, nil
	}

	// 然后在父栏目下查找子栏目
	qb := database.NewQueryBuilder(m.db, "arctype")
	qb.Select("*")
	qb.Where("typedir = ? AND reid = ?", subDir, parentCategory.ID)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("根据路径获取子栏目失败", "parentDir", parentDir, "subDir", subDir, "error", err)
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	// 转换为栏目对象
	category := &Category{}

	// 处理ID字段 - 支持多种类型
	if id, ok := result["id"].(int64); ok {
		category.ID = id
	} else if id, ok := result["id"].(int); ok {
		category.ID = int64(id)
	} else if id, ok := result["id"].(string); ok {
		if idInt, err := strconv.ParseInt(id, 10, 64); err == nil {
			category.ID = idInt
		}
	}

	// 处理ParentID字段
	if reid, ok := result["reid"].(int64); ok {
		category.ParentID = reid
	} else if reid, ok := result["reid"].(int); ok {
		category.ParentID = int64(reid)
	} else if reid, ok := result["reid"].(string); ok {
		if reidInt, err := strconv.ParseInt(reid, 10, 64); err == nil {
			category.ParentID = reidInt
		}
	}

	category.TypeName, _ = result["typename"].(string)
	category.TypeDir, _ = result["typedir"].(string)

	// 处理整数字段
	if ishidden, ok := result["ishidden"].(int64); ok {
		category.IsHidden = int(ishidden)
	} else if ishidden, ok := result["ishidden"].(int); ok {
		category.IsHidden = ishidden
	}

	if channeltype, ok := result["channeltype"].(int64); ok {
		category.ChannelType = int(channeltype)
	} else if channeltype, ok := result["channeltype"].(int); ok {
		category.ChannelType = channeltype
	}

	if crossid, ok := result["crossid"].(int64); ok {
		category.CrossID = crossid
	} else if crossid, ok := result["crossid"].(int); ok {
		category.CrossID = int64(crossid)
	}

	if sortrank, ok := result["sortrank"].(int64); ok {
		category.SortRank = int(sortrank)
	} else if sortrank, ok := result["sortrank"].(int); ok {
		category.SortRank = sortrank
	}

	category.Description, _ = result["description"].(string)
	category.Keywords, _ = result["keywords"].(string)
	category.ListTpl, _ = result["templist"].(string)
	category.ArticleTpl, _ = result["temparticle"].(string)

	// 处理整数字段
	if status, ok := result["status"].(int64); ok {
		category.Status = int(status)
	} else if status, ok := result["status"].(int); ok {
		category.Status = status
	} else {
		// 如果数据库中没有status字段，设置默认值为启用
		category.Status = 1
	}

	// 处理日期
	if createtime, ok := result["createtime"].(time.Time); ok {
		category.CreateTime = createtime
	}
	if updatetime, ok := result["updatetime"].(time.Time); ok {
		category.UpdateTime = updatetime
	}

	return category, nil
}

// GetSubCategories 获取子栏目
func (m *CategoryModel) GetSubCategories(parentID int64) ([]*Category, error) {
	return m.GetChildCategories(parentID)
}

// GetCount 获取栏目总数
func (m *CategoryModel) GetCount() (int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "arctype")

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("获取栏目总数失败", "error", err)
		return 0, err
	}

	return count, nil
}

// GetArticleCount 获取栏目文章数量
func (m *CategoryModel) GetArticleCount(categoryID int64) (int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "archives")
	qb.From(m.db.TableName("archives"))
	qb.Where("typeid = ?", categoryID)

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("获取栏目文章数量失败", "categoryid", categoryID, "error", err)
		return 0, err
	}

	return count, nil
}
