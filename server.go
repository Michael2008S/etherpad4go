package poker

import (
	"bytes"
	"encoding/json"
	"github.com/Michael2008S/etherpad4go/api"
	"github.com/Michael2008S/etherpad4go/model"
	bgStore "github.com/Michael2008S/etherpad4go/store"
	"github.com/Michael2008S/etherpad4go/utils/changeset"
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
	ID string

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
				log.Printf("主动断开链接1：error: %v", err)
			}
			log.Printf("主动断开链接2：error: %v", err)
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		//q.Q("readPump_message_TrimSpace:", string(message))
		c.hub.broadcast <- InboundMsg{from: c, message: message}
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

	// 使用Sec-WebSocket-Key当链接key
	id := r.Header.Get("Sec-WebSocket-Key")

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	log.Printf("header: Sec-WebSocket-Key is \" %v \" \n", id)
	client.ID = id
	client.dbStore = dbStore
	client.hub.register <- client

	if websocket.IsWebSocketUpgrade(r) {
		log.Println("收到websocket链接")

		// 发送 CLIENT_VARS 数据
		//sendClientVars(hub, client, client.dbStore)

	} else {
		log.Println("您这也不是websocket啊")
		w.Write([]byte(`您这也不是websocket啊`))
		return
	}

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

// 发送客户端数据：
func sendClientVars(hub *Hub, client *Client, db bgStore.Store) {

	// TODO pid
	// load the pad-object from the database
	pad := model.NewPad("q", "", db)

	atext := pad.AText
	//atext.Attribs = "|4+2l"
	//atext.Text = "Welcome to Etherpad!\n\nThis pad text is synchronized~ https://github.com/ether/etherpad-lite\n\n"
	translated, newPool := changeset.PrepareForWire(atext.Attribs, pad.Pool)
	//apool := newPool.ToJsonAble
	atext.Attribs = translated

	//historicalAuthorData :=  map[string]

	clientVarsData := api.ClientVarsDataResp{}
	clientVarsData.SkinName = "no-skin"
	clientVarsData.PadID = pad.Id
	collabClientVars := api.CollabClientVars{
		InitialAttributedText: atext,
		ClientIP:              "127.0.0.1",
		PadID:                 pad.Id,
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
		Apool: newPool,
		Rev:   pad.GetHeadRevisionNumber(),
		Time:  time.Now().Unix(),
	}
	clientVarsData.UserID = sessionInfo[client.ID].author
	clientVarsData.CollabClientVars = collabClientVars
	clientVarsData.ColorPalette = model.ColorPalette
	clientVarsResp := api.WarpMsgResp{
		Type: api.MsgTypeClientVars,
		Data: clientVarsData,
	}
	w, err := client.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		log.Println(err)
		return
	}
	//q.Q(clientVarsResp)
	resp, _ := json.Marshal(clientVarsResp)
	w.Write(resp)
	if err := w.Close(); err != nil {
		return
	}

	//sessionInfo[client.ID].rev =pad.GetHeadRevisionNumber()
	//sessionInfo[client.ID].author =

	// prepare the notification for the other users on the pad, that this user joined
	userInfo := api.UserInfo{
		Ip:        "127.0.0.1",
		ColorId:   0,
		UserAgent: "Anonymous",
		UserId:    sessionInfo[client.ID].author,
	}
	authInfo, ok := sessionInfo[client.ID]
	if ok {
		authInfo.rev = pad.GetHeadRevisionNumber()
		userInfo.UserId = authInfo.author
	}
	messageToTheOtherUsers := api.UserNewInfoResp{
		Type: "COLLABROOM",
		Data: struct {
			Type     string       `json:"type"`
			UserInfo api.UserInfo `json:"userInfo"`
		}{Type: "USER_NEWINFO",
			UserInfo: userInfo},
	}
	messageToTheOtherUsersResp, _ := json.Marshal(&messageToTheOtherUsers)

	// 发给我的消息
	messageToMe := api.UserNewInfoResp{
		Type: "COLLABROOM",
		Data: struct {
			Type     string       `json:"type"`
			UserInfo api.UserInfo `json:"userInfo"`
		}{Type: "USER_NEWINFO",
		},
	}

	// notify all existing users about new user
	// TODO client.broadcast.to(padIds.padId).json.send(messageToTheOtherUsers);
	// Get sessions for this pad and update them (in parallel)
	roomClients := _getRoomClients(pad.Id)
	for _, clientID := range roomClients {
		roomCli := hub.clients[clientID]
		if clientID == client.ID || roomCli == nil {
			continue
		}
		// 发送我到pad 房间其他人
		w, err := roomCli.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			log.Println(err)
			continue
		}
		w.Write(messageToTheOtherUsersResp)
		if err := w.Close(); err != nil {
			log.Println(err)
			//return
		}

		// 发送其他已有的用户到我的 client 中
		toMyUserInfo := api.UserInfo{
			Ip:        "127.0.0.1",
			ColorId:   0,
			UserAgent: "Anonymous",
			UserId:    sessionInfo[clientID].author,
		}
		myWrite, err := client.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			log.Println(err)
		}
		messageToMe.Data.UserInfo = toMyUserInfo
		messageToMeResp, _ := json.Marshal(&messageToMe)
		myWrite.Write(messageToMeResp)
		if err := myWrite.Close(); err != nil {
			log.Println(err)
		}
	}
}
