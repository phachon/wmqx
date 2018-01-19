package container

import (
	"github.com/spf13/viper"
	"github.com/spf13/pflag"
	"flag"
	"fmt"
	"os"
	"strings"
	"github.com/fatih/color"
	"rmqc/mq"
	"rmqc/pools"
	"errors"
)

var (
	Ctx = NewContainer()
)

type Container struct {
	Config *viper.Viper
	QMessage *mq.QMessage
	RabbitMQPools *pools.RabbitMQ
	ConsumerWorkChan chan ConsumerWorker
}

type ConsumerWorker struct {
	Action string
	QueueName string
}

// return container object
func NewContainer() *Container {
	ct := &Container{
		Config: viper.New(),
		QMessage: mq.NewQMessage(),
		RabbitMQPools: pools.NewRabbitMQPools(),
		ConsumerWorkChan: make(chan ConsumerWorker, 1),
	}
	return ct
}

// init config
func (container *Container) InitConfig() {
	cfg := container.Config

	cfg.SetDefault("wmq.version", "1.5")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.String("listen-api", "0.0.0.0:3302", "api service listening port")
	pflag.String("listen-publish", "0.0.0.0:3303", "publish service listening port")
	pflag.String("api-token", "guest", "access api token")
	configFile := pflag.String("config", "", "config file path")
	pflag.Bool("api-disable", false, "disable api service")
	pflag.String("level", "debug", "console log level,should be one of debug,info,warn,error")
	version := pflag.Bool("version", false, "show version about current WMQ")
	example := pflag.Bool("data-example", false, "print example of data-file")
	pflag.StringSlice("ignore-headers", []string{}, "these http headers will be ignored when access to consumer's url , multiple splitted by comma(,)")
	pflag.String("realip-header", "X-Forwarded-For", "the publisher's real ip will be set in this http header when access to consumer's url")
	pflag.Int("fail-wait", 50, "access consumer url  fail and then how many seconds to sleep  and retry")
	pflag.Int("go-fail-wait", 3, "consumer's goroutine occur error and then how many seconds to sleep and retry")
	pflag.String("mq-host", "127.0.0.1", "which host be used when connect to RabbitMQ")
	pflag.Int("mq-port", 5672, "which port be used when connect to RabbitMQ")
	pflag.String("mq-username", "guest", "which username be used when connect to RabbitMQ")
	pflag.String("mq-password", "guest", "which password be used when connect to RabbitMQ")
	pflag.String("mq-vhost", "/", "which vhost be used when connect to RabbitMQ")
	pflag.String("mq-prefix", "wmq.", "the queue and exchange default prefix")
	pflag.String("data-file", "message.json", "which file will store messages")
	pflag.String("log-dir", "log", "the directory which store log files")
	pflag.Bool("log-access", true, "access log on or off")
	pflag.Bool("log-post", false, "log post data on or off")

	pflag.Int64("log-max-size", 102400000, "log file max size(bytes) for rotate")
	pflag.Int("log-max-count", 3, "log file max count for rotate to remain")
	pflag.StringSlice("log-level", []string{"info", "error", "debug"}, "log to file level,multiple splitted by comma(,)")
	pflag.Parse()
	if *version {
		fmt.Printf("WMQ v%s - https://github.com/snail007/wmq\n", cfg.GetString("wmq.version"))
		os.Exit(0)
	}
	if *example {
		printExample()
		os.Exit(0)
	}
	cfg.BindPFlag("listen.api", pflag.Lookup("listen-api"))
	cfg.BindPFlag("listen.publish", pflag.Lookup("listen-publish"))
	cfg.BindPFlag("api.token", pflag.Lookup("api-token"))
	cfg.BindPFlag("api.disable", pflag.Lookup("api-disable"))
	cfg.BindPFlag("publish.IgnoreHeaders", pflag.Lookup("ignore-headers"))
	cfg.BindPFlag("publish.RealIpHeader", pflag.Lookup("realip-header"))
	cfg.BindPFlag("consume.FailWait", pflag.Lookup("fail-wait"))
	cfg.BindPFlag("consume.GoFailWait", pflag.Lookup("go-fail-wait"))
	cfg.BindPFlag("consume.DataFile", pflag.Lookup("data-file"))
	cfg.BindPFlag("rabbitmq.host", pflag.Lookup("mq-host"))
	cfg.BindPFlag("rabbitmq.port", pflag.Lookup("mq-port"))
	cfg.BindPFlag("rabbitmq.username", pflag.Lookup("mq-username"))
	cfg.BindPFlag("rabbitmq.password", pflag.Lookup("mq-password"))
	cfg.BindPFlag("rabbitmq.vhost", pflag.Lookup("mq-vhost"))
	cfg.BindPFlag("rabbitmq.prefix", pflag.Lookup("mq-prefix"))
	cfg.BindPFlag("log.dir", pflag.Lookup("log-dir"))
	cfg.BindPFlag("log.level", pflag.Lookup("log-level"))
	cfg.BindPFlag("log.access", pflag.Lookup("log-access"))
	cfg.BindPFlag("log.post", pflag.Lookup("log-post"))
	cfg.BindPFlag("log.console-level", pflag.Lookup("level"))
	cfg.BindPFlag("log.fileMaxSize", pflag.Lookup("log-max-size"))
	cfg.BindPFlag("log.maxCount", pflag.Lookup("log-max-count"))
	cfg.SetDefault("default.IgnoreHeaders", []string{"Token", "RouteKey", "Host", "Expect", "Accept-Encoding", "Content-Length", "Connection"})
	fmt.Printf("%s", *configFile)
	if *configFile != "" {
		cfg.SetConfigFile(*configFile)
	} else {
		cfg.SetConfigName("config")
		cfg.AddConfigPath("/etc/wmq/")
		cfg.AddConfigPath("$HOME/.wmq")
		cfg.AddConfigPath("../wmq")
		cfg.AddConfigPath("../")
	}
	err := cfg.ReadInConfig()
	file := cfg.ConfigFileUsed()
	if err != nil && !strings.Contains(err.Error(), "Not") {
		fmt.Printf("%s", err)
	} else if file != "" {
		fmt.Printf("use config file : %s\n", file)
	}
	err = nil
	cfg.Set("publish.IgnoreHeaders", append(cfg.GetStringSlice("default.IgnoreHeaders"), cfg.GetStringSlice("publish.IgnoreHeaders")...))
}

// load message file: read message.json file to qmessage
func (container *Container) LoadMessageFile() error {
	messageFile := container.Config.GetString("consume.Datafile")
	file, err := os.OpenFile(messageFile, os.O_RDONLY|os.O_CREATE, 0766)
	if err != nil {
		return err
	}
	defer file.Close()

	err = container.QMessage.ReadFrom(file)
	if err != nil {
		return err
	}
	return nil
}

// reset message file: reset messages write message.json file
func (container *Container) ResetMessageFile() error {
	messageFile := container.Config.GetString("consume.Datafile")
	file, err := os.OpenFile(messageFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0766)
	if err != nil {
		return err
	}
	defer file.Close()

	err = container.QMessage.WriteTo(file, true)
	if err != nil {
		return err
	}
	return nil
}

// init RabbitMQ pools
func (container *Container) InitRabbitMQPools() {
	container.RabbitMQPools.Init(20)
}

// init RabbitMQ all exchanges
func (container *Container) InitExchanges() error {
	rabbitMq, err := container.RabbitMQPools.GetMQ()
	if err != nil {
		return errors.New("rabbitmq pools faild: " + err.Error())
	}
	defer container.RabbitMQPools.Recover(rabbitMq)

	container.QMessage.Lock.Lock()
	defer container.QMessage.Lock.Unlock()

	messages := container.QMessage.Messages
	for _, message := range messages {
		// declare exchange
		err := rabbitMq.DeclareExchange(message.Name, message.Mode, message.Durable)
		if err != nil {
			return errors.New("Declare exchange faild: "+err.Error())
		}
		// declare queue
		for _, consumer := range message.Consumers {
			consumerKey := message.Name +"-"+ consumer.ID
			err := rabbitMq.DeclareQueue(consumerKey, message.Durable)
			if err != nil {
				return errors.New("Declare queue faild: "+err.Error())
			}
			// bind queue to exchange
			err = rabbitMq.BindQueueToExchange(consumerKey, message.Name, consumer.RouteKey)
			if err != nil {
				return errors.New("bind queue exchange fail: "+err.Error())
			}

			container.SendConsumerSign("insert", consumerKey)
		}
	}
	return nil
}

// declare a Exchange
func (container *Container) DeclareExchange(name string) error {
	rabbitMq, err := container.RabbitMQPools.GetMQ()
	if err != nil {
		return errors.New("rabbitmq pools faild: " + err.Error())
	}
	defer container.RabbitMQPools.Recover(rabbitMq)

	message, err := container.QMessage.GetMessageByName(name)
	if err != nil {
		return err
	}
	// declare exchange
	err = rabbitMq.DeclareExchange(message.Name, message.Mode, message.Durable)
	if err != nil {
		return errors.New("Declare exchange "+message.Name+" faild: "+err.Error())
	}

	return nil
}

// delete a exchange
func (container *Container) DeleteExchange(name string) error {
	rabbitMq, err := container.RabbitMQPools.GetMQ()
	if err != nil {
		return errors.New("rabbitmq pools faild: " + err.Error())
	}
	defer container.RabbitMQPools.Recover(rabbitMq)

	consumers := container.QMessage.GetConsumersByMessageName(name)
	if len(consumers) > 0 {
		//unbind queue exchange
		for _, consumer := range consumers {
			consumerKey := name+"-"+consumer.ID
			routeKey := consumer.RouteKey
			err := rabbitMq.UnBindQueueToExchange(consumerKey, name, routeKey)
			if err != nil {
				return errors.New("Unbind exchange "+name+" consumer id "+consumer.ID+" faild: "+err.Error())
			}
		}
	}
	// delete exchange
	err = rabbitMq.DeleteExchange(name)
	if err != nil {
		return errors.New("Delete exchange "+name+" faild: "+err.Error())
	}

	return nil
}

// declare a consumer
func (container *Container) DeclareConsumer(name string, consumerId string) error {
	rabbitMq, err := container.RabbitMQPools.GetMQ()
	if err != nil {
		return errors.New("rabbitmq pools faild: " + err.Error())
	}
	defer container.RabbitMQPools.Recover(rabbitMq)

	message, err := container.QMessage.GetMessageByName(name)
	if err != nil {
		return err
	}
	consumer, err := container.QMessage.GetConsumerById(name, consumerId)
	if err != nil {
		return err
	}
	consumerKey := message.Name +"-"+ consumer.ID
	// declare queue
	err = rabbitMq.DeclareQueue(consumerKey, message.Durable)
	if err != nil {
		return errors.New("Declare queue faild: "+err.Error())
	}
	// bind queue to exchange
	err = rabbitMq.BindQueueToExchange(consumerKey, message.Name, consumer.RouteKey)
	if err != nil {
		return errors.New("bind queue exchange fail: "+err.Error())
	}

	return nil
}

// init consumer
func (container *Container) InitConsumer() {

	go func() {
		consumerWorker := <- container.ConsumerWorkChan

		switch consumerWorker.Action {
		case "insert":
			go func() {
				rabbitMq, _ := container.RabbitMQPools.GetMQ()
				defer container.RabbitMQPools.Recover(rabbitMq)
				channel, _ := rabbitMq.Conn.Channel()
				defer channel.Close()

				autoAck := false // is auto ack
				exclusive := false //
				noLocal := false
				noWait := false
				delivery, _ := channel.Consume(consumerWorker.QueueName, "", autoAck, exclusive, noLocal, noWait, nil)
				for {
					select {
					case d := <-delivery:
						file, _ := os.OpenFile("./consumer.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
						file.WriteString(string(d.Body) + "\r\n")
						file.Close()
						d.Ack(true)
					//d.Nack(false, true)
					}
				}
			}()
		case "update":
			return
		case "delete":
			return
		case "status":
			return
		}
	}()
}

// publish message to mq
func (container *Container) Publish(body string, exchangeName string, token string, routeKey string) (err error) {
	qMessage, err := container.QMessage.GetMessageByName(exchangeName)
	if err != nil {
		return
	}
	if qMessage.IsNeedToken && qMessage.Token != token {
		return errors.New("token error")
	}

	rabbitMq, err := container.RabbitMQPools.GetMQ()
	if err != nil {
		return errors.New("rabbitmq pools faild: "+err.Error())
	}
	defer container.RabbitMQPools.Recover(rabbitMq)

	return rabbitMq.Publish(exchangeName, routeKey, body)
}

// send consumer sign
func (container *Container) SendConsumerSign(action string, consumerKey string) {
	container.ConsumerWorkChan <- ConsumerWorker{
		Action: action,
		QueueName: consumerKey,
	}
}

func printExample() {
	fmt.Println(`[{
    "Comment": "",
    "Consumers": [{
            "Comment": "",
            "ID": "111",
            "Code": 200,
            "CheckCode": true,
            "RouteKey": "#",
            "Timeout": 5000,
            "URL": "http://test.com/wmq.php"
        }
    ],
    "Durable": false,
    "IsNeedToken": true,
    "Mode": "topic",
    "Name": "test",
    "Token": "JQJsUOqYzYZZgn8gUvs7sIinrJ0tDD8J"
}]`)
}

func poster() string {
	fg := color.New(color.FgHiYellow).SprintFunc()
	return fg(`
██╗    ██╗███╗   ███╗ ██████╗
██║    ██║████╗ ████║██╔═══██╗
██║ █╗ ██║██╔████╔██║██║   ██║
██║███╗██║██║╚██╔╝██║██║▄▄ ██║
╚███╔███╔╝██║ ╚═╝ ██║╚██████╔╝
 ╚══╝╚══╝ ╚═╝     ╚═╝ ╚══▀▀═╝
Author: snail
Link  : https://github.com/snail007/wmq
`)
}