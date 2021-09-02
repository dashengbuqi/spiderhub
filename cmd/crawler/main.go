package main

import (
	"github.com/dashengbuqi/spiderhub/internal/crawler"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	crawler.RunServer()
}
