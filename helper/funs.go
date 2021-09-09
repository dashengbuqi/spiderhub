package helper

import (
	"encoding/json"
	"fmt"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"github.com/robertkrimen/otto"
	"net/url"
	"strings"
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

func FmtConsole(argumentList []otto.Value) string {
	output := []string{}
	for _, argument := range argumentList {
		output = append(output, fmt.Sprintf("%v", argument))
	}
	return strings.Join(output, " ")
}

func FmtUrl(urlStr string) string {
	urlStr = strings.Replace(urlStr, "\\", "", -1)
	urlStr = strings.Replace(urlStr, "\"", "", -1)
	urlStr, _ = url.QueryUnescape(urlStr)
	return urlStr
}
