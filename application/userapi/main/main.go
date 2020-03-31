package main

import (
	"awesomeProject/application/userapi/logic"
	"awesomeProject/application/userapi/middleware"
	users_model "awesomeProject/application/userapi/mod"
	user_config "awesomeProject/application/userserver/main/userconfig"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/micro/go-micro/v2/web"
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
	//etcds := etcdv3.NewRegistry(
	//	registry.Addrs(conf.ETCDConfig.Addrs),
	//)
	service := web.NewService(
		web.Name(conf.ServerName),
		//web.Registry(etcds),
		web.Address(":8080"),
	)
	r := gin.Default()
	userGroup := r.Group("/user")
	userGroup.POST("/login", logic.Login)
	userGroup.POST("/registry",logic.Register )
	userGroup.Use(middleware.ValidTokenMiddleware)

	service.Handle("/", r)

	// Init will parse the command line flags.
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}
	// Run the server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
