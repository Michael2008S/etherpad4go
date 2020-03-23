package poker

import (
	"encoding/json"
	"github.com/Michael2008S/etherpad4go/api"
	bgStore "github.com/Michael2008S/etherpad4go/store"
	"github.com/y0ssar1an/q"
	"log"
)
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

	// 存储层
	dbStore bgStore.Store
}

func NewHub(db bgStore.Store) *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		dbStore:    db,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			//q.Q("Hub Run client<-register", client)
			log.Println("Hub Run client<-register:%+v", client)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			log.Println("Hub Run client<-unregister", client)
		case message := <-h.broadcast:

			// 消息判断分发处理
			q.Q("Run_broadcast:", string(message))

			responseMsg := h.handleMessage(message)
			q.Q("Send_broadcast:", string(responseMsg))

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

func (h *Hub) handleMessage(message []byte) []byte {
	msgType := gojsonq.New().FromString(string(message)).Find("type")
	//
	//accessStatus =  securityManager.checkAccess(padID, sessionCookie, token, password)

	//if (accessStatus !== "grant") {
	//
	//}

	if msgType.(string) == "CLIENT_READY" {
		//	handleClientReady(client, message);

		//createSessionInfo()
		q.Q("createSessionInfo")
	} else if msgType.(string) == "CHANGESET_REQ" {
		//	handleChangesetRequest(client, message);
	} else if msgType.(string) == "COLLABROOM" {
		msgDataType := gojsonq.New().FromString(string(message)).Find("type")
		if msgDataType == "USER_CHANGES" {
			// TODO padChannels.emit(message.padId, {client: client, message: message}); // add to pad queue
			handleUserChanges(h.dbStore, message)
		}
	}

	return []byte(`{"type":"error"}`)
}

func handleClientReady() {

}

func handleChangesetRequest() {

}

func handleUserChanges(db bgStore.Store, message []byte) []byte {
	reqMsg := api.CollabRoomReqMessage{}
	json.Unmarshal(message, &reqMsg)

	// get all Vars we need
	//baseRev := reqMsg.Data.BaseRev
	//wireApool := reqMsg.Data.Apool
	//changeset := reqMsg.Data.Changeset

	// The client might disconnect between our callbacks. We should still
	// finish processing the changeset, so keep a reference to the session.
	//thisSession := SessionInfo[client.ID]
	return []byte("handleUserChanges")
}
