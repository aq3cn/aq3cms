package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// Article 文章模型
type Article struct {
	ID           int64     `json:"id"`
	TypeID       int64     `json:"typeid"`
	Title        string    `json:"title"`
	ShortTitle   string    `json:"shorttitle"`
	Color        string    `json:"color"`
	Writer       string    `json:"writer"`
	Source       string    `json:"source"`
	Author       string    `json:"author"` // 作者
	LitPic       string    `json:"litpic"`
	PubDate      time.Time `json:"pubdate"`
	SendDate     time.Time `json:"senddate"`
	UpdateDate   time.Time `json:"updatedate"` // 更新时间
	Keywords     string    `json:"keywords"`
	Description  string    `json:"description"`
	Filename     string    `json:"filename"`
	IsTop        int       `json:"istop"`
	IsRecommend  int       `json:"isrecommend"`
	IsHot        int       `json:"ishot"`
	ArcRank      int       `json:"arcrank"`
	Click        int       `json:"click"`
	Status       int       `json:"status"`   // 状态
	MemberID     int64     `json:"memberid"` // 会员ID
	Body         string    `json:"body"`
	Content      string    `json:"content"` // 内容（兼容旧版）
	TypeName     string    `json:"typename"`
	TypeDir      string    `json:"typedir"`
	CategoryName string    `json:"categoryname"` // 栏目名称（兼容模板）
	TemplateFile string    `json:"templatefile"` // 自定义模板文件
	Tags         string    `json:"tags"`         // 标签
}

// ArticleModel 文章模型操作
type ArticleModel struct {
	db *database.DB
}

// NewArticleModel 创建文章模型
func NewArticleModel(db *database.DB) *ArticleModel {
	return &ArticleModel{
		db: db,
	}
}

// parseFlags 解析flag字段
func parseFlags(flagStr string) (isTop, isRecommend, isHot int) {
	if flagStr == "" {
		return 0, 0, 0
	}

	flags := strings.Split(flagStr, ",")
	for _, flag := range flags {
		switch flag {
		case "c":
			isTop = 1
		case "h":
			isRecommend = 1
		case "p":
			isHot = 1
		}
	}

	return isTop, isRecommend, isHot
}

// GetByID 根据ID获取文章
func (m *ArticleModel) GetByID(id int64) (*Article, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "archives")
	qb.Select("a.id", "a.typeid", "a.title", "a.shorttitle", "a.color", "a.writer", "a.source", "a.litpic", "a.pubdate", "a.senddate", "a.keywords", "a.description", "a.filename", "a.flag", "a.arcrank", "a.click", "t.typename", "t.typedir", "ad.body")
	qb.From(m.db.TableName("archives") + " AS a")
	qb.LeftJoin(m.db.TableName("arctype")+" AS t", "a.typeid = t.id")
	qb.LeftJoin(m.db.TableName("addonarticle")+" AS ad", "a.id = ad.aid")
	qb.Where("a.id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("查询文章失败", "id", id, "error", err)
		return nil, err
	}

	if result == nil {
		return nil, fmt.Errorf("文章不存在")
	}

	// 转换为文章对象
	article := &Article{}

	// 处理ID字段 - 使用通用转换方法
	if idVal := result["id"]; idVal != nil {
		switch v := idVal.(type) {
		case int64:
			article.ID = v
		case int:
			article.ID = int64(v)
		case int32:
			article.ID = int64(v)
		case []uint8: // MySQL可能返回字节数组
			if len(v) > 0 {
				if id, err := strconv.ParseInt(string(v), 10, 64); err == nil {
					article.ID = id
				}
			}
		case string:
			if id, err := strconv.ParseInt(v, 10, 64); err == nil {
				article.ID = id
			}
		}
	}

	// 处理TypeID字段 - 使用通用转换方法
	if typeidVal := result["typeid"]; typeidVal != nil {
		switch v := typeidVal.(type) {
		case int64:
			article.TypeID = v
		case int:
			article.TypeID = int64(v)
		case int32:
			article.TypeID = int64(v)
		case []uint8:
			if len(v) > 0 {
				if typeid, err := strconv.ParseInt(string(v), 10, 64); err == nil {
					article.TypeID = typeid
				}
			}
		case string:
			if typeid, err := strconv.ParseInt(v, 10, 64); err == nil {
				article.TypeID = typeid
			}
		}
	}
	article.Title, _ = result["title"].(string)
	article.ShortTitle, _ = result["shorttitle"].(string)
	article.Color, _ = result["color"].(string)
	article.Writer, _ = result["writer"].(string)
	article.Source, _ = result["source"].(string)
	article.LitPic, _ = result["litpic"].(string)

	// 处理日期 - 数据库中存储的是Unix时间戳
	if pubdate, ok := result["pubdate"].(int64); ok {
		article.PubDate = time.Unix(pubdate, 0)
	} else if pubdate, ok := result["pubdate"].(int); ok {
		article.PubDate = time.Unix(int64(pubdate), 0)
	} else if pubdate, ok := result["pubdate"].(time.Time); ok {
		article.PubDate = pubdate
	}

	if senddate, ok := result["senddate"].(int64); ok {
		article.SendDate = time.Unix(senddate, 0)
	} else if senddate, ok := result["senddate"].(int); ok {
		article.SendDate = time.Unix(int64(senddate), 0)
	} else if senddate, ok := result["senddate"].(time.Time); ok {
		article.SendDate = senddate
	}

	article.Keywords, _ = result["keywords"].(string)
	article.Description, _ = result["description"].(string)
	article.Filename, _ = result["filename"].(string)
	article.TemplateFile, _ = result["filename"].(string) // filename字段对应模板文件

	// 解析flag字段
	if flagStr, ok := result["flag"].(string); ok {
		article.IsTop, article.IsRecommend, article.IsHot = parseFlags(flagStr)
	}

	// 处理整数字段
	if arcrank, ok := result["arcrank"].(int64); ok {
		article.ArcRank = int(arcrank)
	} else if arcrank, ok := result["arcrank"].(int); ok {
		article.ArcRank = arcrank
	}

	if click, ok := result["click"].(int64); ok {
		article.Click = int(click)
	} else if click, ok := result["click"].(int); ok {
		article.Click = click
	}

	article.Body, _ = result["body"].(string)
	article.Content = article.Body // 设置Content字段与Body相同
	article.TypeName, _ = result["typename"].(string)
	article.TypeDir, _ = result["typedir"].(string)

	// 设置CategoryName字段以兼容模板
	article.CategoryName = article.TypeName

	return article, nil
}

// GetList 获取文章列表
func (m *ArticleModel) GetList(typeid int64, page, pageSize int) ([]*Article, int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "archives")
	qb.Select("a.id", "a.typeid", "a.title", "a.shorttitle", "a.color", "a.writer", "a.source", "a.litpic", "a.pubdate", "a.senddate", "a.keywords", "a.description", "a.filename", "a.flag", "a.arcrank", "a.click", "t.typename", "t.typedir")
	qb.From(m.db.TableName("archives") + " AS a")
	qb.LeftJoin(m.db.TableName("arctype")+" AS t", "a.typeid = t.id")

	// 添加条件
	qb.Where("a.arcrank > -1")
	if typeid > 0 {
		qb.Where("a.typeid = ?", typeid)
	}

	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("查询文章总数失败", "error", err)
		return nil, 0, err
	}

	// 设置分页
	offset := (page - 1) * pageSize
	qb.OrderBy("a.pubdate DESC")
	qb.Limit(pageSize, offset)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("查询文章列表失败", "error", err)
		return nil, 0, err
	}

	// 转换为文章对象
	articles := make([]*Article, 0, len(results))
	for _, result := range results {
		article := &Article{}

		// 处理ID字段 - 使用通用转换方法
		if idVal := result["id"]; idVal != nil {
			switch v := idVal.(type) {
			case int64:
				article.ID = v
			case int:
				article.ID = int64(v)
			case int32:
				article.ID = int64(v)
			case []uint8: // MySQL可能返回字节数组
				if len(v) > 0 {
					if id, err := strconv.ParseInt(string(v), 10, 64); err == nil {
						article.ID = id
					}
				}
			case string:
				if id, err := strconv.ParseInt(v, 10, 64); err == nil {
					article.ID = id
				}
			default:
				logger.Error("ID字段类型不支持", "value", idVal, "type", fmt.Sprintf("%T", idVal))
			}
		}

		// 处理TypeID字段 - 使用通用转换方法
		if typeidVal := result["typeid"]; typeidVal != nil {
			switch v := typeidVal.(type) {
			case int64:
				article.TypeID = v
			case int:
				article.TypeID = int64(v)
			case int32:
				article.TypeID = int64(v)
			case []uint8:
				if len(v) > 0 {
					if typeid, err := strconv.ParseInt(string(v), 10, 64); err == nil {
						article.TypeID = typeid
					}
				}
			case string:
				if typeid, err := strconv.ParseInt(v, 10, 64); err == nil {
					article.TypeID = typeid
				}
			}
		}
		article.Title, _ = result["title"].(string)
		article.ShortTitle, _ = result["shorttitle"].(string)
		article.Color, _ = result["color"].(string)
		article.Writer, _ = result["writer"].(string)
		article.Source, _ = result["source"].(string)
		article.LitPic, _ = result["litpic"].(string)

		// 处理日期 - 数据库中存储的是Unix时间戳
		if pubdate, ok := result["pubdate"].(int64); ok {
			article.PubDate = time.Unix(pubdate, 0)
		} else if pubdate, ok := result["pubdate"].(int); ok {
			article.PubDate = time.Unix(int64(pubdate), 0)
		} else if pubdate, ok := result["pubdate"].(time.Time); ok {
			article.PubDate = pubdate
		}

		if senddate, ok := result["senddate"].(int64); ok {
			article.SendDate = time.Unix(senddate, 0)
		} else if senddate, ok := result["senddate"].(int); ok {
			article.SendDate = time.Unix(int64(senddate), 0)
		} else if senddate, ok := result["senddate"].(time.Time); ok {
			article.SendDate = senddate
		}

		article.Keywords, _ = result["keywords"].(string)
		article.Description, _ = result["description"].(string)
		article.Filename, _ = result["filename"].(string)
		article.TemplateFile, _ = result["filename"].(string) // filename字段对应模板文件

		// 解析flag字段
		if flagStr, ok := result["flag"].(string); ok {
			article.IsTop, article.IsRecommend, article.IsHot = parseFlags(flagStr)
		}

		// 处理整数字段
		if arcrank, ok := result["arcrank"].(int64); ok {
			article.ArcRank = int(arcrank)
		} else if arcrank, ok := result["arcrank"].(int); ok {
			article.ArcRank = arcrank
		}

		if click, ok := result["click"].(int64); ok {
			article.Click = int(click)
		} else if click, ok := result["click"].(int); ok {
			article.Click = click
		}

		article.TypeName, _ = result["typename"].(string)
		article.TypeDir, _ = result["typedir"].(string)

		// 设置CategoryName字段以兼容模板
		article.CategoryName = article.TypeName

		articles = append(articles, article)
	}

	return articles, total, nil
}

// Create 创建文章
func (m *ArticleModel) Create(article *Article) (int64, error) {
	// 开始事务
	tx, err := m.db.Begin()
	if err != nil {
		logger.Error("开始事务失败", "error", err)
		return 0, err
	}
	defer tx.Rollback()

	// 构建flag字段
	var flags []string
	if article.IsTop == 1 {
		flags = append(flags, "c") // 置顶
	}
	if article.IsRecommend == 1 {
		flags = append(flags, "h") // 推荐
	}
	if article.IsHot == 1 {
		flags = append(flags, "p") // 热门
	}
	flagStr := strings.Join(flags, ",")

	// 转换时间为Unix时间戳，处理无效时间
	var pubdate, senddate int64

	// 验证并转换pubdate
	if article.PubDate.IsZero() || article.PubDate.Year() < 1970 {
		pubdate = time.Now().Unix()
	} else {
		pubdate = article.PubDate.Unix()
	}

	// 验证并转换senddate
	if article.SendDate.IsZero() || article.SendDate.Year() < 1970 {
		senddate = time.Now().Unix()
	} else {
		senddate = article.SendDate.Unix()
	}

	// 执行插入
	result, err := tx.Exec(
		"INSERT INTO "+m.db.TableName("archives")+" (typeid, title, shorttitle, color, writer, source, litpic, pubdate, senddate, keywords, description, filename, flag, arcrank, click) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		article.TypeID, article.Title, article.ShortTitle, article.Color, article.Writer, article.Source, article.LitPic, pubdate, senddate, article.Keywords, article.Description, article.Filename, flagStr, article.ArcRank, article.Click,
	)
	if err != nil {
		logger.Error("插入文章主表失败", "error", err)
		return 0, err
	}

	// 获取插入ID
	id, err := result.LastInsertId()
	if err != nil {
		logger.Error("获取插入ID失败", "error", err)
		return 0, err
	}

	// 插入附加表
	_, err = tx.Exec(
		"INSERT INTO "+m.db.TableName("addonarticle")+" (aid, body) VALUES (?, ?)",
		id, article.Body,
	)
	if err != nil {
		logger.Error("插入文章附加表失败", "error", err)
		return 0, err
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		logger.Error("提交事务失败", "error", err)
		return 0, err
	}

	return id, nil
}

// Update 更新文章
func (m *ArticleModel) Update(article *Article) error {
	// 开始事务
	tx, err := m.db.Begin()
	if err != nil {
		logger.Error("开始事务失败", "error", err)
		return err
	}
	defer tx.Rollback()

	// 构建flag字段
	var flags []string
	if article.IsTop == 1 {
		flags = append(flags, "c") // 置顶
	}
	if article.IsRecommend == 1 {
		flags = append(flags, "h") // 推荐
	}
	if article.IsHot == 1 {
		flags = append(flags, "p") // 热门
	}
	flagStr := strings.Join(flags, ",")

	// 转换时间为Unix时间戳，处理无效时间
	var pubdate, senddate int64

	// 验证并转换pubdate
	if article.PubDate.IsZero() || article.PubDate.Year() < 1970 {
		pubdate = time.Now().Unix()
	} else {
		pubdate = article.PubDate.Unix()
	}

	// 验证并转换senddate
	if article.SendDate.IsZero() || article.SendDate.Year() < 1970 {
		senddate = time.Now().Unix()
	} else {
		senddate = article.SendDate.Unix()
	}

	// 更新主表
	_, err = tx.Exec(
		"UPDATE "+m.db.TableName("archives")+" SET typeid=?, title=?, shorttitle=?, color=?, writer=?, source=?, litpic=?, pubdate=?, senddate=?, keywords=?, description=?, filename=?, flag=?, arcrank=?, click=? WHERE id=?",
		article.TypeID, article.Title, article.ShortTitle, article.Color, article.Writer, article.Source, article.LitPic, pubdate, senddate, article.Keywords, article.Description, article.Filename, flagStr, article.ArcRank, article.Click, article.ID,
	)
	if err != nil {
		logger.Error("更新文章主表失败", "error", err)
		return err
	}

	// 更新附加表
	_, err = tx.Exec(
		"UPDATE "+m.db.TableName("addonarticle")+" SET body=? WHERE aid=?",
		article.Body, article.ID,
	)
	if err != nil {
		logger.Error("更新文章附加表失败", "error", err)
		return err
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		logger.Error("提交事务失败", "error", err)
		return err
	}

	return nil
}

// Delete 删除文章
func (m *ArticleModel) Delete(id int64) error {
	// 开始事务
	tx, err := m.db.Begin()
	if err != nil {
		logger.Error("开始事务失败", "error", err)
		return err
	}
	defer tx.Rollback()

	// 删除主表
	_, err = tx.Exec(
		"DELETE FROM "+m.db.TableName("archives")+" WHERE id=?",
		id,
	)
	if err != nil {
		logger.Error("删除文章主表失败", "error", err)
		return err
	}

	// 删除附加表
	_, err = tx.Exec(
		"DELETE FROM "+m.db.TableName("addonarticle")+" WHERE aid=?",
		id,
	)
	if err != nil {
		logger.Error("删除文章附加表失败", "error", err)
		return err
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		logger.Error("提交事务失败", "error", err)
		return err
	}

	return nil
}

// Search 搜索文章
func (m *ArticleModel) Search(keyword string, page, pageSize int) ([]*Article, int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "archives")
	qb.Select("a.*", "t.typename", "t.typedir")
	qb.From(m.db.TableName("archives") + " AS a")
	qb.LeftJoin(m.db.TableName("arctype")+" AS t", "a.typeid = t.id")

	// 添加搜索条件
	qb.Where("a.arcrank > -1")
	qb.Where("(a.title LIKE ? OR a.keywords LIKE ? OR a.description LIKE ?)",
		"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")

	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("搜索文章总数失败", "error", err)
		return nil, 0, err
	}

	// 设置分页
	offset := (page - 1) * pageSize
	qb.OrderBy("a.pubdate DESC")
	qb.Limit(pageSize, offset)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("搜索文章列表失败", "error", err)
		return nil, 0, err
	}

	// 转换为文章对象
	articles := make([]*Article, 0, len(results))
	for _, result := range results {
		article := &Article{}
		article.ID, _ = result["id"].(int64)
		article.TypeID, _ = result["typeid"].(int64)
		article.Title, _ = result["title"].(string)
		article.ShortTitle, _ = result["shorttitle"].(string)
		article.Color, _ = result["color"].(string)
		article.Writer, _ = result["writer"].(string)
		article.Source, _ = result["source"].(string)
		article.LitPic, _ = result["litpic"].(string)

		// 处理日期
		if pubdate, ok := result["pubdate"].(time.Time); ok {
			article.PubDate = pubdate
		}
		if senddate, ok := result["senddate"].(time.Time); ok {
			article.SendDate = senddate
		}

		article.Keywords, _ = result["keywords"].(string)
		article.Description, _ = result["description"].(string)
		article.Filename, _ = result["filename"].(string)

		// 解析flag字段
		if flagStr, ok := result["flag"].(string); ok {
			article.IsTop, article.IsRecommend, article.IsHot = parseFlags(flagStr)
		}

		// 处理整数字段
		if arcrank, ok := result["arcrank"].(int64); ok {
			article.ArcRank = int(arcrank)
		} else if arcrank, ok := result["arcrank"].(int); ok {
			article.ArcRank = arcrank
		}

		if click, ok := result["click"].(int64); ok {
			article.Click = int(click)
		} else if click, ok := result["click"].(int); ok {
			article.Click = click
		}

		article.TypeName, _ = result["typename"].(string)
		article.TypeDir, _ = result["typedir"].(string)

		articles = append(articles, article)
	}

	return articles, total, nil
}

// GetByTypeID 根据栏目ID获取文章列表
func (m *ArticleModel) GetByTypeID(typeid int64, page, pageSize int) ([]*Article, int, error) {
	return m.GetList(typeid, page, pageSize)
}

// IncrementClick 增加文章点击量
func (m *ArticleModel) IncrementClick(id int64) error {
	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("archives")+" SET click = click + 1 WHERE id = ?",
		id,
	)
	if err != nil {
		logger.Error("增加文章点击量失败", "id", id, "error", err)
		return err
	}

	return nil
}

// GetPrevArticle 获取上一篇文章
func (m *ArticleModel) GetPrevArticle(id int64, typeid int64) (*Article, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "archives")
	qb.Select("a.*", "t.typename", "t.typedir")
	qb.From(m.db.TableName("archives") + " AS a")
	qb.LeftJoin(m.db.TableName("arctype")+" AS t", "a.typeid = t.id")
	qb.Where("a.arcrank > -1")
	qb.Where("a.id < ?", id)

	if typeid > 0 {
		qb.Where("a.typeid = ?", typeid)
	}

	qb.OrderBy("a.id DESC")
	qb.Limit(1, 0)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("查询上一篇文章失败", "id", id, "error", err)
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	// 转换为文章对象
	article := &Article{}
	article.ID, _ = result["id"].(int64)
	article.TypeID, _ = result["typeid"].(int64)
	article.Title, _ = result["title"].(string)
	article.ShortTitle, _ = result["shorttitle"].(string)
	article.Color, _ = result["color"].(string)
	article.Writer, _ = result["writer"].(string)
	article.Source, _ = result["source"].(string)
	article.LitPic, _ = result["litpic"].(string)

	// 处理日期
	if pubdate, ok := result["pubdate"].(time.Time); ok {
		article.PubDate = pubdate
	}
	if senddate, ok := result["senddate"].(time.Time); ok {
		article.SendDate = senddate
	}

	article.Keywords, _ = result["keywords"].(string)
	article.Description, _ = result["description"].(string)
	article.Filename, _ = result["filename"].(string)

	// 处理整数字段
	if istop, ok := result["istop"].(int64); ok {
		article.IsTop = int(istop)
	}
	if isrecommend, ok := result["isrecommend"].(int64); ok {
		article.IsRecommend = int(isrecommend)
	}
	if ishot, ok := result["ishot"].(int64); ok {
		article.IsHot = int(ishot)
	}
	if arcrank, ok := result["arcrank"].(int64); ok {
		article.ArcRank = int(arcrank)
	}
	if click, ok := result["click"].(int64); ok {
		article.Click = int(click)
	}

	article.TypeName, _ = result["typename"].(string)
	article.TypeDir, _ = result["typedir"].(string)

	return article, nil
}

// GetNextArticle 获取下一篇文章
func (m *ArticleModel) GetNextArticle(id int64, typeid int64) (*Article, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "archives")
	qb.Select("a.*", "t.typename", "t.typedir")
	qb.From(m.db.TableName("archives") + " AS a")
	qb.LeftJoin(m.db.TableName("arctype")+" AS t", "a.typeid = t.id")
	qb.Where("a.arcrank > -1")
	qb.Where("a.id > ?", id)

	if typeid > 0 {
		qb.Where("a.typeid = ?", typeid)
	}

	qb.OrderBy("a.id ASC")
	qb.Limit(1, 0)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("查询下一篇文章失败", "id", id, "error", err)
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	// 转换为文章对象
	article := &Article{}
	article.ID, _ = result["id"].(int64)
	article.TypeID, _ = result["typeid"].(int64)
	article.Title, _ = result["title"].(string)
	article.ShortTitle, _ = result["shorttitle"].(string)
	article.Color, _ = result["color"].(string)
	article.Writer, _ = result["writer"].(string)
	article.Source, _ = result["source"].(string)
	article.LitPic, _ = result["litpic"].(string)

	// 处理日期
	if pubdate, ok := result["pubdate"].(time.Time); ok {
		article.PubDate = pubdate
	}
	if senddate, ok := result["senddate"].(time.Time); ok {
		article.SendDate = senddate
	}

	article.Keywords, _ = result["keywords"].(string)
	article.Description, _ = result["description"].(string)
	article.Filename, _ = result["filename"].(string)

	// 处理整数字段
	if istop, ok := result["istop"].(int64); ok {
		article.IsTop = int(istop)
	}
	if isrecommend, ok := result["isrecommend"].(int64); ok {
		article.IsRecommend = int(isrecommend)
	}
	if ishot, ok := result["ishot"].(int64); ok {
		article.IsHot = int(ishot)
	}
	if arcrank, ok := result["arcrank"].(int64); ok {
		article.ArcRank = int(arcrank)
	}
	if click, ok := result["click"].(int64); ok {
		article.Click = int(click)
	}

	article.TypeName, _ = result["typename"].(string)
	article.TypeDir, _ = result["typedir"].(string)

	return article, nil
}

// GetRelatedArticles 获取相关文章
func (m *ArticleModel) GetRelatedArticles(keywords string, id int64, limit int) ([]*Article, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "archives")
	qb.Select("a.*", "t.typename", "t.typedir")
	qb.From(m.db.TableName("archives") + " AS a")
	qb.LeftJoin(m.db.TableName("arctype")+" AS t", "a.typeid = t.id")
	qb.Where("a.arcrank > -1")
	qb.Where("a.id <> ?", id)

	// 如果有关键词，则按关键词匹配
	if keywords != "" {
		keywordList := strings.Split(keywords, ",")
		for _, keyword := range keywordList {
			if keyword != "" {
				qb.Where("a.keywords LIKE ?", "%"+keyword+"%")
			}
		}
	}

	qb.OrderBy("a.pubdate DESC")
	qb.Limit(limit, 0)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("查询相关文章失败", "id", id, "error", err)
		return nil, err
	}

	// 转换为文章对象
	articles := make([]*Article, 0, len(results))
	for _, result := range results {
		article := &Article{}
		article.ID, _ = result["id"].(int64)
		article.TypeID, _ = result["typeid"].(int64)
		article.Title, _ = result["title"].(string)
		article.ShortTitle, _ = result["shorttitle"].(string)
		article.Color, _ = result["color"].(string)
		article.Writer, _ = result["writer"].(string)
		article.Source, _ = result["source"].(string)
		article.LitPic, _ = result["litpic"].(string)

		// 处理日期
		if pubdate, ok := result["pubdate"].(time.Time); ok {
			article.PubDate = pubdate
		}
		if senddate, ok := result["senddate"].(time.Time); ok {
			article.SendDate = senddate
		}

		article.Keywords, _ = result["keywords"].(string)
		article.Description, _ = result["description"].(string)
		article.Filename, _ = result["filename"].(string)

		// 处理整数字段
		if istop, ok := result["istop"].(int64); ok {
			article.IsTop = int(istop)
		}
		if isrecommend, ok := result["isrecommend"].(int64); ok {
			article.IsRecommend = int(isrecommend)
		}
		if ishot, ok := result["ishot"].(int64); ok {
			article.IsHot = int(ishot)
		}
		if arcrank, ok := result["arcrank"].(int64); ok {
			article.ArcRank = int(arcrank)
		}
		if click, ok := result["click"].(int64); ok {
			article.Click = int(click)
		}

		article.TypeName, _ = result["typename"].(string)
		article.TypeDir, _ = result["typedir"].(string)

		articles = append(articles, article)
	}

	return articles, nil
}

// IncrementCommentCount 增加文章评论数
func (m *ArticleModel) IncrementCommentCount(id int64) error {
	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("archives")+" SET commentcount = commentcount + 1 WHERE id = ?",
		id,
	)
	if err != nil {
		logger.Error("增加文章评论数失败", "id", id, "error", err)
		return err
	}

	return nil
}

// DecrementCommentCount 减少文章评论数
func (m *ArticleModel) DecrementCommentCount(id int64) error {
	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("archives")+" SET commentcount = GREATEST(commentcount - 1, 0) WHERE id = ?",
		id,
	)
	if err != nil {
		logger.Error("减少文章评论数失败", "id", id, "error", err)
		return err
	}

	return nil
}

// GetMemberArticles 获取会员文章
func (m *ArticleModel) GetMemberArticles(memberID int64, page, pageSize int) ([]*Article, int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "archives")
	qb.Select("a.*", "t.typename", "t.typedir")
	qb.From(m.db.TableName("archives") + " AS a")
	qb.LeftJoin(m.db.TableName("arctype")+" AS t", "a.typeid = t.id")
	qb.Where("a.memberid = ?", memberID)
	qb.Where("a.arcrank > -1")

	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("获取会员文章总数失败", "memberID", memberID, "error", err)
		return nil, 0, err
	}

	// 设置分页
	offset := (page - 1) * pageSize
	qb.OrderBy("a.pubdate DESC")
	qb.Limit(pageSize, offset)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取会员文章列表失败", "memberID", memberID, "error", err)
		return nil, 0, err
	}

	// 转换为文章对象
	articles := make([]*Article, 0, len(results))
	for _, result := range results {
		article := &Article{}
		article.ID, _ = result["id"].(int64)
		article.TypeID, _ = result["typeid"].(int64)
		article.Title, _ = result["title"].(string)
		article.ShortTitle, _ = result["shorttitle"].(string)
		article.Color, _ = result["color"].(string)
		article.Writer, _ = result["writer"].(string)
		article.Source, _ = result["source"].(string)
		article.LitPic, _ = result["litpic"].(string)

		// 处理日期
		if pubdate, ok := result["pubdate"].(time.Time); ok {
			article.PubDate = pubdate
		}
		if senddate, ok := result["senddate"].(time.Time); ok {
			article.SendDate = senddate
		}

		article.Keywords, _ = result["keywords"].(string)
		article.Description, _ = result["description"].(string)
		article.Filename, _ = result["filename"].(string)

		// 处理整数字段
		if istop, ok := result["istop"].(int64); ok {
			article.IsTop = int(istop)
		}
		if isrecommend, ok := result["isrecommend"].(int64); ok {
			article.IsRecommend = int(isrecommend)
		}
		if ishot, ok := result["ishot"].(int64); ok {
			article.IsHot = int(ishot)
		}
		if arcrank, ok := result["arcrank"].(int64); ok {
			article.ArcRank = int(arcrank)
		}
		if click, ok := result["click"].(int64); ok {
			article.Click = int(click)
		}

		article.TypeName, _ = result["typename"].(string)
		article.TypeDir, _ = result["typedir"].(string)

		articles = append(articles, article)
	}

	return articles, total, nil
}

// GetByMemberID 获取会员文章
func (m *ArticleModel) GetByMemberID(memberID int64, page, pageSize int) ([]*Article, int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "archives")
	qb.Select("a.*", "t.typename", "t.typedir")
	qb.From(m.db.TableName("archives") + " AS a")
	qb.LeftJoin(m.db.TableName("arctype")+" AS t", "a.typeid = t.id")
	qb.Where("a.memberid = ?", memberID)
	qb.Where("a.arcrank > -1")

	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("获取会员文章总数失败", "memberID", memberID, "error", err)
		return nil, 0, err
	}

	// 设置分页
	offset := (page - 1) * pageSize
	qb.OrderBy("a.pubdate DESC")
	qb.Limit(pageSize, offset)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取会员文章列表失败", "memberID", memberID, "error", err)
		return nil, 0, err
	}

	// 转换为文章对象
	articles := make([]*Article, 0, len(results))
	for _, result := range results {
		article := &Article{}
		article.ID, _ = result["id"].(int64)
		article.TypeID, _ = result["typeid"].(int64)
		article.Title, _ = result["title"].(string)
		article.ShortTitle, _ = result["shorttitle"].(string)
		article.Color, _ = result["color"].(string)
		article.Writer, _ = result["writer"].(string)
		article.Source, _ = result["source"].(string)
		article.LitPic, _ = result["litpic"].(string)

		// 处理日期
		if pubdate, ok := result["pubdate"].(time.Time); ok {
			article.PubDate = pubdate
		}
		if senddate, ok := result["senddate"].(time.Time); ok {
			article.SendDate = senddate
		}

		article.Keywords, _ = result["keywords"].(string)
		article.Description, _ = result["description"].(string)
		article.Filename, _ = result["filename"].(string)

		// 处理整数字段
		if istop, ok := result["istop"].(int64); ok {
			article.IsTop = int(istop)
		}
		if isrecommend, ok := result["isrecommend"].(int64); ok {
			article.IsRecommend = int(isrecommend)
		}
		if ishot, ok := result["ishot"].(int64); ok {
			article.IsHot = int(ishot)
		}
		if arcrank, ok := result["arcrank"].(int64); ok {
			article.ArcRank = int(arcrank)
		}
		if click, ok := result["click"].(int64); ok {
			article.Click = int(click)
		}

		article.TypeName, _ = result["typename"].(string)
		article.TypeDir, _ = result["typedir"].(string)

		articles = append(articles, article)
	}

	return articles, total, nil
}

// GetByTag 根据标签获取文章
func (m *ArticleModel) GetByTag(tag string, page, pageSize int) ([]*Article, int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "archives")
	qb.Select("a.*", "t.typename", "t.typedir")
	qb.From(m.db.TableName("archives") + " AS a")
	qb.LeftJoin(m.db.TableName("arctype")+" AS t", "a.typeid = t.id")
	qb.Where("a.arcrank > -1")
	qb.Where("a.keywords LIKE ?", "%"+tag+"%")

	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("查询标签文章总数失败", "tag", tag, "error", err)
		return nil, 0, err
	}

	// 设置分页
	offset := (page - 1) * pageSize
	qb.OrderBy("a.pubdate DESC")
	qb.Limit(pageSize, offset)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("查询标签文章列表失败", "tag", tag, "error", err)
		return nil, 0, err
	}

	// 转换为文章对象
	articles := make([]*Article, 0, len(results))
	for _, result := range results {
		article := &Article{}
		article.ID, _ = result["id"].(int64)
		article.TypeID, _ = result["typeid"].(int64)
		article.Title, _ = result["title"].(string)
		article.ShortTitle, _ = result["shorttitle"].(string)
		article.Color, _ = result["color"].(string)
		article.Writer, _ = result["writer"].(string)
		article.Source, _ = result["source"].(string)
		article.LitPic, _ = result["litpic"].(string)

		// 处理日期
		if pubdate, ok := result["pubdate"].(time.Time); ok {
			article.PubDate = pubdate
		}
		if senddate, ok := result["senddate"].(time.Time); ok {
			article.SendDate = senddate
		}

		article.Keywords, _ = result["keywords"].(string)
		article.Description, _ = result["description"].(string)
		article.Filename, _ = result["filename"].(string)

		// 处理整数字段
		if istop, ok := result["istop"].(int64); ok {
			article.IsTop = int(istop)
		}
		if isrecommend, ok := result["isrecommend"].(int64); ok {
			article.IsRecommend = int(isrecommend)
		}
		if ishot, ok := result["ishot"].(int64); ok {
			article.IsHot = int(ishot)
		}
		if arcrank, ok := result["arcrank"].(int64); ok {
			article.ArcRank = int(arcrank)
		}
		if click, ok := result["click"].(int64); ok {
			article.Click = int(click)
		}

		article.TypeName, _ = result["typename"].(string)
		article.TypeDir, _ = result["typedir"].(string)

		articles = append(articles, article)
	}

	return articles, total, nil
}

// GetByDateRange 根据日期范围获取文章
func (m *ArticleModel) GetByDateRange(typeID int64, startDate, endDate time.Time) ([]*Article, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "archives")
	qb.Select("a.*", "t.typename", "t.typedir")
	qb.From(m.db.TableName("archives") + " AS a")
	qb.LeftJoin(m.db.TableName("arctype")+" AS t", "a.typeid = t.id")

	// 添加条件
	if typeID > 0 {
		qb.Where("a.typeid = ?", typeID)
	}

	qb.Where("a.pubdate >= ?", startDate)
	qb.Where("a.pubdate <= ?", endDate)
	qb.Where("a.arcrank > 0")
	qb.OrderBy("a.pubdate DESC")

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("查询文章列表失败", "error", err)
		return nil, err
	}

	// 转换为文章对象
	articles := make([]*Article, 0, len(results))
	for _, result := range results {
		article := &Article{}
		article.ID, _ = result["id"].(int64)
		article.TypeID, _ = result["typeid"].(int64)
		article.Title, _ = result["title"].(string)
		article.ShortTitle, _ = result["shorttitle"].(string)
		article.Color, _ = result["color"].(string)
		article.Writer, _ = result["writer"].(string)
		article.Source, _ = result["source"].(string)
		article.LitPic, _ = result["litpic"].(string)

		// 处理日期
		if pubdate, ok := result["pubdate"].(time.Time); ok {
			article.PubDate = pubdate
		}
		if senddate, ok := result["senddate"].(time.Time); ok {
			article.SendDate = senddate
		}

		article.Keywords, _ = result["keywords"].(string)
		article.Description, _ = result["description"].(string)
		article.Filename, _ = result["filename"].(string)

		// 处理整数字段
		if istop, ok := result["istop"].(int64); ok {
			article.IsTop = int(istop)
		}
		if isrecommend, ok := result["isrecommend"].(int64); ok {
			article.IsRecommend = int(isrecommend)
		}
		if ishot, ok := result["ishot"].(int64); ok {
			article.IsHot = int(ishot)
		}
		if arcrank, ok := result["arcrank"].(int64); ok {
			article.ArcRank = int(arcrank)
		}
		if click, ok := result["click"].(int64); ok {
			article.Click = int(click)
		}

		article.TypeName, _ = result["typename"].(string)
		article.TypeDir, _ = result["typedir"].(string)

		articles = append(articles, article)
	}

	return articles, nil
}

// GetCount 获取文章总数
func (m *ArticleModel) GetCount() (int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "archives")
	qb.Where("arcrank > 0")

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("获取文章总数失败", "error", err)
		return 0, err
	}

	return count, nil
}

// GetAll 获取所有文章
func (m *ArticleModel) GetAll(page, pageSize int) ([]*Article, int, error) {
	return m.GetList(0, page, pageSize)
}

// GetLatestArticles 获取最新文章
func (m *ArticleModel) GetLatestArticles(limit int) ([]*Article, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "archives")
	qb.Select("a.*", "t.typename", "t.typedir")
	qb.From(m.db.TableName("archives") + " AS a")
	qb.LeftJoin(m.db.TableName("arctype")+" AS t", "a.typeid = t.id")
	qb.Where("a.arcrank > -1")
	qb.OrderBy("a.pubdate DESC")
	qb.Limit(limit, 0)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取最新文章失败", "error", err)
		return nil, err
	}

	// 转换为文章对象
	articles := make([]*Article, 0, len(results))
	for _, result := range results {
		article := &Article{}
		article.ID, _ = result["id"].(int64)
		article.TypeID, _ = result["typeid"].(int64)
		article.Title, _ = result["title"].(string)
		article.ShortTitle, _ = result["shorttitle"].(string)
		article.Color, _ = result["color"].(string)
		article.Writer, _ = result["writer"].(string)
		article.Source, _ = result["source"].(string)
		article.LitPic, _ = result["litpic"].(string)

		// 处理日期
		if pubdate, ok := result["pubdate"].(time.Time); ok {
			article.PubDate = pubdate
		}
		if senddate, ok := result["senddate"].(time.Time); ok {
			article.SendDate = senddate
		}

		article.Keywords, _ = result["keywords"].(string)
		article.Description, _ = result["description"].(string)
		article.Filename, _ = result["filename"].(string)

		// 处理整数字段
		if istop, ok := result["istop"].(int64); ok {
			article.IsTop = int(istop)
		}
		if isrecommend, ok := result["isrecommend"].(int64); ok {
			article.IsRecommend = int(isrecommend)
		}
		if ishot, ok := result["ishot"].(int64); ok {
			article.IsHot = int(ishot)
		}
		if arcrank, ok := result["arcrank"].(int64); ok {
			article.ArcRank = int(arcrank)
		}
		if click, ok := result["click"].(int64); ok {
			article.Click = int(click)
		}

		article.TypeName, _ = result["typename"].(string)
		article.TypeDir, _ = result["typedir"].(string)

		articles = append(articles, article)
	}

	return articles, nil
}

// GetTotalCount 获取文章总数
func (m *ArticleModel) GetTotalCount() (int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "archives")
	qb.Where("arcrank > -1")

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("获取文章总数失败", "error", err)
		return 0, err
	}

	return count, nil
}

// GetTodayCount 获取今日新增文章数
func (m *ArticleModel) GetTodayCount() (int, error) {
	// 获取今日开始时间
	today := time.Now().Format("2006-01-02")
	todayStart := today + " 00:00:00"

	// 构建查询
	qb := database.NewQueryBuilder(m.db, "archives")
	qb.Where("arcrank > -1")
	qb.Where("pubdate >= ?", todayStart)

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("获取今日文章数失败", "error", err)
		return 0, err
	}

	return count, nil
}

// GetPendingCount 获取待审核文章数
func (m *ArticleModel) GetPendingCount() (int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "archives")
	qb.Where("arcrank = -1") // 待审核状态

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("获取待审核文章数失败", "error", err)
		return 0, err
	}

	return count, nil
}
