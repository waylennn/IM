package main

import (
	"awesomeProject/application/imserver/concrol"
	"awesomeProject/application/imserver/logic"
	"awesomeProject/application/imserver/main/config"
	"awesomeProject/application/imserver/protoc"
	"github.com/micro/go-plugins/broker/rabbitmq/v2"

	"fmt"
	"github.com/juju/ratelimit"
	micro "github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/transport/grpc"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/micro/go-plugins/registry/etcdv3"
	ratelimit2 "github.com/micro/go-plugins/wrapper/ratelimiter/ratelimit/v2"
)

func main() {
	//初始化配置文件
	conf, err := config.InitConfig("../config/config.json")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(conf)

	//选择etcd作为注册中心
	etcds := etcdv3.NewRegistry(
		registry.Addrs(conf.ETCDConfig.Addrs),
	)
	//限流
	rl := ratelimit.NewBucketWithRate(float64(conf.RateLimit.Rate), int64(conf.RateLimit.Rate))
	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name(conf.ServerName),
		micro.Registry(etcds),
		micro.Transport(grpc.NewTransport()),
		micro.WrapHandler(ratelimit2.NewHandlerWrapper(rl, false)),
	)

	// Init will parse the command line flags.
	service.Init()
	//不同的topic和服务名 对应不同的MQ 这里MQ可以集群部署 也可以单台
	publisherServerMap := make(map[string]*logic.RabbitMqBroker)
	for _, item := range conf.ImServerList {
		//amqbAddress := item.AmqbAddress
		//p, err := server.NewRabbitMqBroker(
		//	item.Topic,
		//	rabbitmq.NewBroker(func(options *broker.Options) {
		//		options.Addrs = amqbAddress
		//	}),
		//)
		//if err != nil {
		//	log.Fatal(err)
		//}

		//这里没有MQ就使用默认的
		rmqBroker,err := logic.NewRMqBroker(item.Topic,rabbitmq.NewBroker())
		if err != nil {
				log.Fatal(err)
			}
		publisherServerMap[item.ServerName+item.Topic] = rmqBroker
	}

	imServer := &concrol.ImServer{
		RabbitMq:publisherServerMap,
	}
	if err := im.RegisterImHandler(service.Server(), imServer); err != nil {
		log.Fatal(err)
	}
	// Run the server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
