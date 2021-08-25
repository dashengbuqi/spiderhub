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

func (this *crawlerPool) Set(key string, sp *spider.Spider) bool {
	mu.Lock()
	defer mu.Unlock()
	this.data[key] = sp
	_, ok := this.data[key]
	return ok
}

func (this *crawlerPool) Delete(key string) {
	if _, ok := this.data[key]; ok {
		mu.Lock()
		delete(this.data, key)
		mu.Unlock()
	}
}

func (this *crawlerPool) Stop(key string) {
	if _, ok := this.data[key]; ok {
		mu.Lock()
		this.data[key].Stop()
		delete(this.data, key)
		mu.Unlock()
	}
}

func (this *crawlerPool) Exist(key string) bool {
	_, ok := this.data[key]
	return ok
}

func (this *crawlerPool) Get(key string) *spider.Spider {
	return this.data[key]
}
