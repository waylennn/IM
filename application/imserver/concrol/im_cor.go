package concrol

import (
	"awesomeProject/application/imserver/logic"
	im "awesomeProject/application/imserver/protoc"
	"context"
	"encoding/json"
	"fmt"
	"github.com/micro/go-micro/v2/broker"
)

type ImServer struct {
	RabbitMq  map[string]*logic.RabbitMqBroker
}

func (i *ImServer) PublishMessage(ctx context.Context, in *im.PublishMessageRequest, out *im.PublishMessageResponse) error {
	body,err := json.Marshal(in);
	if err != nil {
		return err
	}

	rmq := i.RabbitMq[in.ServerName+in.Topic]
	if err := rmq.Publisher(&broker.Message{Body:body}) ; err != nil {
		return err
	}
	return nil
}
