package cleaner

import (
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"github.com/dashengbuqi/spiderhub/internal/crawler/spider"
	"sync"
)

type cleanerPool struct {
	data map[string]*spider.Spider
}

var (
	mu    sync.RWMutex
	Spool = &cleanerPool{
		data: make(map[string]*spider.Spider),
	}
)

//启动蜘蛛
func (this *cleanerPool) Start(key string, spd *spider.Spider) {
	mu.Lock()
	defer mu.Unlock()
	has := this.set(key, spd)
	if has {
		this.data[key].Run()
	}
}

//停止蜘蛛
func (this *cleanerPool) Stop(key string) {
	has := this.Exist(key)
	if has {
		mu.Lock()
		this.data[key].Stop()
		this.delete(key)
		mu.Unlock()
	}
}
func (this *cleanerPool) Exist(key string) bool {
	mu.RLock()
	defer mu.RUnlock()
	_, ok := this.data[key]
	return ok
}
func (this *cleanerPool) set(key string, spd *spider.Spider) bool {
	mu.Lock()
	defer mu.Unlock()
	this.data[key] = spd
	_, ok := this.data[key]
	return ok
}
func (this *cleanerPool) delete(key string) {
	if _, ok := this.data[key]; ok {
		mu.Lock()
		delete(this.data, key)
		mu.Unlock()
	}
}

func (this *cleanerPool) Get(key string) *spider.Spider {
	mu.RLock()
	defer mu.RUnlock()
	return this.data[key]
}

func (this *cleanerPool) SpiderStop(cm common.Communication) {
	key := helper.NewToken(cm.UserId, cm.AppId, cm.DebugId).Crawler().ToString()
	this.Stop(key)
}
