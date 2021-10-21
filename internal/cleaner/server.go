package cleaner

import (
	"encoding/json"
	"github.com/dashengbuqi/spiderhub"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"github.com/dashengbuqi/spiderhub/middleware/queue"
	"time"
)

func RunServer() {
	c := make(chan []byte)
	go queue.RabbitConn.Consume(&queue.CleanerChannel, c)

	for {
		select {
		case msg := <-c:
			var cm common.Communication
			err := json.Unmarshal(msg, &cm)
			if err != nil {
				spiderhub.Logger.Error("%s", err.Error())
				continue
			}
			if cm.Abort {
				go CleanPool.CleanStop(cm)
			} else {
				if cm.Method == common.SCHEDULE_METHOD_DEBUG {
					time.Sleep(time.Second * 3)
				}
				go NewSchedule(cm).Run()
			}
		}
	}
}
