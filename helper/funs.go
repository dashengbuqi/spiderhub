package helper

import (
	"crypto/md5"
	"encoding/hex"
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

func StringToM5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

var weekDayArr = map[string]string{
	"Monday":    "星期一",
	"Tuesday":   "星期二",
	"Wednesday": "星期三",
	"Thursday":  "星期四",
	"Friday":    "星期五",
	"Saturday":  "星期六",
	"Sunday":    "星期天",
}

//转换成时间，星期，年月日
func FmtDateTime(t int64) string {
	tm := time.Now()
	now := tm.Unix()
	yesterday_start := time.Date(tm.Year(), tm.Month(), tm.Day()-1, 0, 0, 0, 0, tm.Location()).Unix()
	yesterday_end := time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, tm.Location()).Unix()
	diff := now - t

	day := int64(3600 * 24)
	week := day * 7
	if t > yesterday_start && t < yesterday_end {
		hm := time.Unix(t, 0).Format("15:04")
		return "昨天" + "" + hm
	} else {
		if diff <= day {
			res := time.Unix(t, 0).Format("15:04")
			return res
		} else if diff <= week {
			wd := time.Unix(t, 0).Weekday().String()
			hm := time.Unix(t, 0).Format("15:04")
			return weekDayArr[wd] + " " + hm
		} else {
			res := time.Unix(t, 0).Format("2006-01-02 15:04")
			return res
		}
	}
}
