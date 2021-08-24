package spider

import (
	"bytes"
	"encoding/json"
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"github.com/dashengbuqi/spiderhub/internal/crawler/rules"
	"github.com/gocolly/colly"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//请求前回调函数
func (this *Spider) onRequest(r *colly.Request) {
	if this.method == common.SCHEDULE_METHOD_DEBUG {
		nowTime := time.Now().Unix()
		if this.tm+600 < nowTime {
			this.log <- helper.FmtLog(common.LOG_INFO, "调试模式只允许执行10分钟", common.LOG_LEVEL_DEBUG, common.LOG_TYPE_SYSTEM)
			this.abort = true
		}
	}
	if this.abort == true {
		this.log <- helper.FmtLog(common.LOG_INFO, "停止爬虫[url] "+r.URL.String(), common.LOG_LEVEL_INFO, common.LOG_TYPE_SYSTEM)
		r.Abort()
	}
	queueCount, _ := this.queue.Size()
	this.log <- helper.FmtLog(common.LOG_INFO, "爬虫剩余"+strconv.Itoa(queueCount), common.LOG_LEVEL_INFO, common.LOG_TYPE_SYSTEM)
	//检查当前请求是否合法
	var allowUrl bool
	if len(this.rules[rules.DOMAIN].([]string)) > 0 {
		for _, regex := range this.rules[rules.DOMAIN].([]string) {
			if match, err := regexp.MatchString(regex, r.URL.Host); err == nil && match == true {
				allowUrl = true
			}
		}
	}
	if allowUrl == false {
		this.log <- helper.FmtLog(common.LOG_ERROR, "不在请求范围内停止执行[url] "+r.URL.String(), common.LOG_LEVEL_ERROR, common.LOG_TYPE_SYSTEM)
		r.Abort()
	}
	if _, ok := this.rules[rules.COOKIE]; ok {
		cookie := this.rules[rules.COOKIE].(string)
		r.Headers.Add("cookie", cookie)
	}
	this.log <- helper.FmtLog(common.LOG_INFO, "加载页面:"+r.URL.String(), common.LOG_LEVEL_INFO, common.LOG_TYPE_SYSTEM)
	//请求前添加header
	if res, err := this.container.Call(rules.FUNC_BEFORE_CRAWL, nil, r); err == nil {
		if res.IsDefined() {
			if res.IsObject() {
				keys := res.Object().Keys()
				for _, key := range keys {
					val, _ := res.Object().Get(key)
					r.Headers.Add(key, val.String())
				}
			}
		}
	}
	//请求前链接需要需要时间戳或其它参数重写此方法更新当前链接即可
	if res, err := this.container.Call(rules.FUNC_BEFORE_DOWNLOAD_PAGE, nil, r.URL.String()); err == nil {
		if res.IsDefined() && res.IsObject() {
			keys := res.Object().Keys()
			params := make(map[string]interface{})
			for _, key := range keys {
				key = strings.ToUpper(key)
				if key == "HEADER" {
					value, _ := res.Object().Get(key)
					if value.IsObject() {
						hkeys := value.Object().Keys()
						hmap := make(map[string]string)
						for _, hk := range hkeys {
							hval, _ := value.Object().Get(hk)
							if hval.IsDefined() {
								hmap[hk] = hval.String()
							}
						}
						params[key] = hmap
					}
				} else if key == "BODY" {
					value, _ := res.Object().Get(key)
					if value.IsObject() {
						bKeys := value.Object().Keys()
						for _, k1 := range bKeys {
							v1, _ := value.Object().Get(k1)
							if v1.IsObject() {
								v1keys := v1.Object().Keys()
								subParams := make(map[string]interface{})
								for _, k2 := range v1keys {
									v2, _ := v1.Object().Get(k2)
									if v2.IsObject() {
										v2Keys := v2.Object().Keys()
										subSubParams := make(map[string]interface{})
										for _, k3 := range v2Keys {
											v3, _ := v2.Object().Get(k3)
											if v3.IsString() {
												subSubParams[k3], _ = v3.Export()
											}
										}
										subParams[k2] = subSubParams
									} else {
										subParams[k2], err = v2.Export()
										if err != nil {
											this.log <- helper.FmtLog(common.LOG_ERROR, err.Error(), common.LOG_LEVEL_ERROR, common.LOG_TYPE_SYSTEM)
										}
									}
								}
								params[k1] = subParams
							} else {
								params[k1], err = v1.Export()
								if err != nil {
									this.log <- helper.FmtLog(common.LOG_ERROR, err.Error(), common.LOG_LEVEL_ERROR, common.LOG_TYPE_SYSTEM)
								}
							}
						}
					} else {
						params[key], _ = value.Export()
					}
				} else {
					value, _ := res.Object().Get(key)
					params[key], _ = value.Export()
				}
			}
			//加载请求
			if len(params) > 0 {
				if _, ok := params["METHOD"]; ok {
					method := strings.ToUpper(params["METHOD"].(string))
					if method == http.MethodGet {
						if _, okey := params["URL"]; okey && len(params["URL"].(string)) > 0 {
							uri := params["URL"].(string)
							err := this.queue.AddURL(uri)
							if err != nil {
								this.log <- helper.FmtLog(common.LOG_ERROR, err.Error(), common.LOG_LEVEL_ERROR, common.LOG_TYPE_SYSTEM)
							}
						}
					} else if method == http.MethodPost || method == http.MethodPut {
						re := new(colly.Request)
						re.Method = method
						delete(params, "METHOD")
						if _, okey := params["URL"]; okey && len(params["URL"].(string)) > 0 {
							uri := params["URL"].(string)
							if u, err := url.Parse(uri); err == nil {
								re.URL = u
							}
							delete(params, "URL")
						}
						if _, okey := params["HEADER"]; okey {
							h := &http.Header{}
							for k, v := range params["HEADER"].(map[string]string) {
								h.Set(k, v)
							}
							re.Headers = h
							delete(params, "HEADER")
						}
						//生成请求

						if _, ok := params["BODY"]; ok {
							bodyType := reflect.TypeOf(params["BODY"])
							if bodyType.Kind() == reflect.String {
								re.Body = bytes.NewReader([]byte(params["BODY"].(string)))
							} else {
								bs, err := json.Marshal(params)
								if err == nil {
									re.Body = bytes.NewReader(bs)
								}
							}
						}
						err = this.queue.AddRequest(re)
						if err == nil {
							this.log <- helper.FmtLog(common.LOG_INFO, "加载页面:"+r.URL.String(), common.LOG_LEVEL_INFO, common.LOG_TYPE_SYSTEM)
							r.Abort()
						}
					}
				}
			}
		}
	}
}
