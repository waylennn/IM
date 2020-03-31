package main

import (
	hello "awesomeProject/protos"
	serverConfig "awesomeProject/server/config"
	"context"
	"fmt"
	"github.com/juju/ratelimit"
	micro "github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/transport/grpc"

	"github.com/micro/go-plugins/registry/etcdv3"
	ratelimit2 "github.com/micro/go-plugins/wrapper/ratelimiter/ratelimit/v2"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	return nil
}

func main() {
	if err := broker.Init(); err != nil {
		fmt.Println(111)
		fmt.Println(err)
		return
	}
	if err := broker.Connect(); err != nil {
		fmt.Println(err)
		return
	}
	go publish()
	go subscribe()

	//初始化配置文件
	conf := initConfig()
	//选择etcd作为注册中心
	etcds := etcdv3.NewRegistry(
		registry.Addrs(conf.ETCDConfig.Addrs),
	)
	//限流
	br := ratelimit.NewBucketWithRate(float64(2), int64(2))
	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name(conf.ServerName),
		micro.Registry(etcds),
		micro.Transport(grpc.NewTransport()),
		micro.WrapHandler(ratelimit2.NewHandlerWrapper(br, false)),
		//micro.Broker(rb),
	)

	// Init will parse the command line flags.
	service.Init()

	// Register handler
	hello.RegisterGreeterHandler(service.Server(), new(Greeter))

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

func initConfig() *serverConfig.ServerConfig {
	//初始化配置文件
	conf := serverConfig.NewServerConfig()
	if err := config.LoadFile("./userconfig/userconfig.json"); err != nil {
		fmt.Println(err)
		panic(err)
	}
	if err := config.Scan(conf); err != nil {
		fmt.Println(err)
		panic(err)
	}

	return conf
}

func publish() {
	a := []byte("dd")
	err := broker.Publish("demo.broker", &broker.Message{Header: map[string]string{"go": "gogo"}, Body: a})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func subscribe() {
	subscribe, err := broker.Subscribe("demo.broker", func(event broker.Event) error {
		fmt.Println(event.Topic(), string(event.Message().Body))
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println()
	fmt.Printf("%v", subscribe)
}
