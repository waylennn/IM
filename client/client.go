package main

import (
	hello "awesomeProject/protos"
	"context"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	micro "github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/transport/grpc"
	"github.com/micro/go-plugins/registry/etcdv3"
	"time"
)

type clientWrapper2 struct {
	client.Client
}

func (c *clientWrapper2) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	return hystrix.Do(req.Service()+"."+req.Endpoint(), func() error {
		return c.Client.Call(ctx, req, rsp, opts...)
	}, nil)
}

// NewClientWrapper returns a hystrix client Wrapper.
func NewClientWrapper2() client.Wrapper {
	return func(c client.Client) client.Client {
		return &clientWrapper2{c}
	}
}

func main() {
	etcds := etcdv3.NewRegistry(
		registry.Addrs("127.0.0.1:2379"),
	)

	// Create a new service
	service := micro.NewService(
		micro.Name("greeter.client"),
		micro.Registry(etcds),
		micro.Transport(grpc.NewTransport()),
		micro.WrapClient(NewClientWrapper2()),
	)
	// Initialise the client and parse command line flags
	service.Init()

	t := time.NewTicker(time.Millisecond * 200)
	for _ = range t.C {
		// Create new greeter client
		greeter := hello.NewGreeterService("greeter", service.Client())
		// Call the greeter
		rsp, err := greeter.Hello(context.TODO(), &hello.Request{Name: "John"})
		if err != nil {
			fmt.Println(err)
			continue
		}
		// Print response
		fmt.Println(rsp.Msg)
	}

}
