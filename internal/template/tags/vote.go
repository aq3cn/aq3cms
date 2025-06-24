package tags

import (
	"bytes"
	"fmt"
	"strconv"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// VoteTag 投票标签处理器
type VoteTag struct {
	DB *database.DB
}

// Handle 处理标签
func (t *VoteTag) Handle(attrs map[string]string, content string, data interface{}) (string, error) {
	// 解析属性
	id := 0
	if idStr, ok := attrs["id"]; ok {
		if i, err := strconv.Atoi(idStr); err == nil && i > 0 {
			id = i
		}
	}

	if id == 0 {
		return "", fmt.Errorf("投票标签缺少id属性")
	}

	// 获取投票信息
	vote, err := t.getVote(id)
	if err != nil {
		return "", err
	}

	// 获取投票选项
	options, err := t.getVoteOptions(id)
	if err != nil {
		return "", err
	}

	// 生成投票表单
	var result bytes.Buffer

	// 表单开始
	result.WriteString(fmt.Sprintf("<form name=\"voteform%d\" method=\"post\" action=\"/vote.php\" target=\"_blank\">\n", id))
	result.WriteString(fmt.Sprintf("<input type=\"hidden\" name=\"id\" value=\"%d\">\n", id))

	// 投票标题
	result.WriteString(fmt.Sprintf("<div class=\"votetitle\">%s</div>\n", vote["title"]))

	// 投票选项
	result.WriteString("<div class=\"voteoptions\">\n")

	// 判断是单选还是多选
	inputType := "radio"
	if isMultiple, ok := vote["ismore"].(int64); ok && isMultiple == 1 {
		inputType = "checkbox"
	}

	// 生成选项
	for _, option := range options {
		optionID, _ := option["id"].(int64)
		optionName, _ := option["name"].(string)

		result.WriteString(fmt.Sprintf("<div class=\"voteoption\"><input type=\"%s\" name=\"voteitem\" value=\"%d\"> %s</div>\n",
			inputType, optionID, optionName))
	}

	result.WriteString("</div>\n")

	// 提交按钮
	result.WriteString("<div class=\"votesubmit\"><input type=\"submit\" name=\"votesubmit\" value=\"投票\"></div>\n")

	// 表单结束
	result.WriteString("</form>\n")

	return result.String(), nil
}

// 获取投票信息
func (t *VoteTag) getVote(id int) (map[string]interface{}, error) {
	// 构建查询
	qb := database.NewQueryBuilder(t.DB, "vote")
	qb.Select("*")
	qb.Where("id = ?", id)

	// 执行查询
	vote, err := qb.First()
	if err != nil {
		logger.Error("查询投票失败", "id", id, "error", err)
		return nil, err
	}

	if vote == nil {
		return nil, fmt.Errorf("投票不存在")
	}

	return vote, nil
}

// 获取投票选项
func (t *VoteTag) getVoteOptions(id int) ([]map[string]interface{}, error) {
	// 构建查询
	qb := database.NewQueryBuilder(t.DB, "vote_option")
	qb.Select("*")
	qb.Where("voteid = ?", id)
	qb.OrderBy("sortid ASC")

	// 执行查询
	options, err := qb.Get()
	if err != nil {
		logger.Error("查询投票选项失败", "voteid", id, "error", err)
		return nil, err
	}

	return options, nil
}
