package common

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
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
	PREFIX_CRAWL_LOG  = "crawlLog"
	PREFIX_CRAWL_DATA = "crawlData"

	//设置前缀
	PREFIX_CLEAN_LOG  = "cleanLog"
	PREFIX_CLEAN_DATA = "cleanData"
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
	Method     int                `json:"method"`
	DebugId    int                `json:"debug_id"`
	UserId     int                `json:"user_id"`
	AppId      primitive.ObjectID `json:"app_id"`
	Abort      bool               `json:"abort"`
	Auto       bool               `json:"auto"`
	CrawlField interface{}        `json:"crawl_field"`
	Token      string             `json:"token"`
	Content    string             `json:"content"`
}

type FieldData struct {
	Alias string      `json:"alias"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}
