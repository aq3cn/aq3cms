package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// Tag 标签模型
type Tag struct {
	ID        int64     `json:"id"`
	Tag       string    `json:"tag"`       // 标签名
	Count     int       `json:"count"`     // 使用次数
	Rank      int       `json:"rank"`      // 排序
	IsHot     int       `json:"ishot"`     // 是否热门
	AddTime   time.Time `json:"addtime"`   // 添加时间
	LastUse   time.Time `json:"lastuse"`   // 最后使用时间
	TagPinyin string    `json:"tagpinyin"` // 标签拼音
}

// TagModel 标签模型操作
type TagModel struct {
	db *database.DB
}

// NewTagModel 创建标签模型
func NewTagModel(db *database.DB) *TagModel {
	return &TagModel{
		db: db,
	}
}

// GetByID 根据ID获取标签
func (m *TagModel) GetByID(id int64) (*Tag, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "tagindex")
	qb.Select("*")
	qb.Where("id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("查询标签失败", "id", id, "error", err)
		return nil, err
	}

	if result == nil {
		return nil, fmt.Errorf("标签不存在")
	}

	// 转换为标签对象
	tag := &Tag{}

	// 处理ID字段 - 支持多种类型
	if id, ok := result["id"].(int64); ok {
		tag.ID = id
	} else if id, ok := result["id"].([]byte); ok {
		if idStr := string(id); idStr != "" {
			if idInt, err := strconv.ParseInt(idStr, 10, 64); err == nil {
				tag.ID = idInt
			}
		}
	} else if id, ok := result["id"].(string); ok {
		if idInt, err := strconv.ParseInt(id, 10, 64); err == nil {
			tag.ID = idInt
		}
	}

	tag.Tag, _ = result["tag"].(string)

	// 处理整数字段 - 支持多种类型
	tag.Count = convertToInt(result["count"])
	tag.Rank = convertToInt(result["rank"])
	tag.IsHot = convertToInt(result["ishot"])

	// 处理日期
	if addTime, ok := result["addtime"].(time.Time); ok {
		tag.AddTime = addTime
	}
	if lastUse, ok := result["lastuse"].(time.Time); ok {
		tag.LastUse = lastUse
	}

	tag.TagPinyin, _ = result["tagpinyin"].(string)

	return tag, nil
}

// GetByName 根据名称获取标签
func (m *TagModel) GetByName(name string) (*Tag, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "tagindex")
	qb.Select("*")
	qb.Where("tag = ?", name)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("查询标签失败", "name", name, "error", err)
		return nil, err
	}

	if result == nil {
		return nil, fmt.Errorf("标签不存在")
	}

	// 转换为标签对象
	tag := &Tag{}

	// 处理ID字段 - 支持多种类型
	if id, ok := result["id"].(int64); ok {
		tag.ID = id
	} else if id, ok := result["id"].([]byte); ok {
		if idStr := string(id); idStr != "" {
			if idInt, err := strconv.ParseInt(idStr, 10, 64); err == nil {
				tag.ID = idInt
			}
		}
	} else if id, ok := result["id"].(string); ok {
		if idInt, err := strconv.ParseInt(id, 10, 64); err == nil {
			tag.ID = idInt
		}
	}

	tag.Tag, _ = result["tag"].(string)

	// 处理整数字段 - 支持多种类型
	tag.Count = convertToInt(result["count"])
	tag.Rank = convertToInt(result["rank"])
	tag.IsHot = convertToInt(result["ishot"])

	// 处理日期
	if addTime, ok := result["addtime"].(time.Time); ok {
		tag.AddTime = addTime
	}
	if lastUse, ok := result["lastuse"].(time.Time); ok {
		tag.LastUse = lastUse
	}

	tag.TagPinyin, _ = result["tagpinyin"].(string)

	return tag, nil
}

// GetList 获取标签列表
func (m *TagModel) GetList(page, pageSize int) ([]*Tag, int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "tagindex")
	qb.Select("*")

	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("查询标签总数失败", "error", err)
		return nil, 0, err
	}

	// 设置分页
	offset := (page - 1) * pageSize
	qb.OrderBy("count DESC")
	qb.Limit(pageSize, offset)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("查询标签列表失败", "error", err)
		return nil, 0, err
	}

	// 转换为标签对象
	tags := make([]*Tag, 0, len(results))
	for _, result := range results {
		tag := &Tag{}

		// 处理ID字段 - 支持多种类型
		if id, ok := result["id"].(int64); ok {
			tag.ID = id
		} else if id, ok := result["id"].([]byte); ok {
			if idStr := string(id); idStr != "" {
				if idInt, err := strconv.ParseInt(idStr, 10, 64); err == nil {
					tag.ID = idInt
				}
			}
		} else if id, ok := result["id"].(string); ok {
			if idInt, err := strconv.ParseInt(id, 10, 64); err == nil {
				tag.ID = idInt
			}
		}

		tag.Tag, _ = result["tag"].(string)

		// 处理整数字段 - 支持多种类型
		tag.Count = convertToInt(result["count"])
		tag.Rank = convertToInt(result["rank"])
		tag.IsHot = convertToInt(result["ishot"])

		// 处理日期
		if addTime, ok := result["addtime"].(time.Time); ok {
			tag.AddTime = addTime
		}
		if lastUse, ok := result["lastuse"].(time.Time); ok {
			tag.LastUse = lastUse
		}

		tag.TagPinyin, _ = result["tagpinyin"].(string)

		tags = append(tags, tag)
	}

	return tags, total, nil
}

// GetHotTags 获取热门标签
func (m *TagModel) GetHotTags(limit int) ([]*Tag, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "tagindex")
	qb.Select("*")
	qb.Where("ishot = 1")
	qb.OrderBy("count DESC")
	qb.Limit(limit)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("查询热门标签失败", "error", err)
		return nil, err
	}

	// 转换为标签对象
	tags := make([]*Tag, 0, len(results))
	for _, result := range results {
		tag := &Tag{}

		// 处理ID字段 - 支持多种类型
		if id, ok := result["id"].(int64); ok {
			tag.ID = id
		} else if id, ok := result["id"].([]byte); ok {
			if idStr := string(id); idStr != "" {
				if idInt, err := strconv.ParseInt(idStr, 10, 64); err == nil {
					tag.ID = idInt
				}
			}
		} else if id, ok := result["id"].(string); ok {
			if idInt, err := strconv.ParseInt(id, 10, 64); err == nil {
				tag.ID = idInt
			}
		}

		tag.Tag, _ = result["tag"].(string)

		// 处理整数字段 - 支持多种类型
		tag.Count = convertToInt(result["count"])
		tag.Rank = convertToInt(result["rank"])
		tag.IsHot = convertToInt(result["ishot"])

		// 处理日期
		if addTime, ok := result["addtime"].(time.Time); ok {
			tag.AddTime = addTime
		}
		if lastUse, ok := result["lastuse"].(time.Time); ok {
			tag.LastUse = lastUse
		}

		tag.TagPinyin, _ = result["tagpinyin"].(string)

		tags = append(tags, tag)
	}

	return tags, nil
}

// GetArticlesByTag 获取标签相关文章
func (m *TagModel) GetArticlesByTag(tagName string, page, pageSize int) ([]*Article, int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "taglist")
	qb.Select("tl.aid")
	qb.Where("tl.tag = ?", tagName)

	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("查询标签文章总数失败", "tag", tagName, "error", err)
		return nil, 0, err
	}

	// 设置分页
	offset := (page - 1) * pageSize
	qb.OrderBy("tl.aid DESC")
	qb.Limit(pageSize, offset)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("查询标签文章列表失败", "tag", tagName, "error", err)
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

// Create 创建标签
func (m *TagModel) Create(tag *Tag) (int64, error) {
	// 检查标签是否已存在
	existingTag, err := m.GetByName(tag.Tag)
	if err == nil && existingTag != nil {
		// 标签已存在，更新使用次数
		existingTag.Count++
		existingTag.LastUse = time.Now()

		err = m.Update(existingTag)
		if err != nil {
			logger.Error("更新标签失败", "tag", tag.Tag, "error", err)
			return 0, err
		}

		return existingTag.ID, nil
	}

	// 标签不存在，创建新标签
	if tag.AddTime.IsZero() {
		tag.AddTime = time.Now()
	}
	if tag.LastUse.IsZero() {
		tag.LastUse = time.Now()
	}

	// 生成拼音
	if tag.TagPinyin == "" {
		tag.TagPinyin = generatePinyin(tag.Tag)
	}

	// 构建数据
	data := map[string]interface{}{
		"tag":       tag.Tag,
		"count":     tag.Count,
		"rank":      tag.Rank,
		"ishot":     tag.IsHot,
		"addtime":   tag.AddTime.Format("2006-01-02 15:04:05"),
		"lastuse":   tag.LastUse.Format("2006-01-02 15:04:05"),
		"tagpinyin": tag.TagPinyin,
	}

	// 执行插入
	qb := database.NewQueryBuilder(m.db, "tagindex")
	id, err := qb.Insert(data)
	if err != nil {
		logger.Error("创建标签失败", "tag", tag.Tag, "error", err)
		return 0, err
	}

	return id, nil
}

// Update 更新标签
func (m *TagModel) Update(tag *Tag) error {
	// 构建数据
	data := map[string]interface{}{
		"tag":       tag.Tag,
		"count":     tag.Count,
		"rank":      tag.Rank,
		"ishot":     tag.IsHot,
		"lastuse":   tag.LastUse.Format("2006-01-02 15:04:05"),
		"tagpinyin": tag.TagPinyin,
	}

	// 构建查询
	qb := database.NewQueryBuilder(m.db, "tagindex")
	qb.Where("id = ?", tag.ID)

	// 执行更新
	_, err := qb.Update(data)
	if err != nil {
		logger.Error("更新标签失败", "id", tag.ID, "error", err)
		return err
	}

	return nil
}

// Delete 删除标签
func (m *TagModel) Delete(id int64) error {
	// 开始事务
	tx, err := m.db.Begin()
	if err != nil {
		logger.Error("开始事务失败", "error", err)
		return err
	}
	defer tx.Rollback()

	// 获取标签信息
	tag, err := m.GetByID(id)
	if err != nil {
		logger.Error("获取标签信息失败", "id", id, "error", err)
		return err
	}

	// 删除标签索引
	_, err = tx.Exec(
		"DELETE FROM "+m.db.TableName("tagindex")+" WHERE id = ?",
		id,
	)
	if err != nil {
		logger.Error("删除标签索引失败", "id", id, "error", err)
		return err
	}

	// 删除标签关联
	_, err = tx.Exec(
		"DELETE FROM "+m.db.TableName("taglist")+" WHERE tag = ?",
		tag.Tag,
	)
	if err != nil {
		logger.Error("删除标签关联失败", "tag", tag.Tag, "error", err)
		return err
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		logger.Error("提交事务失败", "error", err)
		return err
	}

	return nil
}

// AddArticleTag 添加文章标签
func (m *TagModel) AddArticleTag(aid int64, tagName string) error {
	// 检查标签是否已存在
	tag, err := m.GetByName(tagName)
	if err != nil {
		// 标签不存在，创建新标签
		tag = &Tag{
			Tag:     tagName,
			Count:   1,
			Rank:    0,
			IsHot:   0,
			AddTime: time.Now(),
			LastUse: time.Now(),
		}

		_, err = m.Create(tag)
		if err != nil {
			logger.Error("创建标签失败", "tag", tagName, "error", err)
			return err
		}
	} else {
		// 标签已存在，更新使用次数
		tag.Count++
		tag.LastUse = time.Now()

		err = m.Update(tag)
		if err != nil {
			logger.Error("更新标签失败", "tag", tagName, "error", err)
			return err
		}
	}

	// 检查文章标签关联是否已存在
	qb := database.NewQueryBuilder(m.db, "taglist")
	qb.Select("*")
	qb.Where("aid = ?", aid)
	qb.Where("tag = ?", tagName)

	result, err := qb.First()
	if err != nil {
		logger.Error("查询文章标签关联失败", "aid", aid, "tag", tagName, "error", err)
		return err
	}

	if result != nil {
		// 关联已存在，不需要重复添加
		return nil
	}

	// 添加文章标签关联
	data := map[string]interface{}{
		"aid": aid,
		"tag": tagName,
	}

	qb = database.NewQueryBuilder(m.db, "taglist")
	_, err = qb.Insert(data)
	if err != nil {
		logger.Error("添加文章标签关联失败", "aid", aid, "tag", tagName, "error", err)
		return err
	}

	return nil
}

// GetByAID 根据文章ID获取标签列表
func (m *TagModel) GetByAID(aid int64) ([]*Tag, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "taglist")
	qb.Select("tag")
	qb.Where("aid = ?", aid)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("查询文章标签失败", "aid", aid, "error", err)
		return nil, err
	}

	// 提取标签名
	tagNames := make([]string, 0, len(results))
	for _, result := range results {
		if tagName, ok := result["tag"].(string); ok {
			tagNames = append(tagNames, tagName)
		}
	}

	// 如果没有标签，返回空列表
	if len(tagNames) == 0 {
		return []*Tag{}, nil
	}

	// 查询标签详情
	tags := make([]*Tag, 0, len(tagNames))
	for _, tagName := range tagNames {
		tag, err := m.GetByName(tagName)
		if err != nil {
			logger.Error("查询标签详情失败", "tag", tagName, "error", err)
			continue
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

// RemoveArticleTag 移除文章标签
func (m *TagModel) RemoveArticleTag(aid int64, tagName string) error {
	// 删除文章标签关联
	qb := database.NewQueryBuilder(m.db, "taglist")
	qb.Where("aid = ?", aid)
	qb.Where("tag = ?", tagName)

	_, err := qb.Delete()
	if err != nil {
		logger.Error("删除文章标签关联失败", "aid", aid, "tag", tagName, "error", err)
		return err
	}

	// 更新标签使用次数
	tag, err := m.GetByName(tagName)
	if err != nil {
		logger.Error("获取标签信息失败", "tag", tagName, "error", err)
		return err
	}

	tag.Count--
	if tag.Count < 0 {
		tag.Count = 0
	}

	err = m.Update(tag)
	if err != nil {
		logger.Error("更新标签失败", "tag", tagName, "error", err)
		return err
	}

	return nil
}

// UpdateArticleTags 更新文章标签
func (m *TagModel) UpdateArticleTags(aid int64, tags string) error {
	// 开始事务
	tx, err := m.db.Begin()
	if err != nil {
		logger.Error("开始事务失败", "error", err)
		return err
	}
	defer tx.Rollback()

	// 删除文章所有标签关联
	_, err = tx.Exec(
		"DELETE FROM "+m.db.TableName("taglist")+" WHERE aid = ?",
		aid,
	)
	if err != nil {
		logger.Error("删除文章标签关联失败", "aid", aid, "error", err)
		return err
	}

	// 如果标签为空，直接提交事务
	if tags == "" {
		if err := tx.Commit(); err != nil {
			logger.Error("提交事务失败", "error", err)
			return err
		}
		return nil
	}

	// 分割标签
	tagList := strings.Split(tags, ",")

	// 添加新标签
	for _, tagName := range tagList {
		tagName = strings.TrimSpace(tagName)
		if tagName == "" {
			continue
		}

		// 检查标签是否已存在
		var tagID int64
		err = tx.QueryRow(
			"SELECT id FROM "+m.db.TableName("tagindex")+" WHERE tag = ?",
			tagName,
		).Scan(&tagID)

		if err != nil {
			// 标签不存在，创建新标签
			now := time.Now().Format("2006-01-02 15:04:05")
			_, err = tx.Exec(
				"INSERT INTO "+m.db.TableName("tagindex")+" (tag, count, `rank`, ishot, addtime, lastuse) VALUES (?, ?, ?, ?, ?, ?)",
				tagName, 1, 0, 0, now, now,
			)
			if err != nil {
				logger.Error("创建标签失败", "tag", tagName, "error", err)
				return err
			}
		} else {
			// 标签已存在，更新使用次数
			_, err = tx.Exec(
				"UPDATE "+m.db.TableName("tagindex")+" SET count = count + 1, lastuse = ? WHERE id = ?",
				time.Now().Format("2006-01-02 15:04:05"), tagID,
			)
			if err != nil {
				logger.Error("更新标签失败", "tag", tagName, "error", err)
				return err
			}
		}

		// 添加文章标签关联
		_, err = tx.Exec(
			"INSERT INTO "+m.db.TableName("taglist")+" (aid, tag) VALUES (?, ?)",
			aid, tagName,
		)
		if err != nil {
			logger.Error("添加文章标签关联失败", "aid", aid, "tag", tagName, "error", err)
			return err
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		logger.Error("提交事务失败", "error", err)
		return err
	}

	return nil
}

// DeleteByAID 删除文章所有标签
func (m *TagModel) DeleteByAID(aid int64) error {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "taglist")
	qb.Where("aid = ?", aid)

	// 执行删除
	_, err := qb.Delete()
	if err != nil {
		logger.Error("删除文章标签失败", "aid", aid, "error", err)
		return err
	}

	return nil
}

// UpdateTags 更新文章标签（别名）
func (m *TagModel) UpdateTags(aid int64, tags string) error {
	return m.UpdateArticleTags(aid, tags)
}

// GetArticleTags 获取文章标签
func (m *TagModel) GetArticleTags(aid int64) ([]string, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "taglist")
	qb.Select("tag")
	qb.Where("aid = ?", aid)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("查询文章标签失败", "aid", aid, "error", err)
		return nil, err
	}

	// 提取标签
	tags := make([]string, 0, len(results))
	for _, result := range results {
		if tag, ok := result["tag"].(string); ok {
			tags = append(tags, tag)
		}
	}

	return tags, nil
}

// GetNewTags 获取最新标签
func (m *TagModel) GetNewTags(limit int) ([]*Tag, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "tagindex")
	qb.Select("*")
	qb.OrderBy("addtime DESC")
	qb.Limit(limit, 0)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取最新标签失败", "error", err)
		return nil, err
	}

	// 转换为标签对象
	tags := make([]*Tag, 0, len(results))
	for _, result := range results {
		tag := &Tag{}
		tag.ID, _ = result["id"].(int64)
		tag.Tag, _ = result["tag"].(string)

		// 处理整数字段
		if count, ok := result["count"].(int64); ok {
			tag.Count = int(count)
		}
		if rank, ok := result["rank"].(int64); ok {
			tag.Rank = int(rank)
		}
		if ishot, ok := result["ishot"].(int64); ok {
			tag.IsHot = int(ishot)
		}

		// 处理日期
		if addtime, ok := result["addtime"].(time.Time); ok {
			tag.AddTime = addtime
		}
		if lastuse, ok := result["lastuse"].(time.Time); ok {
			tag.LastUse = lastuse
		}

		tag.TagPinyin, _ = result["tagpinyin"].(string)

		tags = append(tags, tag)
	}

	return tags, nil
}

// GetArticles 获取标签文章
func (m *TagModel) GetArticles(tagName string, page, pageSize int) ([]*Article, int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "taglist")
	qb.Select("tl.aid")
	qb.Where("tl.tag = ?", tagName)

	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("获取标签文章总数失败", "tag", tagName, "error", err)
		return nil, 0, err
	}

	// 设置分页
	offset := (page - 1) * pageSize
	qb.OrderBy("tl.addtime DESC")
	qb.Limit(pageSize, offset)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取标签文章列表失败", "tag", tagName, "error", err)
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

// 生成拼音
func generatePinyin(text string) string {
	// 简单实现，实际应使用拼音库
	return strings.ToLower(text)
}

// convertToInt 将interface{}转换为int，支持多种类型
func convertToInt(value interface{}) int {
	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case int32:
		return int(v)
	case []byte:
		if str := string(v); str != "" {
			if i, err := strconv.Atoi(str); err == nil {
				return i
			}
		}
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return 0
}
