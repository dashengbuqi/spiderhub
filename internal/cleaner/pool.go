package cleaner

import (
	"github.com/dashengbuqi/spiderhub/helper"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"sync"
)

type cleanerPool struct {
	data map[string]*Cleaner
}

var (
	mu        sync.RWMutex
	CleanPool = &cleanerPool{
		data: make(map[string]*Cleaner),
	}
)

//启动
func (this *cleanerPool) Start(key string, cl *Cleaner) {
	mu.Lock()
	defer mu.Unlock()
	has := this.set(key, cl)
	if has {
		this.data[key].Run()
	}
}

//停止
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
func (this *cleanerPool) set(key string, cl *Cleaner) bool {
	mu.Lock()
	defer mu.Unlock()
	this.data[key] = cl
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

func (this *cleanerPool) Get(key string) *Cleaner {
	mu.RLock()
	defer mu.RUnlock()
	return this.data[key]
}

func (this *cleanerPool) CleanStop(cm common.Communication) {
	key := helper.NewToken(cm.UserId, cm.AppId, cm.DebugId).Pool().ToString()
	this.Stop(key)
}
