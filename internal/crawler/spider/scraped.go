package spider

import (
	"github.com/dashengbuqi/spiderhub/internal/crawler"
	"github.com/gocolly/colly"
)

//响应回调函数
func (this *Spider) onScraped(r *colly.Response) {
	this.container.Call(crawler.FUNC_AFTER_CRAWL, nil, r.Request.URL.String())
}
