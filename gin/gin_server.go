package main

import (
	serverConfig "awesomeProject/server/config"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/web"
	"time"

	"github.com/micro/go-plugins/registry/etcdv3"
)

type Greeter struct{}

func main() {
	//初始化配置文件
	conf := initConfig()
	//选择etcd作为注册中心
	etcds := etcdv3.NewRegistry(
		registry.Addrs(conf.ETCDConfig.Addrs),
	)
	// Create a new service. Optionally include some options here.
	service := web.NewService(
		web.Name(conf.ServerName),
		web.Registry(etcds),
		web.Address(":8811"),
	)
	r := gin.Default()
	r.GET("/user/login", Login)

	userGroup := r.Group("/user")
	userGroup.Use(ValidTokenMiddleware)
	userGroup.GET("/list", List)

	service.Handle("/", r)

	if err := service.Init(); err != nil {
		fmt.Println(err)
	}

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

func generateJWT() (tokenStr string, err error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(30 * time.Second).Unix(),
		Issuer:    "test",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if tokenStr, err = token.SignedString([]byte("jey.sign")); err != nil {
		return
	}
	return
}

func validToken(tokenStr string) error {
	_, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte("jey.sign"), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				fmt.Println("Timing is everything")
				return err
			} else {
				fmt.Println("token invalid:", err)
				return err
			}
		}
	}
	return nil
}

func Login(c *gin.Context) {
	token, err := generateJWT()
	fmt.Println(token, err)
	if err != nil {
		c.JSON(500, err)
		return
	}
	username := c.Query("username")
	password := c.Query("username")

	c.JSON(200, gin.H{
		"token":   token,
		"message": username,
		"nick":    password,
	})
}

func List(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": 200,
	})
}

//验证TOKEN中间件
func ValidTokenMiddleware(c *gin.Context) {
	tokenStr := c.GetHeader("Auth")
	err := validToken(tokenStr)
	if err != nil {
		c.JSON(429, &gin.H{"message": err})
		c.Abort()
	}
	c.Next()
}
