package main

import (
	"awesomeProject/application/imserver/logic"
	"awesomeProject/application/imserver/main/config"
	"awesomeProject/application/imserver/mod"
	"fmt"
	"github.com/juju/ratelimit"
	micro "github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/transport/grpc"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/micro/go-plugins/broker/rabbitmq/v2"
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
	//etcds := etcdv3.NewRegistry(
	//	registry.Addrs(conf.ETCDConfig.Addrs),
	//)
	//限流
	rl := ratelimit.NewBucketWithRate(float64(conf.RateLimit.Rate), int64(conf.RateLimit.Rate))
	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name(conf.ServerName),
		//micro.Registry(etcds),
		micro.Transport(grpc.NewTransport()),
		micro.WrapHandler(ratelimit2.NewHandlerWrapper(rl, false)),
	)

	// Init will parse the command line flags.
	service.Init()
	err = mod.InitModEngine(conf.DB.Name,conf.DB.Address)
	if err != nil {
		log.Fatal(err)
	}

	rmqBrokerRegistry := rabbitmq.NewBroker()
	rmqBroker, err := logic.NewRMqBroker(conf.RabbitMq.Topic, rmqBrokerRegistry)
	if err != nil {
		log.Fatal(err)
	}

	imServer,err := logic.NewImServer(rmqBroker, logic.ImServerAddress(":7073"))
	if err != nil {
		log.Fatal(err)
	}

	go imServer.Subscribe()
	go imServer.Run()

	// Run the server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
