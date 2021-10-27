package common

import (
	"encoding/json"
	"github.com/kataras/iris/v12/sessions"
	"time"
)

const (
	USER_ID = "UserID"
	USER    = "UserInfo"

	EXEC_STATUS_FINISH  = 0
	EXEC_STATUS_RUNNING = 1

	STORAGE_MODE_INSERT = 1
	STORAGE_MODE_UPDATE = 2
	STORAGE_MODE_APPEND = 3
)
const (
	TARGET_TYPE_CRAWLER = 1
	TARGET_TYPE_CLEAN   = 2
)

//日志级别
const (
	LOG_LEVEL_ALL = iota
	LOG_LEVEL_INFO
	LOG_LEVEL_DEBUG
	LOG_LEVEL_WARN
	LOG_LEVEL_ERROR

	LOG_ERROR   = "[错误]"
	LOG_INFO    = "[信息]"
	LOG_DEBUG   = "[调试]"
	LOG_WARNING = "[警告]"
)
const (
	//请求方式  0调试  1运行
	SCHEDULE_METHOD_DEBUG = iota
	SCHEDULE_METHOD_EXECUTE
)

const (
	STATUS_NORMAL = iota
	STATUS_RUNNING

	METHOD_INSERT = 1
	METHOD_UPDATE = 2
	METHOD_APPEND = 3

	STORAGE_NORMAL = 0
	STORAGE_SERVER = 1
	STORAGE_PAN    = 2
)

const (
	SELECTORTYPE_XPATH = iota
	SELECTORTYPE_JSONPATH
	SELECTORTYPE_REGEX
)

//日志类型
const (
	LOG_TYPE_ALL = iota
	LOG_TYPE_SYSTEM
	LOG_TYPE_USER
	LOG_TYPE_DATA
	LOG_TYPE_URL
	LOG_TYPE_FINISH
	LOG_TYPE_HANDLEND
	LOG_TYPE_HANDLEND_URL

	//设置前缀
	PREFIX_CRAWL_LOG  = "CollectorLog"
	PREFIX_CRAWL_DATA = "CollectorData"

	//设置前缀
	PREFIX_CLEAN_LOG  = "CleanLog"
	PREFIX_CLEAN_DATA = "CleanData"

	METHOD_DEBUG  = 0
	METHOD_EXCUTE = 1
)

var (
	Session = sessions.New(
		sessions.Config{
			Cookie:  "SHSESSION",
			Expires: 24 * time.Hour,
		})
)

//日志输出
type LogLevel struct {
	Level     int    `json:"level"`
	Type      int    `json:"type"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
}

//通信数据
type Communication struct {
	Method     int         `json:"method"`
	DebugId    int64       `json:"debug_id"`
	UserId     int64       `json:"user_id"`
	AppId      int64       `json:"app_id"`
	Abort      bool        `json:"abort"`
	Auto       bool        `json:"auto"`
	CrawlField interface{} `json:"crawl_field"`
	Token      string      `json:"token"`
	Content    string      `json:"content"`
}

type FieldData struct {
	Alias string      `json:"alias"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}

type TableHead struct {
	Name  string `json:"name"`
	Alias string `json:"alias"`
	Type  string `json:"type"`
}

//格式化输出日志
func FmtLog(title, content string, level, tp int) []byte {
	l := &LogLevel{
		Level:     level,
		Type:      tp,
		Title:     title,
		Content:   content,
		CreatedAt: time.Now().Unix(),
	}
	res, _ := json.Marshal(l)
	return res
}
