[![logo](./logo.png)](https://github.com/phachon/wmqx)

[![stable](https://img.shields.io/badge/stable-stable-green.svg)](https://github.com/phachon/wmqx/) 
[![build](https://img.shields.io/shippable/5444c5ecb904a4b21567b0ff.svg)](https://travis-ci.org/phachon/wmqx)
[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/phachon/wmqx)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/phachon/wmqx/master/LICENSE)
[![go_Report](https://goreportcard.com/badge/github.com/phachon/wmqx)](https://goreportcard.com/report/github.com/phachon/wmqx)
[![powered_by](https://img.shields.io/badge/powered_by-Go-3362c2.svg?style=flat)]()
[![platforms](https://img.shields.io/badge/platform-All-yellow.svg?style=flat)]()
[![download_count](https://img.shields.io/github/downloads/phachon/wmqx/total.svg?style=plastic)](https://github.com/phachon/wmqx/releases) 
[![release](https://img.shields.io/github/release/phachon/wmqx.svg?style=flat)](https://github.com/phachon/wmqx/releases) 

WMQX 是一个基于 RabbitMQ 开发的支持 http 协议的 MQ 服务, 他的前身是 [wmq](https://github.com/snail007/wmq), 感谢他的作者同时也是我的好友 [snail007](https://github.com/snail007) , 当然，你也可以理解为 WMQX 是 wmq 的升级版。

# 为什么需要 WMQX?

## RabbitMQ 的使用
RabbitMQ 是一个轻量级的，易于部署在本地和云上，支持多个消息传递协议的消息队列中间件。RabbitMQ 可以应用于许多的场景中，同时也支持多种语言的 SDK。通常你会这样使用 RabbitMQ：
 
1. 作为生产者：
    - 编写代码连接到 RabbitMQ，并打开一个 channel。
    - 编写代码声明一个 exchange，并设置相关属性。
    - 编写代码声明一个 queue，并设置相关属性。
    - 编写代码使用 routing key，在 exchange 和 queue 之间建立好绑定关系。
    - 编写代码发送消息到 RabbitMQ。
2. 作为消费者：
    - 编写代码连接到 RabbitMQ，并打开一个 channel, 开启消费进程，等待接收到消息后，处理消费的业务逻辑。
    
如下图所示：
[![RabbitMQ](./docs/images/rabbitmq.png)](https://github.com/phachon/wmqx)

### 问题：
1. RabbitMQ 的连接、Exchange、Queue 的声明和修改和业务代码耦合在一起，增加了开发和维护的成本。
2. 当修改消费者的业务逻辑，可能会需要频繁的重启消费进程。
3. 对于第一次使用 MQ 的用户，去理解 RabbitMQ 的原理和编写代码实现生产和消费是需要一定的时间和人力成本。
4. 。。。

### 解决与实现：
1. 将 RabbitMQ 的连接、Exchange、Queue 的声明和修改删除等一些和业务无关的操作抽离出来单独提供服务，Exchange、queue 的操作以友好的 API 的方式提供给用户。
2. 帮助用户去实现每一个消息的消费进程，用户只需要提供消费者的 API 接口，消费进程等待有消息后调用对应的消费者 API。消费者业务逻辑修改，只需要修改 API, 消费进程无需重启。
3. 对于第一次使用 MQ 或者不清楚 RabbitMQ 原理的用户，不需要去深入了解 RabbitMQ 的使用和编码实现，只需要通过 http 的方式接入服务，即可快速使用消息队列。 

* 所以 WMQX 也就由此诞生。工作原理如下图所示：*
[![wmqx](./docs/images/wmqx.png)](https://github.com/phachon/wmqx)