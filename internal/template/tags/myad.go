package tags

import (
	"fmt"
	"strconv"
	"time"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// MyAdTag 广告标签处理器
type MyAdTag struct {
	DB *database.DB
}

// Handle 处理标签
func (t *MyAdTag) Handle(attrs map[string]string, content string, data interface{}) (string, error) {
	// 解析属性
	id := 0
	if idStr, ok := attrs["id"]; ok {
		if i, err := strconv.Atoi(idStr); err == nil && i > 0 {
			id = i
		}
	}

	if id == 0 {
		return "", fmt.Errorf("广告标签缺少id属性")
	}

	// 获取广告信息
	ad, err := t.getAd(id)
	if err != nil {
		return "", err
	}

	// 检查广告是否有效
	if !t.isAdValid(ad) {
		return "", nil
	}

	// 增加广告点击量
	go t.incrementAdHits(id)

	// 返回广告内容
	adBody, _ := ad["normbody"].(string)
	return adBody, nil
}

// 获取广告信息
func (t *MyAdTag) getAd(id int) (map[string]interface{}, error) {
	// 构建查询
	qb := database.NewQueryBuilder(t.DB, "myad")
	qb.Select("*")
	qb.Where("aid = ?", id)

	// 执行查询
	ad, err := qb.First()
	if err != nil {
		logger.Error("查询广告失败", "id", id, "error", err)
		return nil, err
	}

	if ad == nil {
		return nil, fmt.Errorf("广告不存在")
	}

	return ad, nil
}

// 检查广告是否有效
func (t *MyAdTag) isAdValid(ad map[string]interface{}) bool {
	// 检查广告是否启用
	if isCheck, ok := ad["ischeck"].(int64); ok && isCheck != 1 {
		return false
	}

	// 检查广告时间
	now := time.Now()

	// 检查开始时间
	if startTime, ok := ad["starttime"].(time.Time); ok && startTime.After(now) {
		return false
	}

	// 检查结束时间
	if endTime, ok := ad["endtime"].(time.Time); ok && endTime.Before(now) {
		return false
	}

	return true
}

// 增加广告点击量
func (t *MyAdTag) incrementAdHits(id int) {
	// 执行更新
	_, err := t.DB.Execute("UPDATE "+t.DB.TableName("myad")+" SET hits = hits + 1 WHERE aid = ?", id)
	if err != nil {
		logger.Error("更新广告点击量失败", "id", id, "error", err)
	}
}
