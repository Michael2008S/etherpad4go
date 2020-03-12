package poker

import "github.com/y0ssar1an/q"
import "github.com/thedevsaddam/gojsonq"

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:

			// 消息判断分发处理
			q.Q("Run_broadcast:", string(message))

			handleMessage(message)
			responseMsg := []byte("{}")

			for client := range h.clients {
				select {
				case client.send <- responseMsg:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func handleMessage(message []byte) string {
	msgType := gojsonq.New().FromString(string(message)).Find("type")
	//
	//accessStatus =  securityManager.checkAccess(padID, sessionCookie, token, password)

	//if (accessStatus !== "grant") {
	//
	//}

	if msgType.(string) == "CLIENT_READY" {
		//	handleClientReady(client, message);
	} else if msgType.(string) == "CHANGESET_REQ" {
		//	handleChangesetRequest(client, message);
	}

	return ""
}

func handleClientReady() {

}

func handleChangesetRequest() {

}
