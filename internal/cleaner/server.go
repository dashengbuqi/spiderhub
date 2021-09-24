package cleaner

import (
	"encoding/json"
	"github.com/dashengbuqi/spiderhub"
	"github.com/dashengbuqi/spiderhub/internal/common"
	"github.com/dashengbuqi/spiderhub/middleware/queue"
)

var (
	CleanerChannel = queue.Channel{
		Exchange:     "Cleaners",
		ExchangeType: "direct",
		RoutingKey:   "Request",
		Reliable:     true,
		Durable:      false,
	}
)

func RunServer() {
	c := make(chan []byte)
	go queue.RabbitConn.Consume(&CleanerChannel, c)

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
				go NewSchedule(cm).Run()
			}
		}
	}
}
