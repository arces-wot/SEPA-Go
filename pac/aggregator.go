package pac

import (
	"github.com/arces-wot/SEPA-Go/sepa"
	"github.com/arces-wot/SEPA-Go/sepa/sparql"
)

type aggregator struct {
	c Consumer
	p Producer
}

func (a aggregator)Consume(handler func(notification *sparql.Notification),data interface{})(sepa.Subscription,error)  {
	return a.c.Consume(handler,data)
}

func (a aggregator)Produce(data interface{}) error {
	return a.p.Produce(data)
}