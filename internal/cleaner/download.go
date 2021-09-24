package cleaner

import (
	"errors"
	"github.com/dashengbuqi/spiderhub"
	"github.com/dashengbuqi/spiderhub/configs"
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
)

//下载附件
type Download struct {
	data      *common.FieldData
	field     FieldStash
	primary   interface{}
	container *otto.Otto
	token     string
}

func NewDownload(d *common.FieldData, p interface{}, f FieldStash, t string, c *otto.Otto) *Download {
	return &Download{
		data:      d,
		field:     f,
		primary:   p,
		container: c,
		token:     t,
	}
}

func (this *Download) Run() error {
	defer func() {
		p := recover()
		if p != nil {
			spiderhub.Logger.Error("下载失败:%v", p.(error).Error())
		}
	}()
	var list []string
	if this.data.Type == TYPE_STRING {
		if len(this.data.Value.(string)) > 0 && this.data.Value.(string) != "undefined" {
			list = append(list, this.data.Value.(string))
		}
	} else if this.data.Type == TYPE_ARRAY || this.data.Type == TYPE_MAP {
		if reflect.TypeOf(this.data.Value).Kind() == reflect.Array {
			for _, uri := range this.data.Value.([]interface{}) {
				list = append(list, uri.(string))
			}
		} else if reflect.TypeOf(this.data.Value).Kind() == reflect.Map {
			for _, uri := range this.data.Value.(map[string]interface{}) {
				list = append(list, uri.(string))
			}
		} else if reflect.TypeOf(this.data.Value).Kind() == reflect.String {
			uri := this.data.Value.(string)
			if strings.Index(uri, "[") > -1 || strings.Index(uri, "]") > -1 {
				uri = strings.Replace(uri, "[", "", -1)
				uri = strings.Replace(uri, "]", "", -1)
				list = strings.Split(uri, ",")
			}
		}
	}
	var err error
	if len(list) > 0 {
		for _, uri := range list {
			if len(uri) == 0 {
				continue
			}
			uri = helper.FmtUrl(uri)
			//下载附件前的确认
			if res, err := this.container.Call(FUNC_BEFORE_HANDLE_IMG, nil, this.field.Name, uri); err == nil {
				if res.IsDefined() {
					uri = res.String()
				}
			}
			u, _ := url.Parse(uri)
			if len(u.Scheme) == 0 {
				u.Scheme = "https"
			}
			uri = u.String()
			var status int
			var header http.Header
			var body []byte
			var contentType string
			//请求目标地址
			deep := 3
			var i int
			for {
				if i++; i > deep {
					break
				}
				status, header, body = this.process(uri)
				contentType = http.DetectContentType(body)
				if status != HTTP_STATUS_SUCCESS {
					err = errors.New("附件地址请求失败")
				}
				//是如果是网页,需要回调用继续获取附件地址,用于下载地址隐藏下内容页中
				if strings.Contains(contentType, "text/html") || strings.Contains(contentType, "text/plain") {
					if res, err := this.container.Call(FUNC_BEFORE_HOST_FILE, nil, contentType, string(body)); err == nil {
						if res.IsDefined() {
							uri = res.String()
						}
					} else {
						//没有设置就退出
						break
					}
				} else {
					//不是文本就退出循环
					break
				}

			}
			//最终还是网页内容就直接返回错误
			if strings.Contains(contentType, "text/html") || strings.Contains(contentType, "text/plain") {
				err = errors.New("未获取到附件地址")
				continue
			}
			//body保存到对应的目录
			err = this.save(uri, header, body)
		}
	}
	return err
}

//保存附件
func (this *Download) save(uri string, header http.Header, body []byte) error {
	base, _ := configs.GetParamsByField("Common", "AttachPath")
	name := helper.GetFileName(uri)
	fullPath := base.(string) + "/" + this.token + "/" + helper.StringToM5(this.primary.(string))
	if helper.HasFile(fullPath) == false {
		if err := helper.MakeDir(fullPath); err != nil {
			return err
		}
	}
	var cd string
	cd = header.Get("Content-Disposition")
	if len(cd) == 0 {
		cd = header.Get("content-disposition")
	}
	if len(cd) > 0 {
		cdArr := strings.Split(cd, ";")
		for _, c := range cdArr {
			if strings.Index(c, "=") > 0 {
				cArr := strings.Split(c, "=")
				if strings.ToLower(cArr[0]) == "filename" {
					name = strings.TrimSpace(strings.Replace(cArr[1], "\"", "", -1))
				}
			}
		}
	}
	fp, err := os.Create(fullPath + "/" + name)
	if err != nil {
		return err
	}
	defer fp.Close()
	err = helper.CopyFF(strings.NewReader(string(body)), fp)
	return err
}

func (this *Download) process(uri string) (int, http.Header, []byte) {
	method := "GET"
	if len(this.field.AttachedMethod) > 0 {
		method = this.field.AttachedMethod
	}
	var params string
	if len(this.field.AttachedUrlParams) > 0 {
		params = this.field.AttachedUrlParams
	}
	req, _ := http.NewRequest(method, uri, strings.NewReader(params))
	if this.field.AttachedHeaders != nil {
		for key, val := range this.field.AttachedHeaders {
			req.Header.Add(key, val.(string))
		}
	}
	cli := &http.Client{}
	resp, _ := cli.Do(req)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, resp.Header, body
}
