package spiderhub

import (
	"github.com/astaxie/beego/logs"
	"github.com/dashengbuqi/spiderhub/configs"
	"github.com/dashengbuqi/spiderhub/helper"
	"strings"
)

var (
	Logger *logs.BeeLogger

	levelMap = map[string]int{
		"emergency": 0,
		"alert":     1,
		"critical":  2,
		"error":     3,
		"warning":   4,
		"info":      5,
		"debug":     6,
	}
)

func init() {
	//path, _ := configs.GetParamsByField("Log", "Path")
	params, _ := configs.GetParams("Log")
	fullPath := helper.CurDir() + helper.GetSeparator() + params.(map[interface{}]interface{})["Path"].(string)

	var name, n string
	var err error
	n, err = helper.GetBinaryCurrentAppName()
	if err != nil {
		name = params.(map[interface{}]interface{})["Name"].(string)
	} else {
		name = n
	}
	levelStr := params.(map[interface{}]interface{})["Level"].(string)
	levelStr = strings.ToLower(levelStr)
	//# debug | info | warning | error
	filename := strings.Replace(fullPath, "\\", "\\\\", -1) + "/" + name + ".log"
	Logger = logs.NewLogger(10000)
	jsonConfig := `{
	"filename": "` + filename + `",
	"color": true
}`
	Logger.SetLogger("file", jsonConfig) // 设置日志记录方式：本地文件记录
	Logger.SetLevel(levelMap[levelStr])  // 设置日志写入缓冲区的等级
	Logger.EnableFuncCallDepth(true)     // 输出log时能显示输出文件名和行号（非必须）
}
