package crawler

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"github.com/gocolly/colly"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func (this *Spider) onResponse(r *colly.Response) {
	var isAllow = false
	if _, ok := this.params[ACCEPT_HTTP_STATUS]; ok && len(this.params[ACCEPT_HTTP_STATUS].([]interface{})) > 0 {
		for _, code := range this.params[ACCEPT_HTTP_STATUS].([]interface{}) {
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
		this.success(r)
	} else if r.StatusCode == HTTP_STATUS_FORBIDDEN {
		this.forbidden(r)
	} else {
		this.otherError(r)
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
		if len(this.params[SCAN_URLS].([]string)) > 0 {
			for _, su := range this.params[SCAN_URLS].([]string) {
				if su == r.Request.URL.String() {
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
								this.autoFindURL(body, r)
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
					this.autoFindURL(body, r)
				}
			}
		}
	}
	//内容页
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
					return
				}
			}
			//内容页
			if res, err := this.container.Call(FUNC_ON_PROCESS_CONTENT_PAGE, nil, body, this.queue); err == nil {
				if state, _ := res.ToBoolean(); state == true {
					this.autoFindURL(body, r)
				}
			}
			this.mu.Lock()
			respData := this.extract(body, r.Request.URL.String())
			if len(respData) > 0 {
				respData["target_url"] = map[bool]interface{}{
					false: r.Request.URL.String(),
				}
				this.outData <- respData
			}
			if this.method == common.SCHEDULE_METHOD_DEBUG {
				if this.runTimes > 10 {
					this.outLog <- common.FmtLog(common.LOG_INFO, "调试模式结束", common.LOG_LEVEL_INFO, common.LOG_TYPE_SYSTEM)
					this.abort = true
				}
				this.runTimes++
			}
			this.mu.Unlock()
		}
	}
}

func (this *Spider) extract(body string, curl string) map[string]interface{} {
	defer func() {
		if p := recover(); p != nil {
			str, ok := p.(string)
			var err error
			if ok {
				err = errors.New(str)
			} else {
				err = errors.New("异常提取内容")
			}
			this.outLog <- common.FmtLog("异常", err.Error(), common.LOG_LEVEL_ERROR, common.LOG_TYPE_SYSTEM)
		}
	}()
	fields := this.params[FIELDS].([]FieldStash)
	result := make(map[string]interface{})
	if len(fields) > 0 {
		for _, field := range fields {
			if field.Temporary == false {
				data := this.recursExtract(body, field, curl)
				if data == nil {
					this.outLog <- common.FmtLog("异常", "字段:"+field.Name, common.LOG_LEVEL_ERROR, common.LOG_TYPE_SYSTEM)
					break
				}
				result[field.Name] = data
			}
		}
	}

	//整个网页完成抽取时回调此函数。一般在此回调中做一些数据整理的操作 FUNC_AFTER_EXTRACT_PAGE
	if res, err := this.container.Call(FUNC_AFTER_EXTRACT_PAGE, nil, result); err == nil {
		if res.IsDefined() {
			keys := res.Object().Keys()
			if len(keys) > 0 {
				for _, key := range keys {
					val, _ := res.Object().Get(key)
					result[key] = val
				}
			}
		}
	}
	return result
}

//递归的提取字段
func (this *Spider) recursExtract(body string, field FieldStash, curl string) map[bool]interface{} {
	defer func() {
		if p := recover(); p != nil {
			str, ok := p.(string)
			var err error
			if ok {
				err = errors.New(str)
			} else {
				err = errors.New("异常提取内容")
			}
			this.outLog <- common.FmtLog("异常", err.Error(), common.LOG_LEVEL_ERROR, common.LOG_TYPE_SYSTEM)
		}
	}()
	data := make(map[bool]interface{})
	valueType := field.Type
	//设置默认值
	if len(valueType) == 0 || helper.ValueInArray(valueType, TypeArr) == false {
		valueType = TYPE_STRING
	}
	//字符串处理
	if valueType == TYPE_STRING {
		if field.ExtractMethod == EXTRACT_ATTACHEDURL {
			//表示需要的数据在另外一个链接（我们叫attachedUrl）的请求结果里面，需要额外再发一次请求来获取数据。
			u := helper.ExtractItem(body, field.Selector, field.SelectorType)
			if len(u.(string)) == 0 {
				data = map[bool]interface{}{
					field.Primary: "",
				}
			} else {
				var full string
				if len(field.Func) > 0 {
					uri, _ := this.container.Call(field.Func, nil, u, curl)
					full = uri.String()
				} else {
					full = u.(string)
				}
				//提取到内容后回调函数
				if res, err := this.container.Call(FUNC_AFTER_EXTRACT_FIELD, nil, full); err == nil {
					if res.IsDefined() {
						full = res.String()
					}
				}
				requestMethod := "GET"
				if len(field.AttachedMethod) > 0 {
					requestMethod = field.AttachedMethod
				}
				var header http.Header
				if len(field.AttachedHeaders) > 0 {
					for k, v := range field.AttachedHeaders {
						header.Add(k, v)
					}
				}
				subBody := this.outSite(full, requestMethod, header)
				if len(subBody) > 0 {
					subData := make(map[string]map[bool]interface{})
					for _, subField := range field.Children {
						subData[subField.Name] = this.recursExtract(subBody, subField, curl)
					}
					data = map[bool]interface{}{
						field.Primary: subData,
					}
				}
			}
		} else if field.ExtractMethod == EXTRACT_NORMAL {
			var item interface{}
			if field.BodyType == BODY_HTML {
				item = helper.ExtractHtml(body, field.Selector, field.SelectorType)
			} else {
				item = helper.ExtractItem(body, field.Selector, field.SelectorType)
			}
			//如果必须要有值且为空的情况直接返回
			if field.Required && item == nil {
				return nil
			}
			//临时字段
			if field.Temporary {
				if res, err := this.container.Call(FUNC_AFTER_EXTRACT_TEMPORARY_FIELD, nil, item); err == nil {
					if res.IsDefined() {
						item = res.String()
					}
				}
			}
			//提取内容回调函数
			if res, err := this.container.Call(FUNC_AFTER_EXTRACT_FIELD, nil, item); err == nil {
				if res.IsDefined() {
					item = res.String()
				}
			}
			//自定义函数调用
			if len(field.Func) > 0 {
				funData, _ := this.container.Call(field.Func, nil, item, curl)
				if funData.IsNull() {
					data = map[bool]interface{}{
						field.Primary: strings.TrimSpace(strings.Trim(item.(string), "\n")),
					}
				} else {
					data = map[bool]interface{}{
						field.Primary: funData.String(),
					}
				}
			} else {
				data = map[bool]interface{}{
					field.Primary: strings.TrimSpace(strings.Trim(item.(string), "\n")),
				}
			}
			//检查是否有子项
			if len(field.Children) > 0 {
				subData := make(map[string]map[bool]interface{})
				for _, subField := range field.Children {
					subData[subField.Name] = this.recursExtract(body, subField, curl)
				}
				data = map[bool]interface{}{
					field.Primary: subData,
				}
			}
		}
	} else if valueType == TYPE_ARRAY {
		item := helper.Extracts(body, field.Selector, field.SelectorType)
		if len(field.Func) > 0 {
			funData, _ := this.container.Call(field.Func, nil, item, curl)
			if funData.IsUndefined() {
				data = map[bool]interface{}{
					field.Primary: item,
				}
			} else {
				keys := funData.Object().Keys()
				var dataArr []interface{}
				if len(keys) > 0 {
					for _, key := range keys {
						val, _ := funData.Object().Get(key)
						dataArr = append(dataArr, strings.TrimSpace(strings.Trim(val.String(), "\n")))
					}
					data = map[bool]interface{}{
						field.Primary: dataArr,
					}
				} else {
					data = map[bool]interface{}{
						field.Primary: funData,
					}
				}
			}
		} else {
			data = map[bool]interface{}{
				field.Primary: item,
			}
		}
		if len(field.Children) > 0 {
			subData := make(map[string]map[bool]interface{})
			for _, subField := range field.Children {
				subData[subField.Name] = this.recursExtract(body, subField, curl)
			}
			data = map[bool]interface{}{
				field.Primary: subData,
			}
		}
	} else if valueType == TYPE_MAP {
		item := helper.Extracts(body, field.Selector, field.SelectorType)
		if len(field.Func) > 0 {
			funData, _ := this.container.Call(field.Func, nil, item, curl)
			if funData.IsUndefined() {
				data = map[bool]interface{}{
					field.Primary: item,
				}
			} else {
				keys := funData.Object().Keys()
				dataMap := make(map[string]interface{})
				if len(keys) > 0 {
					for _, key := range keys {
						val, _ := funData.Object().Get(key)
						dataMap[key] = strings.TrimSpace(strings.Trim(val.String(), "\n"))
						//dataArr = append(dataArr, strings.TrimSpace(strings.Trim(val.String(), "\n")))
					}
					data = map[bool]interface{}{
						field.Primary: dataMap,
					}
				} else {
					data = map[bool]interface{}{
						field.Primary: funData,
					}
				}
			}
		} else {
			data = map[bool]interface{}{
				field.Primary: item,
			}
		}
		if len(field.Children) > 0 {
			subData := make(map[string]map[bool]interface{})
			for _, subField := range field.Children {
				subData[subField.Name] = this.recursExtract(body, subField, curl)
			}
			data = map[bool]interface{}{
				field.Primary: subData,
			}
		}
	} else {
		var item interface{}
		item = helper.ExtractItem(body, field.Selector, field.SelectorType)
		if field.Required && item == nil {
			return nil
		}
		//临时字段
		if field.Temporary {
			if res, err := this.container.Call(FUNC_AFTER_EXTRACT_FIELD, nil, item); err == nil {
				item = res
			}
		}
		//抽取到内容后回调此函数
		if res, err := this.container.Call(FUNC_AFTER_EXTRACT_FIELD, nil, item); err == nil {
			item = res
		}
		data = map[bool]interface{}{
			field.Primary: item,
		}
		if len(field.Children) > 0 {
			subData := make(map[string]map[bool]interface{})
			for _, subFsItem := range field.Children {
				subData[subFsItem.Name] = this.recursExtract(body, subFsItem, curl)
			}
			data = map[bool]interface{}{
				field.Primary: subData,
			}
		}
	}
	if len(data) > 0 {
		this.container.Call(FUNC_ON_DATA_RECEIVED, nil, data[field.Primary], this.queue)
	}
	return data
}

//请求外链
func (this *Spider) outSite(uri string, method string, header http.Header) string {
	req, _ := http.NewRequest(method, uri, nil)
	if len(header) > 0 {
		req.Header = header
	}
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode != HTTP_STATUS_SUCCESS {
		this.outLog <- common.FmtLog("错误", "抓取外链失败", common.LOG_LEVEL_ERROR, common.LOG_TYPE_SYSTEM)
		return ""
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func (this *Spider) autoFindURL(body string, r *colly.Response) {
	urls := helper.AutoFindLinkUrls(body)
	if len(urls) > 0 {
		for _, url := range urls {
			this.outLog <- common.FmtLog("信息", url, common.LOG_LEVEL_INFO, common.LOG_TYPE_URL)
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
