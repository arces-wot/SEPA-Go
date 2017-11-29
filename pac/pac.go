package pac

import (
	"errors"
	"github.com/arces-wot/SEPA-Go/sepa"
	"github.com/arces-wot/SEPA-Go/sepa/sparql"
)

type Producer interface {
	Produce(data interface{}) error
}

type Consumer interface {
	Consume(handler func(notification *sparql.Notification), data interface{}) (sepa.Subscription, error)
}

type Aggregator interface {
	Producer
	Consumer
}

type Application struct {
	client  sepa.Client
	profile Profile
}

const EMPTYDATA  = ""

func newApplication(profile Profile) Application {
	client := sepa.NewClient(profile.Configuration)
	return Application{client, profile}
}

func (a Application) newProducer(updateID string) (Producer, error) {
	if !a.profile.ContainsUpdate(updateID) {
		return producer{}, errors.New("no update query found in profile")
	}
	return producer{updateID, a}, nil
}

func (a Application) newAggrator(updateId string, queryId string) (Aggregator, error) {
	agg := aggregator{}
	var err error
	if agg.p, err = a.newProducer(updateId); err != nil {
		return agg, err
	}
	if agg.c, err = a.newConsumer(queryId); err != nil {
		return agg, err
	}
	return agg, nil
}

func (a Application) newConsumer(queryID string) (Consumer, error) {
	if !a.profile.ContainsQuery(queryID) {
		return consumer{}, errors.New("no query found in profile")
	}
	return consumer{queryID, a}, nil
}