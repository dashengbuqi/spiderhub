package crawler

import (
	"github.com/dashengbuqi/spiderhub/internal/crawler/spider"
	"sync"
)

type crawlerPool struct {
	data map[string]*spider.Spider
}

var (
	mu    sync.RWMutex
	Spool = &crawlerPool{
		data: make(map[string]*spider.Spider),
	}
)

//启动蜘蛛
func (this *crawlerPool) Start(key string, spd *spider.Spider) {
	has := this.set(key, spd)
	if has {
		this.data[key].Run()
	}
}

//停止蜘蛛
func (this *crawlerPool) Stop(key string) {
	has := this.Exist(key)
	if has {
		this.data[key].Stop()
		this.delete(key)
	}
}
func (this *crawlerPool) Exist(key string) bool {
	_, ok := this.data[key]
	return ok
}
func (this *crawlerPool) set(key string, spd *spider.Spider) bool {
	mu.Lock()
	defer mu.Unlock()
	this.data[key] = spd
	_, ok := this.data[key]
	return ok
}
func (this *crawlerPool) delete(key string) {
	if _, ok := this.data[key]; ok {
		mu.Lock()
		delete(this.data, key)
		mu.Unlock()
	}
}

func (this *crawlerPool) Get(key string) *spider.Spider {
	return this.data[key]
}
