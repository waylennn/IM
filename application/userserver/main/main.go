package main

import (
	user_cro "awesomeProject/application/userserver/concrol"
	user_config "awesomeProject/application/userserver/main/userconfig"
	users_model "awesomeProject/application/userserver/mod"
	user_pro "awesomeProject/application/userserver/protoc"
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
	conf, err := user_config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(conf)
	//初始化DB
	err = users_model.InitModEngine(conf.DB.Name,conf.DB.Address)
	if err != nil {
		log.Fatal(err)
	}
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

	// Register handler
	users := &user_cro.User{}

	if err := user_pro.RegisterUserHandler(service.Server(), users); err != nil {
		log.Fatal(err)
	}
	// Run the server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
