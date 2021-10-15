package crawler

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"github.com/gocolly/colly"
	"regexp"
	"strconv"
	"strings"
)

func (this *Spider) onResponse(r *colly.Response) {
	var isAllow = false
	if _, ok := this.params[ACCEPT_HTTP_STATUS]; ok {
		for _, code := range this.params[ACCEPT_HTTP_STATUS].([]int) {
			if r.StatusCode == code {
				isAllow = true
				break
			}
		}
	} else {
		isAllow = true
	}
	if isAllow == false {
		this.outLog <- common.FmtLog(common.LOG_ERROR, "Http状态码不在允许范围内,状态码:"+strconv.Itoa(r.StatusCode), common.LOG_LEVEL_INFO, common.LOG_TYPE_USER)
		return
	}
	if r.StatusCode == HTTP_STATUS_SUCCESS {

	}
}

func (this *Spider) success(r *colly.Response) {
	body := string(r.Body)
	if helper.ValidUTF8(r.Body) == false {
		body = helper.ConvertToString(body, "gbk", "utf-8")
	}
	//入口
	if _, ok := this.params[SCAN_URLS]; ok {
		autoFindUrls := this.params[AUTOFIND_URLS].(bool)
		if len(this.params[SCAN_URLS].([]interface{})) > 0 {
			for _, su := range this.params[SCAN_URLS].([]interface{}) {
				if su.(string) == r.Request.URL.String() {
					if autoFindUrls {
						urls := helper.AutoFindLinkUrls(body)
						if len(urls) > 0 {
							for _, url := range urls {
								for _, regex := range this.params[CONTENT_REGEX].([]string) {
									if m, err := regexp.MatchString(regex, url); err == nil && m == true {
										this.outLog <- common.FmtLog("提示", url, common.LOG_LEVEL_INFO, common.LOG_TYPE_URL)
										if this.abort == false {
											if strings.Contains(url, "http") || strings.Contains(url, "https") {
												this.queue.AddURL(url)
											} else {
												//[scheme:][//[userinfo@]host][/]path[?query][#fragment]
												url = strings.TrimPrefix(url, "/")
												this.queue.AddURL(r.Request.URL.Scheme + "://" + r.Request.URL.Host + "/" + url)
											}
										}
									}
								}

							}
						}
					} else {
						if rs, err := this.container.Call(FUNC_ON_PROCESS_SCAN_PAGE, nil, body, this.queue); err == nil {
							if status, _ := rs.ToBoolean(); status == true {
								urls := helper.AutoFindLinkUrls(body)
								if len(urls) > 0 {
									for _, url := range urls {
										this.outLog <- common.FmtLog("发现新网页", url, common.LOG_LEVEL_INFO, common.LOG_TYPE_URL)
										if this.abort == false {
											if strings.Contains(url, "http") || strings.Contains(url, "https") {
												this.queue.AddURL(url)
											} else {
												//[scheme:][//[userinfo@]host][/]path[?query][#fragment]
												url = strings.TrimPrefix(url, "/")
												this.queue.AddURL(r.Request.URL.Scheme + "://" + r.Request.URL.Host + "/" + url)
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	//列表
	if _, ok := this.params[HELPER_REGEX]; ok {
		var valid bool
		for _, re := range this.params[HELPER_REGEX].([]string) {
			if m, _ := regexp.MatchString(re, r.Request.URL.String()); m == true {
				valid = true
				break
			}
		}
		if valid {
			if res, err := this.container.Call(FUNC_ON_PROCESS_HELPER_PAGE, nil, body, r.Request.URL.String(), this.queue); err == nil {
				if status, _ := res.ToBoolean(); status == true {
					urls := helper.AutoFindLinkUrls(body)
					if len(urls) > 0 {
						for _, url := range urls {
							for _, regex := range this.params[CONTENT_REGEX].([]string) {
								if m, err := regexp.MatchString(regex, url); err == nil && m == true {
									this.outLog <- common.FmtLog("提示", url, common.LOG_LEVEL_INFO, common.LOG_TYPE_SYSTEM)
									if this.abort == false {
										if strings.Contains(url, "http") || strings.Contains(url, "https") {
											this.queue.AddURL(url)
										} else {
											//[scheme:][//[userinfo@]host][/]path[?query][#fragment]
											url = strings.TrimPrefix(url, "/")
											this.queue.AddURL(r.Request.URL.Scheme + "://" + r.Request.URL.Host + "/" + url)
										}
									}
								}
							}

						}
					}
				}
			}
		}
	}
	//内容页 @todo 未完成
	if _, ok := this.params[CONTENT_REGEX]; ok {
		var valid bool
		for _, re := range this.params[CONTENT_REGEX].([]string) {
			if m, _ := regexp.MatchString(re, r.Request.URL.String()); m == true {
				valid = true
				break
			}
		}
		if valid {
			res, _ := this.container.Call(FUNC_IS_ANTI_SPIDER, nil, r.Request.URL.String(), body)
			if res.IsDefined() == true {
				if has, _ := res.ToBoolean(); has {
					this.forbidden(r)
				}
			}
		}
	}
}

//每个被调度的链接下载完成之后，会先判断返回的状态码是否403，如果403，则直接会触发切换代理；
// 如果不是403，则回调此函数，开发者一般需要在此函数中判断返回码或者网页内容，
// 给出是否反爬的判断，如果判断为反爬，需要返回true，否则返回false。
func (this *Spider) forbidden(r *colly.Response) {
	//可能是因为认定为爬虫导致，可切换代理
	this.outLog <- common.FmtLog("异常", "被认定为爬虫停止执行,地址:"+r.Request.URL.String(), common.LOG_LEVEL_ERROR, common.LOG_TYPE_SYSTEM)
	//切换代理IP交需要重新加入队列

	//转给错误请求
	this.failureRequest(r.Request.URL.String())
}

func (this *Spider) otherError(r *colly.Response) {
	this.outLog <- common.FmtLog("错误", "错误码："+strconv.Itoa(r.StatusCode), common.LOG_LEVEL_ERROR, common.LOG_TYPE_SYSTEM)
	//转给错误请求
	this.failureRequest(r.Request.URL.String())
}

func (this *Spider) failureRequest(url string) {
	h := md5.New()
	h.Write([]byte(strings.ToLower(url)))
	key := hex.EncodeToString(h.Sum(nil))
	//如果存在了就直接删除，否则加入执行队列
	if _, ok := this.failure[key]; ok {
		if this.failure[key] >= 3 {
			delete(this.failure, key)
			return
		}
		this.failure[key]++
		this.queue.AddURL(url)
	} else {
		this.failure[key] = 1
		this.queue.AddURL(url)
	}
}
