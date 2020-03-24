package poker

import (
	"encoding/json"
	"errors"
	"github.com/Michael2008S/etherpad4go/api"
	"github.com/Michael2008S/etherpad4go/model"
	bgStore "github.com/Michael2008S/etherpad4go/store"
	"github.com/Michael2008S/etherpad4go/utils/changeset"
	"github.com/y0ssar1an/q"
	"log"
	"strings"
)
import "github.com/thedevsaddam/gojsonq"

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan InboundMsg

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// 存储层
	dbStore bgStore.Store
}

type InboundMsg struct {
	from    *Client
	message []byte
}

func NewHub(db bgStore.Store) *Hub {
	return &Hub{
		broadcast:  make(chan InboundMsg),
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
			q.Q("Run_broadcast:", string(message.message))

			responseMsg := h.handleMessage(message)
			q.Q("Send_broadcast:", string(responseMsg))

			//for client := range h.clients {
			//	select {
			//	case client.send <- responseMsg:
			//	default:
			//		close(client.send)
			//		delete(h.clients, client)
			//	}
			//}
		}
	}
}

func (h *Hub) handleMessage(message InboundMsg) {
	msgType := gojsonq.New().FromString(string(message.message)).Find("type")
	//
	//accessStatus =  securityManager.checkAccess(padID, sessionCookie, token, password)

	//if (accessStatus !== "grant") {
	//
	//}

	if msgType.(string) == "CLIENT_READY" {
		//	handleClientReady(client, message);
		clientReadyReq := api.ClientReadyReq{}
		json.Unmarshal(message.message, &clientReadyReq)
		createSessionInfo(message.from, clientReadyReq)
		q.Q("createSessionInfo")
	} else if msgType.(string) == "CHANGESET_REQ" {
		//	handleChangesetRequest(client, message);
	} else if msgType.(string) == "COLLABROOM" {
		msgDataType := gojsonq.New().FromString(string(message.message)).Find("type")
		if msgDataType == "USER_CHANGES" {
			// TODO padChannels.emit(message.padId, {client: client, message: message}); // add to pad queue
			handleUserChanges(message)
		}
	}

	//return []byte(`{"type":"error"}`)
}

func handleClientReady() {

}

func handleChangesetRequest() {

}

func handleUserChanges(msg InboundMsg) error {
	reqMsg := api.CollabRoomReqMessage{}
	json.Unmarshal(msg.message, &reqMsg)

	// get all Vars we need
	baseRev := reqMsg.Data.BaseRev
	wireApool := reqMsg.Data.Apool
	cs := reqMsg.Data.Changeset

	// The client might disconnect between our callbacks. We should still
	// finish processing the changeset, so keep a reference to the session.
	thisSession := sessionInfo[msg.from.ID]
	pad := model.NewPad(thisSession.padID, "", msg.from.dbStore)
	// Verify that the changeset has valid syntax and is in canonical form
	chgset := changeset.ChangeSet{}
	if err := chgset.CheckRep(cs); err != nil {
		log.Println(err)
		//return nil
	}
	// Verify that the attribute indexes used in the changeset are all
	// defined in the accompanying attribute pool.
	//chgset.EachAttribNumber(cs)

	// Validate all added 'author' attribs to be the same value as the current user

	// ex. applyUserChanges
	apool := pad.Pool
	r := baseRev

	// The client's changeset might not be based on the latest revision,
	// since other clients are sending changes at the same time.
	// Update the changeset so that it can be applied to the latest revision.
	for ; r < pad.GetHeadRevisionNumber(); r++ {
		c := pad.GetRevisionChangeset(r)
		// At this point, both "c" (from the pad) and "changeset" (from the
		// client) are relative to revision r - 1. The follow function
		// rebases "changeset" so that it is relative to revision r
		// and can be applied after "c".
		if baseRev+1 == r && c == cs {
			//FIXME client.json.send({disconnect:"badChangeset"});
			return  errors.New("Won't apply USER_CHANGES, because it contains an already accepted changeset")
		}

	}

	prevText := pad.GetText()

	pad.AppendRevision(cs, thisSession.author)

	//correctionChangeset := _correctMarkersInPad(pad.AText, pad.Pool)
	//if correctionChangeset {
	//	pad.AppendRevision(correctionChangeset)
	//}

	// Make sure the pad always ends with an empty line.
	if strings.LastIndex(pad.GetText(), "\n") != len(pad.GetText())-1 {
		nlChangeset := chgset.MakeSplice(pad.GetText(), len(pad.GetText())-1, 0, "\n", "", "")
		pad.AppendRevision(nlChangeset, "")
	}
	updatePadClients(pad)
}

func updatePadClients(pad model.Pad) {
	// skip this if no-one is on this pad

}

/**
 * TODO Copied from the Etherpad Source Code. Don't know what this method does excatly...
 */
func _correctMarkersInPad(atext changeset.AText, pool changeset.AttributePool) string {
	text := atext.Text
	// collect char positions of line markers (e.g. bullets) in new atext
	// that aren't at the start of a line
	badMarkers := []string{}
	iter := changeset.NewOperatorIterator(atext.Attribs, 0)
	//offset := 0
	for iter.HasNext() {
		//op := iter.Next()

	}
	if len(badMarkers) == 0 {
		return ""
	}
	// create changeset that removes these bad markers
	//offset = 0
	builder := changeset.NewBuilder(len(text))

	return builder.ToString()
}

func DisconnectBadChangeset() {

}
