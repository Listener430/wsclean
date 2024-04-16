package websocket

import (
	"github.com/gorilla/websocket"
	"sync"
)

var NotificationPools = make(map[string]NotificationClient)

type NotificationClient struct {
	Client *websocket.Conn
}

type OnlineHub struct {
	Register   chan *OnlineClient
	Unregister chan *OnlineClient
	Users      map[string]*OnlineClient
	Quit       chan struct{}
	Wg         sync.WaitGroup
}

type OnlineClient struct {
	OnlineHub           *OnlineHub
	WebSocketConnection *websocket.Conn
	Send                chan OnlineEventStructure
	UserId              string
}

type OnlineEventStructure struct {
	EventType string `json:"eventType"`
}
