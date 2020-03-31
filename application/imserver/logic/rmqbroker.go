package logic

import (
	"fmt"
	"github.com/micro/go-micro/v2/broker"
	"log"
)

type (
	RabbitMqBroker struct {
		topic          string
		rabbitMqBroker broker.Broker
	}
)

func NewRMqBroker (topic string, rmqBroker broker.Broker)(*RabbitMqBroker ,error){
	if err := rmqBroker.Init();err != nil {
		return  nil ,err
	}
	if err := rmqBroker.Connect();err != nil {
		return nil,err
	}
	return &RabbitMqBroker{topic:topic,rabbitMqBroker:rmqBroker},nil
}

func (r *RabbitMqBroker) Publisher (msg *broker.Message) error{
	if err := r.rabbitMqBroker.Publish(r.topic, msg); err != nil {
		log.Printf("[publisher %s err] : %+v", r.topic, err)
		return  err
	}
	log.Printf("[publisher %s] : %s", r.topic, string(msg.Body))

	return nil
}

func (r *RabbitMqBroker) Subscribe (senMsg func(topic string,body []byte)) error {
	_,err := r.rabbitMqBroker.Subscribe(r.topic, func(event broker.Event) error {
		fmt.Println(event.Topic(), string(event.Message().Body))
		senMsg(event.Topic(),event.Message().Body)
		return nil
	})
	if err != nil {
		fmt.Println(222)
		return err
	}
	return nil
}