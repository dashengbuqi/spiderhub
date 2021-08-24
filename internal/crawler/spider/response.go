package spider

import (
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"github.com/dashengbuqi/spiderhub/internal/crawler/rules"
	"github.com/gocolly/colly"
	"strconv"
)

func (this *Spider) onResponse(r *colly.Response) {
	var isAllow = false
	if _, ok := this.rules[rules.ACCEPT_HTTP_STATUS]; ok {
		for _, code := range this.rules[rules.ACCEPT_HTTP_STATUS].([]int) {
			if r.StatusCode == code {
				isAllow = true
				break
			}
		}
	} else {
		isAllow = true
	}
	if isAllow == false {
		this.log <- helper.FmtLog(common.LOG_ERROR, "Http状态码不在允许范围内,编码:"+strconv.Itoa(r.StatusCode), common.LOG_LEVEL_INFO, common.LOG_TYPE_USER)
		return
	}
}
