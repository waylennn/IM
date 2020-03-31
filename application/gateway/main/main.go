package main

import (
	"awesomeProject/application/gateway/logic"
	"awesomeProject/application/gateway/main/config"
	gateway_model "awesomeProject/application/gateway/mod"
	im "awesomeProject/application/imserver/protoc"
	user "awesomeProject/application/userserver/protoc"

	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/transport/grpc"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/micro/go-micro/v2/web"
	"github.com/micro/go-plugins/registry/etcdv3"
)

func main() {
	//初始化配置文件
	conf, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(conf)
	//初始化DB
	err = gateway_model.InitModEngine(conf.DB.Name,conf.DB.Address)
	if err != nil {
		log.Fatal(err)
	}
	//选择etcd作为注册中心
	etcds := etcdv3.NewRegistry(
		registry.Addrs(conf.ETCDConfig.Addrs),
	)
	service := micro.NewService(
		micro.Name("client"),
		micro.Registry(etcds),
		micro.Transport(grpc.NewTransport()),
	)
	service.Init()

	logic.UserClient = user.NewUserService(conf.UserRpcServer.ServerName,service.Client())
	logic.ImClient = im.NewImService(conf.ImRpcServer.ServerName,service.Client())
	logic.NewAddressList(conf.ImRpcServer.ImServerList)

	useGate := web.NewService(
		web.Name(conf.ServerName),
		web.Registry(etcds),
		web.Address(conf.Web.Address),
	)
	router := gin.Default()
	userRouterGroup := router.Group("/gateway")
	//userRouterGroup.Use(middleware.ValidateTokenMiddleware)
	{
		userRouterGroup.POST("/send", logic.Send)
		userRouterGroup.POST("/address", logic.GetServerAddress)
	}
	useGate.Handle("/", router)
	if err := useGate.Init(); err != nil {
		log.Fatal(err)
	}
	if err := useGate.Run(); err != nil {
		log.Fatal(err)
	}
}
