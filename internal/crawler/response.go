package crawler

import (
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
}
