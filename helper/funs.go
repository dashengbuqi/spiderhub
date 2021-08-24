package helper

import (
	"encoding/json"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"time"
)

//格式化输出日志
func FmtLog(title, content string, level, tp int) []byte {
	l := &common.LogLevel{
		Level:     level,
		Type:      tp,
		Title:     title,
		Content:   content,
		CreatedAt: time.Now().Unix(),
	}
	res, _ := json.Marshal(l)
	return res
}
