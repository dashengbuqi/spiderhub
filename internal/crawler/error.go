package crawler

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"github.com/gocolly/colly"
	"strings"
	"time"
)

func (this *Spider) onError(resp *colly.Response, e error) {
	this.outLog <- common.FmtLog(common.LOG_ERROR, e.Error()+": "+resp.Request.URL.String(), common.LOG_LEVEL_ERROR, common.LOG_TYPE_SYSTEM)

	h := md5.New()
	h.Write([]byte(strings.ToLower(resp.Request.URL.String())))
	key := hex.EncodeToString(h.Sum(nil))
	this.mu.Lock()
	defer this.mu.Unlock()
	//如果存在了就直接删除，否则加入执行队列
	if _, ok := this.failure[key]; ok {
		if this.failure[key] >= 3 {
			delete(this.failure, key)
			return
		}
		this.failure[key]++
		this.queue.AddURL(resp.Request.URL.String())
		time.Sleep(1 * time.Second)
	} else {
		this.failure[key] = 1
		this.queue.AddURL(resp.Request.URL.String())
		time.Sleep(1 * time.Second)
	}
}
