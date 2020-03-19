package poker

import (
	"bytes"
	"encoding/json"
	"github.com/Michael2008S/etherpad4go/api"
	"github.com/Michael2008S/etherpad4go/model"
	bgStore "github.com/Michael2008S/etherpad4go/store"
	"github.com/Michael2008S/etherpad4go/utils/changeset"
	"github.com/y0ssar1an/q"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

//type PlayServer struct {
//	http.Handler
//}

// func NewPlayerServer() (*PlayServer, error) {
// 	p := new(PlayServer)
// 	router := http.NewServeMux()
// 	router.Handle("/ws", http.HandlerFunc(p.webSocket))
// 	p.Handler = router
// 	return p, nil
// }

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//func (p *PlayServer) webSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
//	conn, err := wsUpgrader.Upgrade(w, r, nil)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
//	client.hub.register <- client
//
//	// Allow collection of memory referenced by the caller by doing all work in
//	// new goroutines.
//	go client.writePump()
//	go client.readPump()
//}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// 存储层
	dbStore bgStore.Store
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		//q.Q("readPump_message:", string(message))
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		//q.Q("readPump_message_TrimSpace:", string(message))
		c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, dbStore bgStore.Store) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client
	client.dbStore = dbStore

	// 发送 CLIENT_VARS 数据
	sendClientVars(conn)

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

// 发送客户端数据：
func sendClientVars(conn *websocket.Conn) {

	// TODO pid
	// load the pad-object from the database
	pad := model.NewPad("q","")

	atext := pad.AText
	attribsForWire := changeset.Pre

	clientVarsData := api.ClientVarsDataResp{}
	clientVarsData.SkinName = "no-skin"
	collabClientVars := api.CollabClientVars{
		InitialAttributedText: struct {
			Text    string `json:"text"`
			Attribs string `json:"attribs"`
		}{Text: "asdfasdaff", Attribs: "|5+6b*0|3+7*0+4*1|2+2*1+c*0+c*1|3+j*0+4*1|3+9*0|1+r|1+1"},
		ClientIP: "127.0.0.1",
		PadID:    "q",
		HistoricalAuthorData: struct {
			APy0WdSkbof4TM4DD struct {
				Name    interface{} `json:"name"`
				ColorID int         `json:"colorId"`
			} `json:"a.Py0WdSkbof4tM4DD"`
			AYjK4P2YxGHx8NNgf struct {
				Name    string `json:"name"`
				ColorID string `json:"colorId"`
			} `json:"a.YjK4P2yxGHx8NNgf"`
		}{},
		Apool: struct {
			NumToAttrib struct {
				Num0 []string `json:"0"`
				Num1 []string `json:"1"`
			} `json:"numToAttrib"`
			NextNum int `json:"nextNum"`
		}{},
		Rev:  0,
		Time: 0,
	}
	clientVarsData.UserID = "123"
	clientVarsData.CollabClientVars = collabClientVars
	clientVarsData.ColorPalette = model.ColorPalette
	clientVarsResp := api.WarpMsgResp{
		Type: api.MsgTypeClientVars,
		Data: clientVarsData,
	}
	w, err := conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}
	q.Q(clientVarsResp)
	resp, _ := json.Marshal(clientVarsResp)
	w.Write(resp)
	if err := w.Close(); err != nil {
		return
	}
}
