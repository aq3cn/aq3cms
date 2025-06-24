package model

import (
	"fmt"
	"time"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// Product 产品模型
type Product struct {
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
	// 产品特有字段
	ProductName  string   `json:"productname"`
	ProductSN    string   `json:"productsn"`
	Price        float64  `json:"price"`
	OldPrice     float64  `json:"oldprice"`
	Units        string   `json:"units"`
	Weight       float64  `json:"weight"`
	Specification string  `json:"specification"`
	Features     string   `json:"features"`
	Parameters   string   `json:"parameters"`
	Images       []string `json:"images"`
	Stock        int      `json:"stock"`
}

// ProductModel 产品模型操作
type ProductModel struct {
	db *database.DB
}

// NewProductModel 创建产品模型
func NewProductModel(db *database.DB) *ProductModel {
	return &ProductModel{
		db: db,
	}
}

// GetByID 根据ID获取产品
func (m *ProductModel) GetByID(id int64) (*Product, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "archives")
	qb.Select("a.*", "t.typename", "t.typedir", "ad.body", "p.*")
	qb.LeftJoin(m.db.TableName("arctype")+" AS t", "a.typeid = t.id")
	qb.LeftJoin(m.db.TableName("addonarticle")+" AS ad", "a.id = ad.aid")
	qb.LeftJoin(m.db.TableName("addonproduct")+" AS p", "a.id = p.aid")
	qb.Where("a.id = ?", id)
	qb.Where("a.channel = 2") // 产品频道ID
	
	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("查询产品失败", "id", id, "error", err)
		return nil, err
	}
	
	if result == nil {
		return nil, fmt.Errorf("产品不存在")
	}
	
	// 转换为产品对象
	product := &Product{}
	product.ID, _ = result["id"].(int64)
	product.TypeID, _ = result["typeid"].(int64)
	product.Title, _ = result["title"].(string)
	product.ShortTitle, _ = result["shorttitle"].(string)
	product.Color, _ = result["color"].(string)
	product.Writer, _ = result["writer"].(string)
	product.Source, _ = result["source"].(string)
	product.LitPic, _ = result["litpic"].(string)
	
	// 处理日期
	if pubdate, ok := result["pubdate"].(time.Time); ok {
		product.PubDate = pubdate
	}
	if senddate, ok := result["senddate"].(time.Time); ok {
		product.SendDate = senddate
	}
	
	product.Keywords, _ = result["keywords"].(string)
	product.Description, _ = result["description"].(string)
	product.Filename, _ = result["filename"].(string)
	
	// 处理整数字段
	if istop, ok := result["istop"].(int64); ok {
		product.IsTop = int(istop)
	}
	if isrecommend, ok := result["isrecommend"].(int64); ok {
		product.IsRecommend = int(isrecommend)
	}
	if ishot, ok := result["ishot"].(int64); ok {
		product.IsHot = int(ishot)
	}
	if arcrank, ok := result["arcrank"].(int64); ok {
		product.ArcRank = int(arcrank)
	}
	if click, ok := result["click"].(int64); ok {
		product.Click = int(click)
	}
	
	product.Body, _ = result["body"].(string)
	product.TypeName, _ = result["typename"].(string)
	product.TypeDir, _ = result["typedir"].(string)
	
	// 处理产品特有字段
	product.ProductName, _ = result["productname"].(string)
	product.ProductSN, _ = result["productsn"].(string)
	
	// 处理价格
	if price, ok := result["price"].(float64); ok {
		product.Price = price
	} else if priceStr, ok := result["price"].(string); ok {
		fmt.Sscanf(priceStr, "%f", &product.Price)
	}
	
	if oldprice, ok := result["oldprice"].(float64); ok {
		product.OldPrice = oldprice
	} else if oldpriceStr, ok := result["oldprice"].(string); ok {
		fmt.Sscanf(oldpriceStr, "%f", &product.OldPrice)
	}
	
	product.Units, _ = result["units"].(string)
	
	// 处理重量
	if weight, ok := result["weight"].(float64); ok {
		product.Weight = weight
	} else if weightStr, ok := result["weight"].(string); ok {
		fmt.Sscanf(weightStr, "%f", &product.Weight)
	}
	
	product.Specification, _ = result["specification"].(string)
	product.Features, _ = result["features"].(string)
	product.Parameters, _ = result["parameters"].(string)
	
	// 处理库存
	if stock, ok := result["stock"].(int64); ok {
		product.Stock = int(stock)
	} else if stockStr, ok := result["stock"].(string); ok {
		fmt.Sscanf(stockStr, "%d", &product.Stock)
	}
	
	// 获取产品图片
	images, err := m.getProductImages(id)
	if err != nil {
		logger.Error("获取产品图片失败", "id", id, "error", err)
	}
	product.Images = images
	
	return product, nil
}

// GetList 获取产品列表
func (m *ProductModel) GetList(typeid int64, page, pageSize int) ([]*Product, int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "archives")
	qb.Select("a.*", "t.typename", "t.typedir", "p.productname", "p.productsn", "p.price", "p.oldprice")
	qb.LeftJoin(m.db.TableName("arctype")+" AS t", "a.typeid = t.id")
	qb.LeftJoin(m.db.TableName("addonproduct")+" AS p", "a.id = p.aid")
	
	// 添加条件
	qb.Where("a.arcrank > -1")
	qb.Where("a.channel = 2") // 产品频道ID
	if typeid > 0 {
		qb.Where("a.typeid = ?", typeid)
	}
	
	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("查询产品总数失败", "error", err)
		return nil, 0, err
	}
	
	// 设置分页
	offset := (page - 1) * pageSize
	qb.OrderBy("a.pubdate DESC")
	qb.Limit(pageSize, offset)
	
	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("查询产品列表失败", "error", err)
		return nil, 0, err
	}
	
	// 转换为产品对象
	products := make([]*Product, 0, len(results))
	for _, result := range results {
		product := &Product{}
		product.ID, _ = result["id"].(int64)
		product.TypeID, _ = result["typeid"].(int64)
		product.Title, _ = result["title"].(string)
		product.ShortTitle, _ = result["shorttitle"].(string)
		product.Color, _ = result["color"].(string)
		product.Writer, _ = result["writer"].(string)
		product.Source, _ = result["source"].(string)
		product.LitPic, _ = result["litpic"].(string)
		
		// 处理日期
		if pubdate, ok := result["pubdate"].(time.Time); ok {
			product.PubDate = pubdate
		}
		if senddate, ok := result["senddate"].(time.Time); ok {
			product.SendDate = senddate
		}
		
		product.Keywords, _ = result["keywords"].(string)
		product.Description, _ = result["description"].(string)
		product.Filename, _ = result["filename"].(string)
		
		// 处理整数字段
		if istop, ok := result["istop"].(int64); ok {
			product.IsTop = int(istop)
		}
		if isrecommend, ok := result["isrecommend"].(int64); ok {
			product.IsRecommend = int(isrecommend)
		}
		if ishot, ok := result["ishot"].(int64); ok {
			product.IsHot = int(ishot)
		}
		if arcrank, ok := result["arcrank"].(int64); ok {
			product.ArcRank = int(arcrank)
		}
		if click, ok := result["click"].(int64); ok {
			product.Click = int(click)
		}
		
		product.TypeName, _ = result["typename"].(string)
		product.TypeDir, _ = result["typedir"].(string)
		
		// 处理产品特有字段
		product.ProductName, _ = result["productname"].(string)
		product.ProductSN, _ = result["productsn"].(string)
		
		// 处理价格
		if price, ok := result["price"].(float64); ok {
			product.Price = price
		} else if priceStr, ok := result["price"].(string); ok {
			fmt.Sscanf(priceStr, "%f", &product.Price)
		}
		
		if oldprice, ok := result["oldprice"].(float64); ok {
			product.OldPrice = oldprice
		} else if oldpriceStr, ok := result["oldprice"].(string); ok {
			fmt.Sscanf(oldpriceStr, "%f", &product.OldPrice)
		}
		
		products = append(products, product)
	}
	
	return products, total, nil
}

// Create 创建产品
func (m *ProductModel) Create(product *Product) (int64, error) {
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
		product.TypeID, 2, product.Title, product.ShortTitle, product.Color, product.Writer, product.Source, product.LitPic, product.PubDate, product.SendDate, product.Keywords, product.Description, product.Filename, product.IsTop, product.IsRecommend, product.IsHot, product.ArcRank, product.Click,
	)
	if err != nil {
		logger.Error("插入产品主表失败", "error", err)
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
		id, product.Body,
	)
	if err != nil {
		logger.Error("插入产品附加表失败", "error", err)
		return 0, err
	}
	
	// 插入产品表
	_, err = tx.Exec(
		"INSERT INTO "+m.db.TableName("addonproduct")+" (aid, productname, productsn, price, oldprice, units, weight, specification, features, parameters, stock) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		id, product.ProductName, product.ProductSN, product.Price, product.OldPrice, product.Units, product.Weight, product.Specification, product.Features, product.Parameters, product.Stock,
	)
	if err != nil {
		logger.Error("插入产品表失败", "error", err)
		return 0, err
	}
	
	// 插入产品图片
	for _, image := range product.Images {
		_, err = tx.Exec(
			"INSERT INTO "+m.db.TableName("product_images")+" (aid, image) VALUES (?, ?)",
			id, image,
		)
		if err != nil {
			logger.Error("插入产品图片失败", "error", err)
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

// Update 更新产品
func (m *ProductModel) Update(product *Product) error {
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
		product.TypeID, product.Title, product.ShortTitle, product.Color, product.Writer, product.Source, product.LitPic, product.PubDate, product.SendDate, product.Keywords, product.Description, product.Filename, product.IsTop, product.IsRecommend, product.IsHot, product.ArcRank, product.Click, product.ID,
	)
	if err != nil {
		logger.Error("更新产品主表失败", "error", err)
		return err
	}
	
	// 更新附加表
	_, err = tx.Exec(
		"UPDATE "+m.db.TableName("addonarticle")+" SET body=? WHERE aid=?",
		product.Body, product.ID,
	)
	if err != nil {
		logger.Error("更新产品附加表失败", "error", err)
		return err
	}
	
	// 更新产品表
	_, err = tx.Exec(
		"UPDATE "+m.db.TableName("addonproduct")+" SET productname=?, productsn=?, price=?, oldprice=?, units=?, weight=?, specification=?, features=?, parameters=?, stock=? WHERE aid=?",
		product.ProductName, product.ProductSN, product.Price, product.OldPrice, product.Units, product.Weight, product.Specification, product.Features, product.Parameters, product.Stock, product.ID,
	)
	if err != nil {
		logger.Error("更新产品表失败", "error", err)
		return err
	}
	
	// 删除旧的产品图片
	_, err = tx.Exec(
		"DELETE FROM "+m.db.TableName("product_images")+" WHERE aid=?",
		product.ID,
	)
	if err != nil {
		logger.Error("删除产品图片失败", "error", err)
		return err
	}
	
	// 插入新的产品图片
	for _, image := range product.Images {
		_, err = tx.Exec(
			"INSERT INTO "+m.db.TableName("product_images")+" (aid, image) VALUES (?, ?)",
			product.ID, image,
		)
		if err != nil {
			logger.Error("插入产品图片失败", "error", err)
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

// Delete 删除产品
func (m *ProductModel) Delete(id int64) error {
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
		logger.Error("删除产品主表失败", "error", err)
		return err
	}
	
	// 删除附加表
	_, err = tx.Exec(
		"DELETE FROM "+m.db.TableName("addonarticle")+" WHERE aid=?",
		id,
	)
	if err != nil {
		logger.Error("删除产品附加表失败", "error", err)
		return err
	}
	
	// 删除产品表
	_, err = tx.Exec(
		"DELETE FROM "+m.db.TableName("addonproduct")+" WHERE aid=?",
		id,
	)
	if err != nil {
		logger.Error("删除产品表失败", "error", err)
		return err
	}
	
	// 删除产品图片
	_, err = tx.Exec(
		"DELETE FROM "+m.db.TableName("product_images")+" WHERE aid=?",
		id,
	)
	if err != nil {
		logger.Error("删除产品图片失败", "error", err)
		return err
	}
	
	// 提交事务
	if err := tx.Commit(); err != nil {
		logger.Error("提交事务失败", "error", err)
		return err
	}
	
	return nil
}

// 获取产品图片
func (m *ProductModel) getProductImages(id int64) ([]string, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "product_images")
	qb.Select("image")
	qb.Where("aid = ?", id)
	qb.OrderBy("id ASC")
	
	// 执行查询
	results, err := qb.Get()
	if err != nil {
		return nil, err
	}
	
	// 提取图片URL
	images := make([]string, 0, len(results))
	for _, result := range results {
		if image, ok := result["image"].(string); ok {
			images = append(images, image)
		}
	}
	
	return images, nil
}
