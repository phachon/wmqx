package container

import (
	"wmqx/message"
	"wmqx/pools"
	"errors"
	"strings"
	"net/http"
	"wmqx/app"
	"time"
	"io/ioutil"
	"strconv"
)

var Ctx = NewContext()

func NewContext() *Context {
	return &Context{
		QMessage:&message.QMessage{},
		RabbitMQPools: &pools.RabbitMQ{},
		ConsumerProcess: message.NewConsumerProcess(),
	}
}

type Context struct {

	// QMessage
	QMessage *message.QMessage

	// RabbitMQ pools
	RabbitMQPools *pools.RabbitMQ

	// Consumer Process
	ConsumerProcess *message.ConsumerProcess
}

// set RabbitMQ pools number and init
func (ctx *Context) SetRabbitMQPools(n int)  {
	ctx.RabbitMQPools = pools.NewRabbitMQPools()
	ctx.RabbitMQPools.Init(n)
}

// get consumerKey by messageName and consumerId
func (ctx *Context) GetConsumerKey(messageName string, consumerId string) string {
	return messageName +"_"+ consumerId
}

// split consumerKey
func (ctx *Context) SplitConsumerKey(consumerKey string) (messageName string, consumerId string){
	d := strings.Split(consumerKey, "_")
	if len(d) == 2 {
		return d[0], d[1]
	}else {
		return "", d[0]
	}
}

// init RabbitMQ all exchanges
func (ctx *Context) InitExchanges() error {
	rabbitMq, err := ctx.RabbitMQPools.GetMQ()
	if err != nil {
		return errors.New("rabbitmq pools faild: " + err.Error())
	}
	defer ctx.RabbitMQPools.Recover(rabbitMq)

	ctx.QMessage.Lock.Lock()
	defer ctx.QMessage.Lock.Unlock()

	messages := ctx.QMessage.Messages
	for _, msg := range messages {
		// declare exchange
		err := rabbitMq.DeclareExchange(msg.Name, msg.Mode, msg.Durable)
		if err != nil {
			return errors.New("Declare exchange faild: "+err.Error())
		}
		// declare queue
		for _, consumer := range msg.Consumers {
			consumerKey := ctx.GetConsumerKey(msg.Name, consumer.ID)
			err := rabbitMq.DeclareQueue(consumerKey, msg.Durable)
			if err != nil {
				return errors.New("Declare queue faild: "+err.Error())
			}
			// bind queue to exchange
			err = rabbitMq.BindQueueToExchange(consumerKey, msg.Name, consumer.RouteKey)
			if err != nil {
				return errors.New("Bind queue exchange fail: "+err.Error())
			}
			Worker.SendConsumerSign(Consumer_Action_Insert, consumerKey)
		}
	}
	return nil
}

// request consumer url by consumerKey
func (ctx *Context) RequestConsumer(consumerKey string, publishMessage message.PublishMessage) (resBody string, respCode int, err error) {

	// get consumer info
	messageName, consumerId := ctx.SplitConsumerKey(consumerKey)
	consumer, err := ctx.QMessage.GetConsumerById(messageName, consumerId)
	if err != nil {
		return
	}

	url := consumer.URL
	timeout := consumer.Timeout
	code := consumer.Code
	checkCode := consumer.CheckCode
	method := publishMessage.Method
	body := publishMessage.Body
	args := publishMessage.Args
	ip := publishMessage.Ip
	headers := publishMessage.Header
	if !strings.Contains(url, "?") {
		url += "?"
	}
	if args != "" {
		url += args
	}
	var req *http.Request
	if method == "POST" {
		req, err = http.NewRequest("POST", url, strings.NewReader(body))
	}else if method == "GET" {
		req, err = http.NewRequest("GET", url, nil)
	}
	if err != nil {
		return
	}
	req.Header.Set(app.Conf.GetString("publish.RealIpHeader"), ip)
	req.Header.Set("User-Agent", "WMQX version" + app.AppVersion + " - https://github.com/phachon/wmqx")
	if (headers != nil) && (len(headers) > 0) {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	client := &http.Client{}
	client.Timeout = time.Duration(time.Duration(timeout) * time.Millisecond)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	respCode = resp.StatusCode
	defer resp.Body.Close()

	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if checkCode && (code != float64(respCode)) {
		err = errors.New("response code error: "+strconv.Itoa(respCode))
		return
	}

	return string(bodyByte), respCode, nil
}