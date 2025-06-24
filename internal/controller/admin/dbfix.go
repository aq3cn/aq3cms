package admin

import (
	"encoding/json"
	"net/http"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// DBFixController 数据库修复控制器
type DBFixController struct {
	db *database.DB
}

// NewDBFixController 创建数据库修复控制器
func NewDBFixController(db *database.DB) *DBFixController {
	return &DBFixController{
		db: db,
	}
}

// FixMemberTable 修复会员表结构
func (c *DBFixController) FixMemberTable(w http.ResponseWriter, r *http.Request) {
	// 只允许POST请求
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	logger.Info("开始修复会员表结构")

	// 修复sex字段
	logger.Info("修复sex字段为enum类型")
	_, err := c.db.Exec("ALTER TABLE " + c.db.TableName("member") + " MODIFY COLUMN sex ENUM('男','女','保密') DEFAULT '保密'")
	if err != nil {
		logger.Error("修复sex字段失败", "error", err)
		c.jsonResponse(w, map[string]interface{}{
			"success": false,
			"message": "修复sex字段失败: " + err.Error(),
		})
		return
	}
	logger.Info("sex字段修复成功")

	// 修复pwd字段长度
	logger.Info("修复pwd字段长度")
	_, err = c.db.Exec("ALTER TABLE " + c.db.TableName("member") + " MODIFY COLUMN pwd VARCHAR(255) DEFAULT ''")
	if err != nil {
		logger.Error("修复pwd字段失败", "error", err)
		c.jsonResponse(w, map[string]interface{}{
			"success": false,
			"message": "修复pwd字段失败: " + err.Error(),
		})
		return
	}
	logger.Info("pwd字段修复成功")

	// 检查修复结果
	sqlRows, err := c.db.DB.Query(`
		SELECT COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, COLUMN_DEFAULT
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = ?
		  AND COLUMN_NAME IN ('sex', 'pwd')
	`, c.db.TableName("member"))
	if err != nil {
		logger.Error("查询字段信息失败", "error", err)
		c.jsonResponse(w, map[string]interface{}{
			"success": false,
			"message": "查询字段信息失败: " + err.Error(),
		})
		return
	}
	defer sqlRows.Close()

	var fields []map[string]interface{}
	for sqlRows.Next() {
		var columnName, columnType, isNullable string
		var columnDefault *string

		err := sqlRows.Scan(&columnName, &columnType, &isNullable, &columnDefault)
		if err != nil {
			logger.Error("扫描行失败", "error", err)
			continue
		}

		defaultVal := "NULL"
		if columnDefault != nil {
			defaultVal = *columnDefault
		}

		fields = append(fields, map[string]interface{}{
			"name":     columnName,
			"type":     columnType,
			"nullable": isNullable,
			"default":  defaultVal,
		})
	}

	logger.Info("会员表结构修复完成")

	c.jsonResponse(w, map[string]interface{}{
		"success": true,
		"message": "会员表结构修复成功",
		"fields":  fields,
	})
}

// CreateFeedbackTable 创建评论表
func (c *DBFixController) CreateFeedbackTable(w http.ResponseWriter, r *http.Request) {
	// 只允许POST请求
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	logger.Info("开始创建评论表")

	// 删除表（如果存在）
	logger.Info("删除现有评论表（如果存在）")
	_, err := c.db.Exec("DROP TABLE IF EXISTS " + c.db.TableName("feedback"))
	if err != nil {
		logger.Error("删除评论表失败", "error", err)
		c.jsonResponse(w, map[string]interface{}{
			"success": false,
			"message": "删除评论表失败: " + err.Error(),
		})
		return
	}

	// 创建评论表
	logger.Info("创建评论表")
	createTableSQL := `
		CREATE TABLE ` + c.db.TableName("feedback") + ` (
		  id int(11) NOT NULL AUTO_INCREMENT COMMENT '评论ID',
		  aid int(11) NOT NULL DEFAULT '0' COMMENT '文章ID',
		  typeid int(11) NOT NULL DEFAULT '0' COMMENT '栏目ID',
		  username varchar(50) NOT NULL DEFAULT '' COMMENT '用户名',
		  mid int(11) NOT NULL DEFAULT '0' COMMENT '会员ID',
		  ip varchar(15) NOT NULL DEFAULT '' COMMENT 'IP地址',
		  ischeck tinyint(1) NOT NULL DEFAULT '0' COMMENT '审核状态：0=待审核，1=已审核，-1=已拒绝',
		  dtime datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '评论时间',
		  content text NOT NULL COMMENT '评论内容',
		  parentid int(11) NOT NULL DEFAULT '0' COMMENT '父评论ID（回复功能）',
		  score int(11) NOT NULL DEFAULT '0' COMMENT '评分',
		  goodcount int(11) NOT NULL DEFAULT '0' COMMENT '点赞数',
		  badcount int(11) NOT NULL DEFAULT '0' COMMENT '踩数',
		  userface varchar(255) NOT NULL DEFAULT '' COMMENT '用户头像',
		  channeltype tinyint(1) NOT NULL DEFAULT '1' COMMENT '频道类型：1=文章',
		  PRIMARY KEY (id),
		  KEY idx_aid (aid),
		  KEY idx_typeid (typeid),
		  KEY idx_mid (mid),
		  KEY idx_ischeck (ischeck),
		  KEY idx_dtime (dtime),
		  KEY idx_parentid (parentid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评论反馈表'
	`

	_, err = c.db.Exec(createTableSQL)
	if err != nil {
		logger.Error("创建评论表失败", "error", err)
		c.jsonResponse(w, map[string]interface{}{
			"success": false,
			"message": "创建评论表失败: " + err.Error(),
		})
		return
	}
	logger.Info("评论表创建成功")

	// 插入示例数据
	logger.Info("插入示例评论数据")
	insertSQL := `
		INSERT INTO ` + c.db.TableName("feedback") + ` (aid, typeid, username, mid, ip, ischeck, dtime, content, parentid, score, goodcount, badcount, userface, channeltype) VALUES
		(1, 1, '张三', 1, '192.168.1.100', 1, '2024-01-15 10:30:00', '这篇文章写得很好，内容详实，对我很有帮助！', 0, 5, 12, 1, '/uploads/face/default.jpg', 1),
		(1, 1, '李四', 2, '192.168.1.101', 1, '2024-01-15 11:15:00', '同意楼上的观点，作者的分析很到位。', 1, 4, 8, 0, '/uploads/face/default.jpg', 1),
		(1, 1, '王五', 0, '192.168.1.102', 0, '2024-01-15 14:20:00', '有些地方还可以更深入一些，期待作者的后续文章。', 0, 4, 3, 0, '/uploads/face/default.jpg', 1),
		(2, 1, '赵六', 3, '192.168.1.103', 1, '2024-01-16 09:45:00', '技术文章就应该这样写，清晰明了，实用性强。', 0, 5, 15, 0, '/uploads/face/default.jpg', 1),
		(2, 1, '钱七', 0, '192.168.1.104', 0, '2024-01-16 16:30:00', '代码示例很棒，已经在项目中应用了，效果不错！', 0, 5, 6, 0, '/uploads/face/default.jpg', 1),
		(3, 2, '孙八', 4, '192.168.1.105', 1, '2024-01-17 08:20:00', '这个解决方案确实有效，感谢分享！', 0, 4, 9, 1, '/uploads/face/default.jpg', 1),
		(3, 2, '周九', 0, '192.168.1.106', -1, '2024-01-17 12:10:00', '垃圾文章，浪费时间！', 0, 1, 0, 8, '/uploads/face/default.jpg', 1),
		(1, 1, '吴十', 5, '192.168.1.107', 1, '2024-01-17 15:45:00', '@李四 我也有同样的感受，作者的思路很清晰。', 2, 4, 5, 0, '/uploads/face/default.jpg', 1),
		(4, 3, '郑十一', 0, '192.168.1.108', 0, '2024-01-18 10:15:00', '这个功能正是我需要的，请问有完整的源码吗？', 0, 4, 2, 0, '/uploads/face/default.jpg', 1),
		(4, 3, '王十二', 6, '192.168.1.109', 1, '2024-01-18 14:30:00', '文档写得很详细，按照步骤操作成功了！', 0, 5, 11, 0, '/uploads/face/default.jpg', 1)
	`

	_, err = c.db.Exec(insertSQL)
	if err != nil {
		logger.Error("插入示例数据失败", "error", err)
		c.jsonResponse(w, map[string]interface{}{
			"success": false,
			"message": "插入示例数据失败: " + err.Error(),
		})
		return
	}
	logger.Info("示例数据插入成功")

	// 检查创建结果
	var totalCount int
	err = c.db.DB.QueryRow("SELECT COUNT(*) FROM " + c.db.TableName("feedback")).Scan(&totalCount)
	if err != nil {
		logger.Error("查询评论总数失败", "error", err)
		totalCount = 0
	}

	logger.Info("评论表创建完成", "总评论数", totalCount)

	c.jsonResponse(w, map[string]interface{}{
		"success":     true,
		"message":     "评论表创建成功",
		"total_count": totalCount,
		"table_name":  c.db.TableName("feedback"),
	})
}

// jsonResponse 返回JSON响应
func (c *DBFixController) jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
