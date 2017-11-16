package sepa

import (
	"github.com/gorilla/websocket"
	"sync"
	"github.com/arces-wot/SEPA-Go/sepa/sparql"
)



type Subscription struct {
	Id string
	client *Client
	unserHandler func(*sparql.Notification)
	connection *websocket.Conn
	unsublock *sync.Mutex
}

func (s Subscription) Unsubscribe()  {
	s.client.unsubscribe(s)
}









