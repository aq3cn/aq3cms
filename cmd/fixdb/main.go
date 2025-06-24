package main

import (
	"database/sql"
	"fmt"
	"log"

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

	// 修复sex字段
	fmt.Println("正在修复sex字段...")
	_, err = db.Exec("ALTER TABLE aq3cms_member MODIFY COLUMN sex ENUM('男','女','保密') DEFAULT '保密'")
	if err != nil {
		log.Printf("修复sex字段失败: %v", err)
	} else {
		fmt.Println("sex字段修复成功!")
	}

	// 修复pwd字段
	fmt.Println("正在修复pwd字段...")
	_, err = db.Exec("ALTER TABLE aq3cms_member MODIFY COLUMN pwd VARCHAR(255) DEFAULT ''")
	if err != nil {
		log.Printf("修复pwd字段失败: %v", err)
	} else {
		fmt.Println("pwd字段修复成功!")
	}

	// 检查字段信息
	fmt.Println("\n=== 检查字段信息 ===")
	rows, err := db.Query(`
		SELECT COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, COLUMN_DEFAULT 
		FROM INFORMATION_SCHEMA.COLUMNS 
		WHERE TABLE_SCHEMA = 'aq3cms' 
		  AND TABLE_NAME = 'aq3cms_member' 
		  AND COLUMN_NAME IN ('sex', 'pwd')
	`)
	if err != nil {
		log.Printf("查询字段信息失败: %v", err)
		return
	}
	defer rows.Close()

	fmt.Printf("%-12s %-30s %-12s %-15s\n", "字段名", "字段类型", "允许NULL", "默认值")
	fmt.Println("----------------------------------------------------------------")

	for rows.Next() {
		var columnName, columnType, isNullable string
		var columnDefault sql.NullString

		err := rows.Scan(&columnName, &columnType, &isNullable, &columnDefault)
		if err != nil {
			log.Printf("扫描行失败: %v", err)
			continue
		}

		defaultVal := "NULL"
		if columnDefault.Valid {
			defaultVal = columnDefault.String
		}

		fmt.Printf("%-12s %-30s %-12s %-15s\n", columnName, columnType, isNullable, defaultVal)
	}

	fmt.Println("\n数据库修复完成!")
}
