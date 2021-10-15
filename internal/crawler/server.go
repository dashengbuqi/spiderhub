package crawler

import (
	"encoding/json"
	"github.com/dashengbuqi/spiderhub"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"github.com/dashengbuqi/spiderhub/middleware/queue"
	"time"
)

func RunServer() {
	c := make(chan []byte)

	go queue.RabbitConn.Consume(&common.CrawlerChannel, c)

	for {
		select {
		case msg := <-c:
			var cm common.Communication
			err := json.Unmarshal(msg, &cm)
			if err != nil {
				spiderhub.Logger.Error("%v", err)
				continue
			}
			if cm.Abort == true {
				Spool.SpiderStop(cm)
			} else {
				//如果是调试模式 等待 3s 再启动
				if cm.Method == common.SCHEDULE_METHOD_DEBUG {
					time.Sleep(time.Second * 1)
				}
				//运行调度器
				go NewSchedule(cm).Run()
			}
		}
	}
}
