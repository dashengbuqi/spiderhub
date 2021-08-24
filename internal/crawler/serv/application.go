package serv

import "github.com/dashengbuqi/spiderhub/middleware/queue"

var (
	CrawlerChannel = queue.Channel{
		Exchange:     "Crawlers",
		ExchangeType: "direct",
		RoutingKey:   "Request",
		Reliable:     true,
		Durable:      false,
	}
)

func RunApp() {

}
