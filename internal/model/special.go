package model

import (
	"fmt"
	"time"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// Special 专题模型
type Special struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`       // 专题标题
	Typeid      int64     `json:"typeid"`      // 栏目ID
	Note        string    `json:"note"`        // 专题描述
	Pic         string    `json:"pic"`         // 专题图片
	PubDate     time.Time `json:"pubdate"`     // 发布时间
	LastUpdate  time.Time `json:"lastupdate"`  // 最后更新时间
	IsHot       int       `json:"ishot"`       // 是否热门
	Click       int       `json:"click"`       // 点击量
	Template    string    `json:"template"`    // 专题模板
	TemplateList string   `json:"templatelist"` // 列表模板
	Filename    string    `json:"filename"`    // 文件名
	Status      int       `json:"status"`      // 状态
	Keywords    string    `json:"keywords"`    // 关键词
	Description string    `json:"description"` // 描述
	Content     string    `json:"content"`     // 专题内容
}

// SpecialModel 专题模型操作
type SpecialModel struct {
	db *database.DB
}

// NewSpecialModel 创建专题模型
func NewSpecialModel(db *database.DB) *SpecialModel {
	return &SpecialModel{
		db: db,
	}
}

// GetByID 根据ID获取专题
func (m *SpecialModel) GetByID(id int64) (*Special, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "special")
	qb.Select("*")
	qb.Where("id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("查询专题失败", "id", id, "error", err)
		return nil, err
	}

	if result == nil {
		return nil, fmt.Errorf("专题不存在")
	}

	// 转换为专题对象
	special := &Special{}
	special.ID, _ = result["id"].(int64)
	special.Title, _ = result["title"].(string)
	special.Typeid, _ = result["typeid"].(int64)
	special.Note, _ = result["note"].(string)
	special.Pic, _ = result["pic"].(string)

	// 处理日期
	if pubdate, ok := result["pubdate"].(time.Time); ok {
		special.PubDate = pubdate
	}
	if lastupdate, ok := result["lastupdate"].(time.Time); ok {
		special.LastUpdate = lastupdate
	}

	// 处理整数字段
	if ishot, ok := result["ishot"].(int64); ok {
		special.IsHot = int(ishot)
	}
	if click, ok := result["click"].(int64); ok {
		special.Click = int(click)
	}

	special.Template, _ = result["template"].(string)
	special.TemplateList, _ = result["templatelist"].(string)
	special.Filename, _ = result["filename"].(string)

	// 处理整数字段
	if status, ok := result["status"].(int64); ok {
		special.Status = int(status)
	}

	special.Keywords, _ = result["keywords"].(string)
	special.Description, _ = result["description"].(string)
	special.Content, _ = result["content"].(string)

	return special, nil
}

// GetByFilename 根据文件名获取专题
func (m *SpecialModel) GetByFilename(filename string) (*Special, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "special")
	qb.Select("*")
	qb.Where("filename = ?", filename)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("查询专题失败", "filename", filename, "error", err)
		return nil, err
	}

	if result == nil {
		return nil, fmt.Errorf("专题不存在")
	}

	// 转换为专题对象
	special := &Special{}
	special.ID, _ = result["id"].(int64)
	special.Title, _ = result["title"].(string)
	special.Typeid, _ = result["typeid"].(int64)
	special.Note, _ = result["note"].(string)
	special.Pic, _ = result["pic"].(string)

	// 处理日期
	if pubdate, ok := result["pubdate"].(time.Time); ok {
		special.PubDate = pubdate
	}
	if lastupdate, ok := result["lastupdate"].(time.Time); ok {
		special.LastUpdate = lastupdate
	}

	// 处理整数字段
	if ishot, ok := result["ishot"].(int64); ok {
		special.IsHot = int(ishot)
	}
	if click, ok := result["click"].(int64); ok {
		special.Click = int(click)
	}

	special.Template, _ = result["template"].(string)
	special.TemplateList, _ = result["templatelist"].(string)
	special.Filename, _ = result["filename"].(string)

	// 处理整数字段
	if status, ok := result["status"].(int64); ok {
		special.Status = int(status)
	}

	special.Keywords, _ = result["keywords"].(string)
	special.Description, _ = result["description"].(string)
	special.Content, _ = result["content"].(string)

	return special, nil
}

// GetAll 获取所有专题
func (m *SpecialModel) GetAll() ([]*Special, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "special")
	qb.Select("*")
	qb.Where("status = 1")

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("查询专题列表失败", "error", err)
		return nil, err
	}



	// 转换为专题对象
	specials := make([]*Special, 0, len(results))
	for _, result := range results {
		special := &Special{}
		special.ID, _ = result["id"].(int64)
		special.Title, _ = result["title"].(string)
		special.Typeid, _ = result["typeid"].(int64)
		special.Note, _ = result["note"].(string)
		special.Pic, _ = result["pic"].(string)

		// 处理日期
		if pubdate, ok := result["pubdate"].(time.Time); ok {
			special.PubDate = pubdate
		}
		if lastupdate, ok := result["lastupdate"].(time.Time); ok {
			special.LastUpdate = lastupdate
		}

		// 处理整数字段
		if ishot, ok := result["ishot"].(int64); ok {
			special.IsHot = int(ishot)
		}
		if click, ok := result["click"].(int64); ok {
			special.Click = int(click)
		}

		special.Template, _ = result["template"].(string)
		special.TemplateList, _ = result["templatelist"].(string)
		special.Filename, _ = result["filename"].(string)

		// 处理整数字段
		if status, ok := result["status"].(int64); ok {
			special.Status = int(status)
		}

		special.Keywords, _ = result["keywords"].(string)
		special.Description, _ = result["description"].(string)

		specials = append(specials, special)
	}

	return specials, nil
}

// GetHotSpecials 获取热门专题
func (m *SpecialModel) GetHotSpecials(limit int) ([]*Special, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "special")
	qb.Select("*")
	qb.Where("status = 1")
	qb.Where("ishot = 1")
	qb.OrderBy("click DESC")
	qb.Limit(limit)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("查询热门专题失败", "error", err)
		return nil, err
	}

	// 转换为专题对象
	specials := make([]*Special, 0, len(results))
	for _, result := range results {
		special := &Special{}
		special.ID, _ = result["id"].(int64)
		special.Title, _ = result["title"].(string)
		special.Typeid, _ = result["typeid"].(int64)
		special.Note, _ = result["note"].(string)
		special.Pic, _ = result["pic"].(string)

		// 处理日期
		if pubdate, ok := result["pubdate"].(time.Time); ok {
			special.PubDate = pubdate
		}
		if lastupdate, ok := result["lastupdate"].(time.Time); ok {
			special.LastUpdate = lastupdate
		}

		// 处理整数字段
		if ishot, ok := result["ishot"].(int64); ok {
			special.IsHot = int(ishot)
		}
		if click, ok := result["click"].(int64); ok {
			special.Click = int(click)
		}

		special.Template, _ = result["template"].(string)
		special.TemplateList, _ = result["templatelist"].(string)
		special.Filename, _ = result["filename"].(string)

		// 处理整数字段
		if status, ok := result["status"].(int64); ok {
			special.Status = int(status)
		}

		special.Keywords, _ = result["keywords"].(string)
		special.Description, _ = result["description"].(string)

		specials = append(specials, special)
	}

	return specials, nil
}

// GetArticles 获取专题文章
func (m *SpecialModel) GetArticles(specialID int64, page, pageSize int) ([]*Article, int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "special_content")
	qb.Select("sc.aid")
	qb.Where("sc.specialid = ?", specialID)

	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("查询专题文章总数失败", "specialid", specialID, "error", err)
		return nil, 0, err
	}

	// 设置分页
	offset := (page - 1) * pageSize
	qb.OrderBy("sc.sortrank ASC")
	qb.Limit(pageSize, offset)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("查询专题文章列表失败", "specialid", specialID, "error", err)
		return nil, 0, err
	}

	// 提取文章ID
	aids := make([]int64, 0, len(results))
	for _, result := range results {
		if aid, ok := result["aid"].(int64); ok {
			aids = append(aids, aid)
		}
	}

	// 如果没有文章，返回空列表
	if len(aids) == 0 {
		return []*Article{}, 0, nil
	}

	// 查询文章详情
	articleModel := NewArticleModel(m.db)
	articles := make([]*Article, 0, len(aids))

	for _, aid := range aids {
		article, err := articleModel.GetByID(aid)
		if err != nil {
			logger.Error("查询文章详情失败", "aid", aid, "error", err)
			continue
		}

		articles = append(articles, article)
	}

	return articles, total, nil
}

// Create 创建专题
func (m *SpecialModel) Create(special *Special) (int64, error) {
	// 检查文件名是否已存在
	if special.Filename != "" {
		existingSpecial, err := m.GetByFilename(special.Filename)
		if err == nil && existingSpecial != nil {
			return 0, fmt.Errorf("文件名已存在")
		}
	}

	// 设置默认值
	if special.PubDate.IsZero() {
		special.PubDate = time.Now()
	}
	if special.LastUpdate.IsZero() {
		special.LastUpdate = time.Now()
	}

	// 构建数据
	data := map[string]interface{}{
		"title":        special.Title,
		"typeid":       special.Typeid,
		"note":         special.Note,
		"pic":          special.Pic,
		"pubdate":      special.PubDate,
		"lastupdate":   special.LastUpdate,
		"ishot":        special.IsHot,
		"click":        special.Click,
		"template":     special.Template,
		"templatelist": special.TemplateList,
		"filename":     special.Filename,
		"status":       special.Status,
		"keywords":     special.Keywords,
		"description":  special.Description,
		"content":      special.Content,
	}

	// 执行插入
	qb := database.NewQueryBuilder(m.db, "special")
	id, err := qb.Insert(data)
	if err != nil {
		logger.Error("创建专题失败", "title", special.Title, "error", err)
		return 0, err
	}

	return id, nil
}

// Update 更新专题
func (m *SpecialModel) Update(special *Special) error {
	// 检查文件名是否已存在
	if special.Filename != "" {
		existingSpecial, err := m.GetByFilename(special.Filename)
		if err == nil && existingSpecial != nil && existingSpecial.ID != special.ID {
			return fmt.Errorf("文件名已存在")
		}
	}

	// 更新最后更新时间
	special.LastUpdate = time.Now()

	// 构建数据
	data := map[string]interface{}{
		"title":        special.Title,
		"typeid":       special.Typeid,
		"note":         special.Note,
		"pic":          special.Pic,
		"lastupdate":   special.LastUpdate,
		"ishot":        special.IsHot,
		"template":     special.Template,
		"templatelist": special.TemplateList,
		"filename":     special.Filename,
		"status":       special.Status,
		"keywords":     special.Keywords,
		"description":  special.Description,
		"content":      special.Content,
	}

	// 构建查询
	qb := database.NewQueryBuilder(m.db, "special")
	qb.Where("id = ?", special.ID)

	// 执行更新
	_, err := qb.Update(data)
	if err != nil {
		logger.Error("更新专题失败", "id", special.ID, "error", err)
		return err
	}

	return nil
}

// Delete 删除专题
func (m *SpecialModel) Delete(id int64) error {
	// 开始事务
	tx, err := m.db.Begin()
	if err != nil {
		logger.Error("开始事务失败", "error", err)
		return err
	}
	defer tx.Rollback()

	// 删除专题
	_, err = tx.Exec(
		"DELETE FROM "+m.db.TableName("special")+" WHERE id = ?",
		id,
	)
	if err != nil {
		logger.Error("删除专题失败", "id", id, "error", err)
		return err
	}

	// 删除专题文章关联
	_, err = tx.Exec(
		"DELETE FROM "+m.db.TableName("special_content")+" WHERE specialid = ?",
		id,
	)
	if err != nil {
		logger.Error("删除专题文章关联失败", "specialid", id, "error", err)
		return err
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		logger.Error("提交事务失败", "error", err)
		return err
	}

	return nil
}

// AddArticle 添加文章到专题
func (m *SpecialModel) AddArticle(specialID, articleID int64, sortRank int) error {
	// 检查文章是否已存在
	qb := database.NewQueryBuilder(m.db, "special_content")
	qb.Select("*")
	qb.Where("specialid = ?", specialID)
	qb.Where("aid = ?", articleID)

	result, err := qb.First()
	if err != nil {
		logger.Error("查询专题文章关联失败", "specialid", specialID, "aid", articleID, "error", err)
		return err
	}

	if result != nil {
		// 文章已存在，更新排序
		data := map[string]interface{}{
			"sortrank": sortRank,
		}

		qb = database.NewQueryBuilder(m.db, "special_content")
		qb.Where("specialid = ?", specialID)
		qb.Where("aid = ?", articleID)

		_, err = qb.Update(data)
		if err != nil {
			logger.Error("更新专题文章排序失败", "specialid", specialID, "aid", articleID, "error", err)
			return err
		}

		return nil
	}

	// 添加文章到专题
	data := map[string]interface{}{
		"specialid": specialID,
		"aid":       articleID,
		"sortrank":  sortRank,
	}

	qb = database.NewQueryBuilder(m.db, "special_content")
	_, err = qb.Insert(data)
	if err != nil {
		logger.Error("添加文章到专题失败", "specialid", specialID, "aid", articleID, "error", err)
		return err
	}

	return nil
}

// RemoveArticle 从专题移除文章
func (m *SpecialModel) RemoveArticle(specialID, articleID int64) error {
	// 删除专题文章关联
	qb := database.NewQueryBuilder(m.db, "special_content")
	qb.Where("specialid = ?", specialID)
	qb.Where("aid = ?", articleID)

	_, err := qb.Delete()
	if err != nil {
		logger.Error("从专题移除文章失败", "specialid", specialID, "aid", articleID, "error", err)
		return err
	}

	return nil
}

// IncrementClick 增加专题点击量
func (m *SpecialModel) IncrementClick(id int64) error {
	// 执行更新
	_, err := m.db.Execute("UPDATE "+m.db.TableName("special")+" SET click = click + 1 WHERE id = ?", id)
	if err != nil {
		logger.Error("更新专题点击量失败", "id", id, "error", err)
		return err
	}

	return nil
}
