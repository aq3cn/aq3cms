package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"aq3cms/config"
	"aq3cms/pkg/logger"

	_ "github.com/go-sql-driver/mysql"
)

// DB 数据库连接实例
type DB struct {
	*sql.DB
	Prefix string
}

// NewConnection 创建新的数据库连接
func NewConnection(cfg config.DatabaseConfig) (*DB, error) {
	var dsn string

	switch cfg.Type {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.Charset)
	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s", cfg.Type)
	}

	db, err := sql.Open(cfg.Type, dsn)
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	db.SetMaxIdleConns(cfg.MaxIdle)
	db.SetMaxOpenConns(cfg.MaxOpen)
	db.SetConnMaxLifetime(time.Hour)

	// 测试连接
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	logger.Info("数据库连接成功", "type", cfg.Type, "host", cfg.Host)

	return &DB{
		DB:     db,
		Prefix: cfg.Prefix,
	}, nil
}

// TableName 获取带前缀的表名
func (db *DB) TableName(name string) string {
	return db.Prefix + name
}

// Query 执行查询并返回结果
func (db *DB) Query(query string, args ...interface{}) ([]map[string]interface{}, error) {
	// 替换查询中的表前缀占位符 #@__
	query = replacePrefix(query, db.Prefix)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)

		entry := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]

			// 处理不同类型的值
			switch v := val.(type) {
			case []byte:
				entry[col] = string(v)
			default:
				entry[col] = v
			}
		}
		tableData = append(tableData, entry)
	}

	return tableData, nil
}

// GetOne 获取单条记录
func (db *DB) GetOne(query string, args ...interface{}) (map[string]interface{}, error) {
	results, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	return results[0], nil
}

// Execute 执行非查询SQL
func (db *DB) Execute(query string, args ...interface{}) (sql.Result, error) {
	query = replacePrefix(query, db.Prefix)
	return db.DB.Exec(query, args...)
}

// 替换SQL中的表前缀占位符
func replacePrefix(query, prefix string) string {
	// 替换 #@__ 为实际的表前缀
	return strings.Replace(query, "#@__", prefix, -1)
}

// Backup 备份数据库
func (db *DB) Backup() (string, error) {
	// 生成备份文件名
	backupFile := fmt.Sprintf("backup_%s.sql", time.Now().Format("20060102150405"))

	// 获取所有表
	_, err := db.Query("SHOW TABLES")
	if err != nil {
		return "", err
	}

	// 简化实现，实际应该使用 mysqldump 或类似工具
	logger.Info("数据库备份", "file", backupFile)

	return backupFile, nil
}

// Restore 恢复数据库
func (db *DB) Restore(file interface{}) error {
	// 简化实现，实际应该解析 SQL 文件并执行
	logger.Info("数据库恢复")

	return nil
}
