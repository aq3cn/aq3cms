package database

import (
	"fmt"
	"strings"
)

// QueryBuilder SQL查询构建器
type QueryBuilder struct {
	db         *DB
	table      string
	from       string
	fields     []string
	where      []string
	orderBy    string
	groupBy    string
	limit      string
	join       []string
	args       []interface{}
	whereCount int
}

// NewQueryBuilder 创建新的查询构建器
func NewQueryBuilder(db *DB, table string) *QueryBuilder {
	return &QueryBuilder{
		db:     db,
		table:  table,
		fields: []string{"*"},
	}
}

// Select 设置查询字段
func (qb *QueryBuilder) Select(fields ...string) *QueryBuilder {
	qb.fields = fields
	return qb
}

// From 设置FROM子句
func (qb *QueryBuilder) From(from string) *QueryBuilder {
	qb.from = from
	return qb
}

// Where 添加WHERE条件
func (qb *QueryBuilder) Where(condition string, args ...interface{}) *QueryBuilder {
	if qb.whereCount == 0 {
		qb.where = append(qb.where, "WHERE "+condition)
	} else {
		qb.where = append(qb.where, "AND "+condition)
	}
	qb.args = append(qb.args, args...)
	qb.whereCount++
	return qb
}

// OrWhere 添加OR WHERE条件
func (qb *QueryBuilder) OrWhere(condition string, args ...interface{}) *QueryBuilder {
	if qb.whereCount == 0 {
		qb.where = append(qb.where, "WHERE "+condition)
	} else {
		qb.where = append(qb.where, "OR "+condition)
	}
	qb.args = append(qb.args, args...)
	qb.whereCount++
	return qb
}

// OrderBy 设置排序
func (qb *QueryBuilder) OrderBy(orderBy string) *QueryBuilder {
	qb.orderBy = "ORDER BY " + orderBy
	return qb
}

// GroupBy 设置分组
func (qb *QueryBuilder) GroupBy(groupBy string) *QueryBuilder {
	qb.groupBy = "GROUP BY " + groupBy
	return qb
}

// Limit 设置限制
func (qb *QueryBuilder) Limit(limit int, offset ...int) *QueryBuilder {
	if len(offset) > 0 {
		qb.limit = fmt.Sprintf("LIMIT %d, %d", offset[0], limit)
	} else {
		qb.limit = fmt.Sprintf("LIMIT %d", limit)
	}
	return qb
}

// Offset 设置偏移量
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	if qb.limit == "" {
		qb.limit = fmt.Sprintf("LIMIT 18446744073709551615 OFFSET %d", offset)
	} else {
		// 如果已经有LIMIT，则替换
		parts := strings.Split(qb.limit, "LIMIT ")
		if len(parts) > 1 {
			limitParts := strings.Split(parts[1], ",")
			if len(limitParts) > 1 {
				// 已经有OFFSET
				qb.limit = fmt.Sprintf("LIMIT %s OFFSET %d", limitParts[0], offset)
			} else {
				// 没有OFFSET
				qb.limit = fmt.Sprintf("LIMIT %s OFFSET %d", parts[1], offset)
			}
		}
	}
	return qb
}

// Join 添加JOIN
func (qb *QueryBuilder) Join(table, condition string) *QueryBuilder {
	qb.join = append(qb.join, fmt.Sprintf("JOIN %s ON %s", table, condition))
	return qb
}

// LeftJoin 添加LEFT JOIN
func (qb *QueryBuilder) LeftJoin(table, condition string) *QueryBuilder {
	qb.join = append(qb.join, fmt.Sprintf("LEFT JOIN %s ON %s", table, condition))
	return qb
}

// RightJoin 添加RIGHT JOIN
func (qb *QueryBuilder) RightJoin(table, condition string) *QueryBuilder {
	qb.join = append(qb.join, fmt.Sprintf("RIGHT JOIN %s ON %s", table, condition))
	return qb
}

// Get 执行查询并返回结果
func (qb *QueryBuilder) Get() ([]map[string]interface{}, error) {
	query := qb.buildQuery()
	return qb.db.Query(query, qb.args...)
}

// First 获取第一条记录
func (qb *QueryBuilder) First() (map[string]interface{}, error) {
	qb.Limit(1)
	query := qb.buildQuery()
	return qb.db.GetOne(query, qb.args...)
}

// Count 获取记录数
func (qb *QueryBuilder) Count() (int, error) {
	oldFields := qb.fields
	qb.fields = []string{"COUNT(*) as count"}

	query := qb.buildQuery()
	result, err := qb.db.GetOne(query, qb.args...)

	// 恢复原始字段
	qb.fields = oldFields

	if err != nil {
		return 0, err
	}

	if result == nil {
		return 0, nil
	}

	count, ok := result["count"]
	if !ok {
		return 0, fmt.Errorf("count字段不存在")
	}

	switch v := count.(type) {
	case int64:
		return int(v), nil
	case int:
		return v, nil
	case string:
		var countInt int
		_, err := fmt.Sscanf(v, "%d", &countInt)
		if err != nil {
			return 0, err
		}
		return countInt, nil
	default:
		return 0, fmt.Errorf("无法转换count字段类型")
	}
}

// Insert 插入记录
func (qb *QueryBuilder) Insert(data map[string]interface{}) (int64, error) {
	var columns []string
	var placeholders []string
	var values []interface{}

	for column, value := range data {
		// 用反引号包围字段名，防止保留字冲突
		columns = append(columns, "`"+column+"`")
		placeholders = append(placeholders, "?")
		values = append(values, value)
	}

	var tableName string
	if qb.from != "" {
		tableName = qb.from
	} else {
		tableName = qb.db.TableName(qb.table)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	result, err := qb.db.Execute(query, values...)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// Update 更新记录
func (qb *QueryBuilder) Update(data map[string]interface{}) (int64, error) {
	var sets []string
	var values []interface{}

	for column, value := range data {
		// 用反引号包围字段名，防止保留字冲突
		sets = append(sets, "`"+column+"` = ?")
		values = append(values, value)
	}

	// 添加where条件的参数
	values = append(values, qb.args...)

	var tableName string
	if qb.from != "" {
		tableName = qb.from
	} else {
		tableName = qb.db.TableName(qb.table)
	}

	query := fmt.Sprintf("UPDATE %s SET %s %s",
		tableName,
		strings.Join(sets, ", "),
		strings.Join(qb.where, " "))

	result, err := qb.db.Execute(query, values...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// Delete 删除记录
func (qb *QueryBuilder) Delete() (int64, error) {
	var tableName string
	if qb.from != "" {
		tableName = qb.from
	} else {
		tableName = qb.db.TableName(qb.table)
	}

	query := fmt.Sprintf("DELETE FROM %s %s",
		tableName,
		strings.Join(qb.where, " "))

	result, err := qb.db.Execute(query, qb.args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// 构建查询SQL
func (qb *QueryBuilder) buildQuery() string {
	var query string

	if qb.from != "" {
		query = fmt.Sprintf("SELECT %s FROM %s",
			strings.Join(qb.fields, ", "),
			qb.from)
	} else {
		query = fmt.Sprintf("SELECT %s FROM %s",
			strings.Join(qb.fields, ", "),
			qb.db.TableName(qb.table))
	}

	if len(qb.join) > 0 {
		query += " " + strings.Join(qb.join, " ")
	}

	if len(qb.where) > 0 {
		query += " " + strings.Join(qb.where, " ")
	}

	if qb.groupBy != "" {
		query += " " + qb.groupBy
	}

	if qb.orderBy != "" {
		query += " " + qb.orderBy
	}

	if qb.limit != "" {
		query += " " + qb.limit
	}

	return query
}
