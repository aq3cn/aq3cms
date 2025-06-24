package model

import (
	"fmt"
	"time"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// Download 下载模型
type Download struct {
	ID          int64     `json:"id"`
	TypeID      int64     `json:"typeid"`
	Title       string    `json:"title"`
	ShortTitle  string    `json:"shorttitle"`
	Color       string    `json:"color"`
	Writer      string    `json:"writer"`
	Source      string    `json:"source"`
	LitPic      string    `json:"litpic"`
	PubDate     time.Time `json:"pubdate"`
	SendDate    time.Time `json:"senddate"`
	Keywords    string    `json:"keywords"`
	Description string    `json:"description"`
	Filename    string    `json:"filename"`
	IsTop       int       `json:"istop"`
	IsRecommend int       `json:"isrecommend"`
	IsHot       int       `json:"ishot"`
	ArcRank     int       `json:"arcrank"`
	Click       int       `json:"click"`
	Body        string    `json:"body"`
	TypeName    string    `json:"typename"`
	TypeDir     string    `json:"typedir"`
	// 下载特有字段
	SoftName     string   `json:"softname"`
	SoftVersion  string   `json:"softversion"`
	SoftLanguage string   `json:"softlanguage"`
	SoftType     string   `json:"softtype"`
	SoftSize     string   `json:"softsize"`
	SoftOS       string   `json:"softos"`
	SoftDeveloper string  `json:"softdeveloper"`
	SoftLicense  string   `json:"softlicense"`
	SoftScore    float64  `json:"softscore"`
	SoftURL      string   `json:"softurl"`
	SoftMirrURL  string   `json:"softmirrurl"`
	Screenshots  []string `json:"screenshots"`
	DownCount    int      `json:"downcount"`
}

// DownloadModel 下载模型操作
type DownloadModel struct {
	db *database.DB
}

// NewDownloadModel 创建下载模型
func NewDownloadModel(db *database.DB) *DownloadModel {
	return &DownloadModel{
		db: db,
	}
}

// GetByID 根据ID获取下载
func (m *DownloadModel) GetByID(id int64) (*Download, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "archives")
	qb.Select("a.*", "t.typename", "t.typedir", "ad.body", "d.*")
	qb.LeftJoin(m.db.TableName("arctype")+" AS t", "a.typeid = t.id")
	qb.LeftJoin(m.db.TableName("addonarticle")+" AS ad", "a.id = ad.aid")
	qb.LeftJoin(m.db.TableName("addondownload")+" AS d", "a.id = d.aid")
	qb.Where("a.id = ?", id)
	qb.Where("a.channel = 3") // 下载频道ID
	
	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("查询下载失败", "id", id, "error", err)
		return nil, err
	}
	
	if result == nil {
		return nil, fmt.Errorf("下载不存在")
	}
	
	// 转换为下载对象
	download := &Download{}
	download.ID, _ = result["id"].(int64)
	download.TypeID, _ = result["typeid"].(int64)
	download.Title, _ = result["title"].(string)
	download.ShortTitle, _ = result["shorttitle"].(string)
	download.Color, _ = result["color"].(string)
	download.Writer, _ = result["writer"].(string)
	download.Source, _ = result["source"].(string)
	download.LitPic, _ = result["litpic"].(string)
	
	// 处理日期
	if pubdate, ok := result["pubdate"].(time.Time); ok {
		download.PubDate = pubdate
	}
	if senddate, ok := result["senddate"].(time.Time); ok {
		download.SendDate = senddate
	}
	
	download.Keywords, _ = result["keywords"].(string)
	download.Description, _ = result["description"].(string)
	download.Filename, _ = result["filename"].(string)
	
	// 处理整数字段
	if istop, ok := result["istop"].(int64); ok {
		download.IsTop = int(istop)
	}
	if isrecommend, ok := result["isrecommend"].(int64); ok {
		download.IsRecommend = int(isrecommend)
	}
	if ishot, ok := result["ishot"].(int64); ok {
		download.IsHot = int(ishot)
	}
	if arcrank, ok := result["arcrank"].(int64); ok {
		download.ArcRank = int(arcrank)
	}
	if click, ok := result["click"].(int64); ok {
		download.Click = int(click)
	}
	
	download.Body, _ = result["body"].(string)
	download.TypeName, _ = result["typename"].(string)
	download.TypeDir, _ = result["typedir"].(string)
	
	// 处理下载特有字段
	download.SoftName, _ = result["softname"].(string)
	download.SoftVersion, _ = result["softversion"].(string)
	download.SoftLanguage, _ = result["softlanguage"].(string)
	download.SoftType, _ = result["softtype"].(string)
	download.SoftSize, _ = result["softsize"].(string)
	download.SoftOS, _ = result["softos"].(string)
	download.SoftDeveloper, _ = result["softdeveloper"].(string)
	download.SoftLicense, _ = result["softlicense"].(string)
	
	// 处理评分
	if softscore, ok := result["softscore"].(float64); ok {
		download.SoftScore = softscore
	} else if softscoreStr, ok := result["softscore"].(string); ok {
		fmt.Sscanf(softscoreStr, "%f", &download.SoftScore)
	}
	
	download.SoftURL, _ = result["softurl"].(string)
	download.SoftMirrURL, _ = result["softmirrurl"].(string)
	
	// 处理下载次数
	if downcount, ok := result["downcount"].(int64); ok {
		download.DownCount = int(downcount)
	} else if downcountStr, ok := result["downcount"].(string); ok {
		fmt.Sscanf(downcountStr, "%d", &download.DownCount)
	}
	
	// 获取截图
	screenshots, err := m.getDownloadScreenshots(id)
	if err != nil {
		logger.Error("获取下载截图失败", "id", id, "error", err)
	}
	download.Screenshots = screenshots
	
	return download, nil
}

// GetList 获取下载列表
func (m *DownloadModel) GetList(typeid int64, page, pageSize int) ([]*Download, int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "archives")
	qb.Select("a.*", "t.typename", "t.typedir", "d.softname", "d.softversion", "d.softsize", "d.softscore", "d.downcount")
	qb.LeftJoin(m.db.TableName("arctype")+" AS t", "a.typeid = t.id")
	qb.LeftJoin(m.db.TableName("addondownload")+" AS d", "a.id = d.aid")
	
	// 添加条件
	qb.Where("a.arcrank > -1")
	qb.Where("a.channel = 3") // 下载频道ID
	if typeid > 0 {
		qb.Where("a.typeid = ?", typeid)
	}
	
	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("查询下载总数失败", "error", err)
		return nil, 0, err
	}
	
	// 设置分页
	offset := (page - 1) * pageSize
	qb.OrderBy("a.pubdate DESC")
	qb.Limit(pageSize, offset)
	
	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("查询下载列表失败", "error", err)
		return nil, 0, err
	}
	
	// 转换为下载对象
	downloads := make([]*Download, 0, len(results))
	for _, result := range results {
		download := &Download{}
		download.ID, _ = result["id"].(int64)
		download.TypeID, _ = result["typeid"].(int64)
		download.Title, _ = result["title"].(string)
		download.ShortTitle, _ = result["shorttitle"].(string)
		download.Color, _ = result["color"].(string)
		download.Writer, _ = result["writer"].(string)
		download.Source, _ = result["source"].(string)
		download.LitPic, _ = result["litpic"].(string)
		
		// 处理日期
		if pubdate, ok := result["pubdate"].(time.Time); ok {
			download.PubDate = pubdate
		}
		if senddate, ok := result["senddate"].(time.Time); ok {
			download.SendDate = senddate
		}
		
		download.Keywords, _ = result["keywords"].(string)
		download.Description, _ = result["description"].(string)
		download.Filename, _ = result["filename"].(string)
		
		// 处理整数字段
		if istop, ok := result["istop"].(int64); ok {
			download.IsTop = int(istop)
		}
		if isrecommend, ok := result["isrecommend"].(int64); ok {
			download.IsRecommend = int(isrecommend)
		}
		if ishot, ok := result["ishot"].(int64); ok {
			download.IsHot = int(ishot)
		}
		if arcrank, ok := result["arcrank"].(int64); ok {
			download.ArcRank = int(arcrank)
		}
		if click, ok := result["click"].(int64); ok {
			download.Click = int(click)
		}
		
		download.TypeName, _ = result["typename"].(string)
		download.TypeDir, _ = result["typedir"].(string)
		
		// 处理下载特有字段
		download.SoftName, _ = result["softname"].(string)
		download.SoftVersion, _ = result["softversion"].(string)
		download.SoftSize, _ = result["softsize"].(string)
		
		// 处理评分
		if softscore, ok := result["softscore"].(float64); ok {
			download.SoftScore = softscore
		} else if softscoreStr, ok := result["softscore"].(string); ok {
			fmt.Sscanf(softscoreStr, "%f", &download.SoftScore)
		}
		
		// 处理下载次数
		if downcount, ok := result["downcount"].(int64); ok {
			download.DownCount = int(downcount)
		} else if downcountStr, ok := result["downcount"].(string); ok {
			fmt.Sscanf(downcountStr, "%d", &download.DownCount)
		}
		
		downloads = append(downloads, download)
	}
	
	return downloads, total, nil
}

// Create 创建下载
func (m *DownloadModel) Create(download *Download) (int64, error) {
	// 开始事务
	tx, err := m.db.Begin()
	if err != nil {
		logger.Error("开始事务失败", "error", err)
		return 0, err
	}
	defer tx.Rollback()
	
	// 插入主表
	_, err = tx.Exec(
		"INSERT INTO "+m.db.TableName("archives")+" (typeid, channel, title, shorttitle, color, writer, source, litpic, pubdate, senddate, keywords, description, filename, istop, isrecommend, ishot, arcrank, click) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		download.TypeID, 3, download.Title, download.ShortTitle, download.Color, download.Writer, download.Source, download.LitPic, download.PubDate, download.SendDate, download.Keywords, download.Description, download.Filename, download.IsTop, download.IsRecommend, download.IsHot, download.ArcRank, download.Click,
	)
	if err != nil {
		logger.Error("插入下载主表失败", "error", err)
		return 0, err
	}
	
	// 获取插入ID
	var id int64
	err = tx.QueryRow("SELECT LAST_INSERT_ID()").Scan(&id)
	if err != nil {
		logger.Error("获取插入ID失败", "error", err)
		return 0, err
	}
	
	// 插入附加表
	_, err = tx.Exec(
		"INSERT INTO "+m.db.TableName("addonarticle")+" (aid, body) VALUES (?, ?)",
		id, download.Body,
	)
	if err != nil {
		logger.Error("插入下载附加表失败", "error", err)
		return 0, err
	}
	
	// 插入下载表
	_, err = tx.Exec(
		"INSERT INTO "+m.db.TableName("addondownload")+" (aid, softname, softversion, softlanguage, softtype, softsize, softos, softdeveloper, softlicense, softscore, softurl, softmirrurl, downcount) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		id, download.SoftName, download.SoftVersion, download.SoftLanguage, download.SoftType, download.SoftSize, download.SoftOS, download.SoftDeveloper, download.SoftLicense, download.SoftScore, download.SoftURL, download.SoftMirrURL, download.DownCount,
	)
	if err != nil {
		logger.Error("插入下载表失败", "error", err)
		return 0, err
	}
	
	// 插入下载截图
	for _, screenshot := range download.Screenshots {
		_, err = tx.Exec(
			"INSERT INTO "+m.db.TableName("download_screenshots")+" (aid, screenshot) VALUES (?, ?)",
			id, screenshot,
		)
		if err != nil {
			logger.Error("插入下载截图失败", "error", err)
			return 0, err
		}
	}
	
	// 提交事务
	if err := tx.Commit(); err != nil {
		logger.Error("提交事务失败", "error", err)
		return 0, err
	}
	
	return id, nil
}

// Update 更新下载
func (m *DownloadModel) Update(download *Download) error {
	// 开始事务
	tx, err := m.db.Begin()
	if err != nil {
		logger.Error("开始事务失败", "error", err)
		return err
	}
	defer tx.Rollback()
	
	// 更新主表
	_, err = tx.Exec(
		"UPDATE "+m.db.TableName("archives")+" SET typeid=?, title=?, shorttitle=?, color=?, writer=?, source=?, litpic=?, pubdate=?, senddate=?, keywords=?, description=?, filename=?, istop=?, isrecommend=?, ishot=?, arcrank=?, click=? WHERE id=?",
		download.TypeID, download.Title, download.ShortTitle, download.Color, download.Writer, download.Source, download.LitPic, download.PubDate, download.SendDate, download.Keywords, download.Description, download.Filename, download.IsTop, download.IsRecommend, download.IsHot, download.ArcRank, download.Click, download.ID,
	)
	if err != nil {
		logger.Error("更新下载主表失败", "error", err)
		return err
	}
	
	// 更新附加表
	_, err = tx.Exec(
		"UPDATE "+m.db.TableName("addonarticle")+" SET body=? WHERE aid=?",
		download.Body, download.ID,
	)
	if err != nil {
		logger.Error("更新下载附加表失败", "error", err)
		return err
	}
	
	// 更新下载表
	_, err = tx.Exec(
		"UPDATE "+m.db.TableName("addondownload")+" SET softname=?, softversion=?, softlanguage=?, softtype=?, softsize=?, softos=?, softdeveloper=?, softlicense=?, softscore=?, softurl=?, softmirrurl=?, downcount=? WHERE aid=?",
		download.SoftName, download.SoftVersion, download.SoftLanguage, download.SoftType, download.SoftSize, download.SoftOS, download.SoftDeveloper, download.SoftLicense, download.SoftScore, download.SoftURL, download.SoftMirrURL, download.DownCount, download.ID,
	)
	if err != nil {
		logger.Error("更新下载表失败", "error", err)
		return err
	}
	
	// 删除旧的下载截图
	_, err = tx.Exec(
		"DELETE FROM "+m.db.TableName("download_screenshots")+" WHERE aid=?",
		download.ID,
	)
	if err != nil {
		logger.Error("删除下载截图失败", "error", err)
		return err
	}
	
	// 插入新的下载截图
	for _, screenshot := range download.Screenshots {
		_, err = tx.Exec(
			"INSERT INTO "+m.db.TableName("download_screenshots")+" (aid, screenshot) VALUES (?, ?)",
			download.ID, screenshot,
		)
		if err != nil {
			logger.Error("插入下载截图失败", "error", err)
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

// Delete 删除下载
func (m *DownloadModel) Delete(id int64) error {
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
		logger.Error("删除下载主表失败", "error", err)
		return err
	}
	
	// 删除附加表
	_, err = tx.Exec(
		"DELETE FROM "+m.db.TableName("addonarticle")+" WHERE aid=?",
		id,
	)
	if err != nil {
		logger.Error("删除下载附加表失败", "error", err)
		return err
	}
	
	// 删除下载表
	_, err = tx.Exec(
		"DELETE FROM "+m.db.TableName("addondownload")+" WHERE aid=?",
		id,
	)
	if err != nil {
		logger.Error("删除下载表失败", "error", err)
		return err
	}
	
	// 删除下载截图
	_, err = tx.Exec(
		"DELETE FROM "+m.db.TableName("download_screenshots")+" WHERE aid=?",
		id,
	)
	if err != nil {
		logger.Error("删除下载截图失败", "error", err)
		return err
	}
	
	// 提交事务
	if err := tx.Commit(); err != nil {
		logger.Error("提交事务失败", "error", err)
		return err
	}
	
	return nil
}

// IncrementDownCount 增加下载次数
func (m *DownloadModel) IncrementDownCount(id int64) error {
	// 执行更新
	_, err := m.db.Execute("UPDATE "+m.db.TableName("addondownload")+" SET downcount = downcount + 1 WHERE aid = ?", id)
	if err != nil {
		logger.Error("更新下载次数失败", "id", id, "error", err)
		return err
	}
	
	return nil
}

// 获取下载截图
func (m *DownloadModel) getDownloadScreenshots(id int64) ([]string, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "download_screenshots")
	qb.Select("screenshot")
	qb.Where("aid = ?", id)
	qb.OrderBy("id ASC")
	
	// 执行查询
	results, err := qb.Get()
	if err != nil {
		return nil, err
	}
	
	// 提取截图URL
	screenshots := make([]string, 0, len(results))
	for _, result := range results {
		if screenshot, ok := result["screenshot"].(string); ok {
			screenshots = append(screenshots, screenshot)
		}
	}
	
	return screenshots, nil
}
