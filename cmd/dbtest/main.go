package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 连接数据库
	db, err := sql.Open("mysql", "aq3cms:aq3cms123@tcp(localhost:3306)/aq3cms?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatal("数据库连接测试失败:", err)
	}

	fmt.Println("数据库连接成功!")

	// 查看会员表结构
	rows, err := db.Query("DESCRIBE aq3cms_member")
	if err != nil {
		log.Fatal("查询表结构失败:", err)
	}
	defer rows.Close()

	fmt.Println("\n=== aq3cms_member 表结构 ===")
	fmt.Printf("%-15s %-20s %-10s %-10s %-10s %-10s\n", "Field", "Type", "Null", "Key", "Default", "Extra")
	fmt.Println(strings.Repeat("-", 80))

	for rows.Next() {
		var field, fieldType, null, key, defaultVal, extra sql.NullString
		err := rows.Scan(&field, &fieldType, &null, &key, &defaultVal, &extra)
		if err != nil {
			log.Fatal("扫描行失败:", err)
		}

		fmt.Printf("%-15s %-20s %-10s %-10s %-10s %-10s\n",
			field.String,
			fieldType.String,
			null.String,
			key.String,
			defaultVal.String,
			extra.String)
	}

	// 特别检查sex字段
	var sexType string
	err = db.QueryRow("SELECT COLUMN_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = 'aq3cms' AND TABLE_NAME = 'aq3cms_member' AND COLUMN_NAME = 'sex'").Scan(&sexType)
	if err != nil {
		log.Fatal("查询sex字段类型失败:", err)
	}

	fmt.Printf("\n=== sex字段详细信息 ===\n")
	fmt.Printf("字段类型: %s\n", sexType)
}
