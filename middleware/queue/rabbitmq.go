package queue

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/dashengbuqi/spiderhub"

	"github.com/dashengbuqi/spiderhub/configs"

	"github.com/streadway/amqp"
	"strconv"
	"sync"
	"time"
)

var (
	uri        string
	queueOnce  sync.Once
	RabbitConn *base
)

type Conn struct {
	Locker     sync.RWMutex
	Connection *amqp.Connection
	rabbitUri  string
}

type Channel struct {
	Exchange     string
	ExchangeType string
	RoutingKey   string
	Reliable     bool
	Durable      bool
	ChannelId    string
	Channel      *amqp.Channel
}

type base struct {
	Conn     *Conn
	Channels map[string]*Channel
}

func init() {
	uris, _ := configs.GetParamsByField("Queue", "Uri")
	if len(uris.(string)) == 0 {
		panic("queue start failure...")
	}
	uri = uris.(string)
	RabbitConn = getQueueConn()
}

func getQueueConn() (bs *base) {
	queueOnce.Do(func() {
		bs = &base{
			Conn:     &Conn{rabbitUri: uri},
			Channels: map[string]*Channel{},
		}
	})
	return bs
}

func (this *base) confirmOne(c <-chan amqp.Confirmation) {
	<-c
}

func (this *base) buildToken(c *Channel) string {
	token := c.Exchange + ":" + c.ExchangeType + ":" + c.RoutingKey + ":" + strconv.FormatBool(c.Durable) + ":" + strconv.FormatBool(c.Reliable)
	m := md5.New()
	m.Write([]byte(token))
	return hex.EncodeToString(m.Sum(nil))
}

func (this *base) refresh(c *Channel) error {
	this.Conn.Locker.Lock()
	defer this.Conn.Locker.Unlock()

	var err error
	if this.Conn.Connection != nil {
		c.Channel, err = this.Conn.Connection.Channel()
	} else {
		err = errors.New("队列链接失败")
	}
	if err != nil {
		for {
			this.Conn.Connection, err = amqp.Dial(this.Conn.rabbitUri)
			if err != nil {
				time.Sleep(3 * time.Second)
			} else {
				c.Channel, _ = this.Conn.Connection.Channel()
				break
			}
		}
	}

	err = c.Channel.ExchangeDeclare(
		c.Exchange,
		c.ExchangeType,
		c.Durable,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	this.Channels[c.ChannelId] = c
	return nil
}

func (this *base) Publish(c *Channel, body []byte) error {
	c.ChannelId = this.buildToken(c)

	if this.Channels[c.ChannelId] == nil {
		err := this.refresh(c)
		if err != nil {
			return err
		}
		c = this.Channels[c.ChannelId]
	} else {
		c = this.Channels[c.ChannelId]
	}

	var times int
	var err error

	for {
		err = c.Channel.Publish(
			c.Exchange,
			c.RoutingKey,
			false,
			false,
			amqp.Publishing{
				ContentType:  "text/plain", //"application/json",
				Body:         body,
				DeliveryMode: amqp.Transient,
				Priority:     0,
			})
		if err != nil {
			time.Sleep(1 * time.Second)
			if times < 3 {
				err = this.refresh(c)
				if err == nil {
					c = this.Channels[c.ChannelId]
				}
			} else {
				err = errors.New("Rabiitmq Conn failure")
				break
			}
			times++
		}
	}
	return err
}

func (this *base) Consume(c *Channel, out chan<- []byte) {
	c.ChannelId = this.buildToken(c)

	if this.Channels[c.ChannelId] == nil {
		err := this.refresh(c)
		if err != nil {
			spiderhub.Logger.Error("%v", err)
		}
	} else {
		c = this.Channels[c.ChannelId]
	}
	queue, err := c.Channel.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		spiderhub.Logger.Error("%v", err)
		return
	}
	err = c.Channel.QueueBind(queue.Name, c.RoutingKey, c.Exchange, false, nil)
	if err != nil {
		spiderhub.Logger.Error("%v", err)
		return
	}
	deli, err := c.Channel.Consume(queue.Name, "", true, false, false, false, nil)

	if err != nil {
		spiderhub.Logger.Error("%v", err)
		return
	}

	ch := make(chan bool)
	go func() {
		for d := range deli {
			out <- d.Body
		}
	}()
	<-ch
}
