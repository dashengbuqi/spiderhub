package crawler

import (
	"github.com/gocolly/colly"
)

//响应回调函数
func (this *Spider) onScraped(r *colly.Response) {
	this.container.Call(FUNC_AFTER_CRAWL, nil, r.Request.URL.String())
}
