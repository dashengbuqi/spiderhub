package crawler

import (
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"sync"
)

type crawlerPool struct {
	data map[string]*Spider
}

var (
	mu    sync.RWMutex
	Spool = &crawlerPool{
		data: make(map[string]*Spider),
	}
)

//启动蜘蛛
func (this *crawlerPool) Start(key string, spd *Spider) {
	has := this.set(key, spd)
	if has {
		this.data[key].Run()
	}
}

//停止蜘蛛
func (this *crawlerPool) Stop(key string) {
	has := this.Exist(key)
	if has {
		mu.Lock()
		this.data[key].Stop()
		mu.Unlock()
	}
}
func (this *crawlerPool) Exist(key string) bool {
	mu.RLock()
	defer mu.RUnlock()
	_, has := this.data[key]
	return has
}
func (this *crawlerPool) set(key string, spd *Spider) bool {
	mu.Lock()
	defer mu.Unlock()
	this.data[key] = spd
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

func (this *crawlerPool) Get(key string) *Spider {
	mu.RLock()
	defer mu.RUnlock()
	return this.data[key]
}

func (this *crawlerPool) SpiderStop(cm common.Communication) {
	key := helper.NewToken(cm.UserId, cm.AppId, cm.DebugId).Pool().ToString()
	this.Stop(key)
}
