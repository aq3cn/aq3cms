package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
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

	// 读取SQL文件
	sqlFile := "sql/create_feedback_table.sql"
	sqlContent, err := ioutil.ReadFile(sqlFile)
	if err != nil {
		log.Fatal("读取SQL文件失败:", err)
	}

	// 分割SQL语句（按分号分割）
	sqlStatements := strings.Split(string(sqlContent), ";")

	fmt.Println("开始执行SQL脚本...")

	// 执行每个SQL语句
	for i, statement := range sqlStatements {
		statement = strings.TrimSpace(statement)
		if statement == "" || strings.HasPrefix(statement, "--") {
			continue
		}

		fmt.Printf("执行SQL语句 %d: %s...\n", i+1, statement[:min(50, len(statement))])

		_, err := db.Exec(statement)
		if err != nil {
			log.Printf("执行SQL语句失败: %v\n语句: %s", err, statement)
			continue
		}
	}

	fmt.Println("SQL脚本执行完成!")

	// 验证表是否创建成功
	fmt.Println("\n验证表结构...")
	rows, err := db.Query("DESCRIBE aq3cms_feedback")
	if err != nil {
		log.Printf("查询表结构失败: %v", err)
		return
	}
	defer rows.Close()

	fmt.Printf("%-15s %-20s %-10s %-10s %-15s %-10s\n", "Field", "Type", "Null", "Key", "Default", "Extra")
	fmt.Println(strings.Repeat("-", 90))

	for rows.Next() {
		var field, fieldType, null, key, defaultVal, extra sql.NullString
		err := rows.Scan(&field, &fieldType, &null, &key, &defaultVal, &extra)
		if err != nil {
			log.Printf("扫描行失败: %v", err)
			continue
		}

		fmt.Printf("%-15s %-20s %-10s %-10s %-15s %-10s\n",
			field.String,
			fieldType.String,
			null.String,
			key.String,
			defaultVal.String,
			extra.String)
	}

	// 检查数据
	fmt.Println("\n检查评论数据...")
	var totalCount int
	err = db.QueryRow("SELECT COUNT(*) FROM aq3cms_feedback").Scan(&totalCount)
	if err != nil {
		log.Printf("查询评论总数失败: %v", err)
		return
	}

	fmt.Printf("评论总数: %d\n", totalCount)

	// 按状态统计
	statusRows, err := db.Query("SELECT ischeck, COUNT(*) as count FROM aq3cms_feedback GROUP BY ischeck")
	if err != nil {
		log.Printf("查询状态统计失败: %v", err)
		return
	}
	defer statusRows.Close()

	fmt.Println("\n按状态统计:")
	for statusRows.Next() {
		var status, count int
		err := statusRows.Scan(&status, &count)
		if err != nil {
			log.Printf("扫描状态行失败: %v", err)
			continue
		}

		statusText := "未知"
		switch status {
		case 0:
			statusText = "待审核"
		case 1:
			statusText = "已审核"
		case -1:
			statusText = "已拒绝"
		}

		fmt.Printf("  %s: %d条\n", statusText, count)
	}

	fmt.Println("\n数据库初始化完成!")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
