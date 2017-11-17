package sepa

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/arces-wot/SEPA-Go/sepa/sparql"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type Client struct {
	host           string
	query_path     string
	update_path    string
	subscribe_path string
}

type subscribeRequest struct {
	Subscribe string `json:"subscribe"`
	Alias     string `json:"alias"`
}

type unsubscribeRequest struct {
	Unsubscribe string `json:"unsubscribe"`
}

type sepaError struct {
	Body string `json:"body"`
	Code int    `json:"code"`
}

func (err sepaError) Error() string {
	return err.Body
}

type subscribeSuccessResponse struct {
	Subscribed string `json:"subscribed"`
	Alias      string `json:"alias"`
}

func NewClient(config Configuration) Client {
	return Client{
		host:           config.Host,
		query_path:     fmt.Sprintf("%s:%d%s", config.Host, config.Ports.Http, config.Paths.Query),
		update_path:    fmt.Sprintf("%s:%d%s", config.Host, config.Ports.Http, config.Paths.Update),
		subscribe_path: fmt.Sprintf("%s:%d%s", config.Host, config.Ports.Ws, config.Paths.Subscribe),
	}
}

/*
	Create a client connected to  a sepa engine with default settings:
	"ports" : {
			"http" : 8000 ,
			"ws" : 9000 ,
		}
		 ,
	"paths" : {
			"update" : "/update" ,
			"query" : "/query" ,
			"subscribe" : "/subscribe"
		}
*/
func NewDefaultClient(host string) Client {
	config := Configuration{
		Host:  host,
		Ports: PortsType{Http: 8000, Ws: 9000},
		Paths: PathsType{Query: "/query", Update: "/update", Subscribe: "/subscribe"},
	}
	return NewClient(config)
}

func (c *Client) Query(sparqlQuery string) (*sparql.Results, error) {
	body := strings.NewReader(sparqlQuery)
	resp, err := http.Post("http://"+c.query_path, "application/sparql-query", body)

	if err != nil {
		return nil, err
	}

	res, err := sparql.ParseFromJson(resp.Body)
	return res, err
}

func (c *Client) Update(sparqlUpdate string) error {
	body := strings.NewReader(sparqlUpdate)
	_, err := http.Post("http://"+c.update_path, "application/sparql-update", body)
	return err
}

func (c *Client) Subscribe(sparqlQuery string, handler func(notification *sparql.Notification)) (Subscription, error) {
	header := http.Header{}
	header.Set("Origin", c.host)
	conn, _, err := websocket.DefaultDialer.Dial("ws://"+c.subscribe_path, header)
	fail_sub := Subscription{}
	if err != nil {
		return fail_sub, err
	}

	req := subscribeRequest{sparqlQuery, ""}
	conn.WriteJSON(req)

	if _, r, err := conn.NextReader(); err == nil {
		resp, err := decodeSubscribeResponse(r)
		if err != nil {
			return fail_sub, err
		}
		var waitForUnsubMessage sync.Mutex
		waitForUnsubMessage.Lock()
		subscription := Subscription{
			client:       c,
			unserHandler: handler,
			connection:   conn,
			Id:           resp.Subscribed,
			unsublock:    &waitForUnsubMessage}

		go notificationReader(conn, subscription)
		return subscription, nil
	}

	return fail_sub, err
}

func (c *Client) unsubscribe(subscription Subscription) error {
	request := unsubscribeRequest{Unsubscribe: subscription.Id}

	err := subscription.connection.WriteJSON(request)

	if err == nil {
		subscription.unsublock.Lock()
		subscription.unsublock.Unlock() //free the lock
	}

	return err
}

func decodeSubscribeResponse(r io.Reader) (subscribeSuccessResponse, error) {
	dec := json.NewDecoder(r)
	resp := subscribeSuccessResponse{"", ""}
	var v map[string]interface{}

	if err := dec.Decode(&v); err != nil {
		return resp, errors.New("Can't decode json")
	}

	if v["subscribed"] == nil {
		strcod, _ := v["Code"].(string)
		code, _ := strconv.Atoi(strcod)
		return resp, sepaError{Body: v["body"].(string), Code: code}
	}

	resp.Subscribed = v["subscribed"].(string)
	resp.Alias = v["alias"].(string)

	return resp, nil
}

func notificationReader(ws_conn *websocket.Conn, sub Subscription) {
	listening := true
	for listening {
		messageType, message, err := ws_conn.ReadMessage()

		if err != nil {
			//TODO: An error from ReadMessage/ nextReader is permanent and every other read will give the same error
			return
		}

		switch messageType {
		case websocket.TextMessage:
			smessage := string(message)

			if ok, _ := regexp.MatchString("^{ *\"results\" *:", smessage); ok {
				reader := strings.NewReader(smessage)

				if not, parse_error := sparql.ParseNotificationJson(reader); parse_error == nil {
					sub.unserHandler(not)
				}

			} else if ok, _ := regexp.MatchString("^{ *\"ping\" *:", smessage); ok {
				fmt.Println("Ping")
			} else if ok, _ := regexp.MatchString("^{ *\"unsubscribed\" *:", smessage); ok {
				//TODO: handle unsubcribe

				//and exit
				//ws_conn.Close()
				listening = false
			}

		}
	}
	//Now the unsuscribe process is complete
	// and the go routine is garbage collected
	sub.unsublock.Unlock()
	return
}
