package poker

import (
	"bytes"
	"encoding/json"
	"github.com/Michael2008S/etherpad4go/api"
	"github.com/y0ssar1an/q"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type PlayServer struct {
	http.Handler
}

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

func (p *PlayServer) webSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

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
		q.Q("readPump_message:", string(message))
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		q.Q("readPump_message_TrimSpace:", string(message))
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
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// 发送 CLIENT_VARS 数据
	sendClientVars(conn)

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

// 发送客户端数据：
func sendClientVars(conn *websocket.Conn) {
	//6384:42["message",{"type":"CLIENT_VARS","data":{"skinName":"no-skin","accountPrivs":{"maxRevisions":100},"automaticReconnectionTimeout":0,"initialRevisionList":[],"initialOptions":{"guestPolicy":"deny"},"savedRevisions":[],"collab_client_vars":{"initialAttributedText":{"text":"Welcome to Etherpad!\n\nThis pad text is synchronized as you type, so that everyone viewing this page sees the same text. This allows you to collaborate seamlessly on documents!\n\nGet involved with Etherpad at http://etherpad.org\n\nasdf\n\nadsf\n\nasdfasdfadsf阿斯顿发阿斯顿发阿斯顿发\nasdfasdfsdaf\nasdf\n给第三方\nad\nasfd\nsafaafdasdf asdf   tehsi. \n\n","attribs":"|5+6b*0|3+7*0+4*1|2+2*1+c*0+c*1|3+j*0+4*1|3+9*0|1+r|1+1"},"clientIp":"127.0.0.1","padId":"q","historicalAuthorData":{"a.Py0WdSkbof4tM4DD":{"name":null,"colorId":36},"a.YjK4P2yxGHx8NNgf":{"name":"","colorId":"#c4e7b1"}},"apool":{"numToAttrib":{"0":["author","a.Py0WdSkbof4tM4DD"],"1":["author","a.YjK4P2yxGHx8NNgf"]},"nextNum":2},"rev":33,"time":1583675281484},"colorPalette":["#ffc7c7","#fff1c7","#e3ffc7","#c7ffd5","#c7ffff","#c7d5ff","#e3c7ff","#ffc7f1","#ffa8a8","#ffe699","#cfff9e","#99ffb3","#a3ffff","#99b3ff","#cc99ff","#ff99e5","#e7b1b1","#e9dcAf","#cde9af","#bfedcc","#b1e7e7","#c3cdee","#d2b8ea","#eec3e6","#e9cece","#e7e0ca","#d3e5c7","#bce1c5","#c1e2e2","#c1c9e2","#cfc1e2","#e0bdd9","#baded3","#a0f8eb","#b1e7e0","#c3c8e4","#cec5e2","#b1d5e7","#cda8f0","#f0f0a8","#f2f2a6","#f5a8eb","#c5f9a9","#ececbb","#e7c4bc","#daf0b2","#b0a0fd","#bce2e7","#cce2bb","#ec9afe","#edabbd","#aeaeea","#c4e7b1","#d722bb","#f3a5e7","#ffa8a8","#d8c0c5","#eaaedd","#adc6eb","#bedad1","#dee9af","#e9afc2","#f8d2a0","#b3b3e6"],"clientIp":"127.0.0.1","userIsGuest":true,"userColor":36,"padId":"q","padOptions":{"noColors":false,"showControls":true,"showChat":true,"showLineNumbers":true,"useMonospaceFont":false,"userName":false,"userColor":false,"rtl":false,"alwaysShowChat":false,"chatAndUsers":false,"lang":"en-gb"},"padShortcutEnabled":{"altF9":true,"altC":true,"cmdShift2":true,"delete":true,"return":true,"esc":true,"cmdS":true,"tab":true,"cmdZ":true,"cmdY":true,"cmdI":true,"cmdB":true,"cmdU":true,"cmd5":true,"cmdShiftL":true,"cmdShiftN":true,"cmdShift1":true,"cmdShiftC":true,"cmdH":true,"ctrlHome":true,"pageUp":true,"pageDown":true},"initialTitle":"Pad: q","opts":{},"chatHead":3,"numConnectedUsers":0,"readOnlyId":"r.c97d8bdd5f1e326473442fc2f55793c3","readonly":false,"serverTimestamp":1583918908724,"userId":"a.Py0WdSkbof4tM4DD","abiwordAvailable":"no","sofficeAvailable":"no","exportAvailable":"no","plugins":{"plugins":{"ep_etherpad-lite":{"parts":[{"name":"express","hooks":{"createServer":"ep_etherpad-lite/node/hooks/express:createServer","restartServer":"ep_etherpad-lite/node/hooks/express:restartServer"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/express"},{"name":"static","hooks":{"expressCreateServer":"ep_etherpad-lite/node/hooks/express/static:expressCreateServer"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/static"},{"name":"i18n","hooks":{"expressCreateServer":"ep_etherpad-lite/node/hooks/i18n:expressCreateServer"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/i18n"},{"name":"specialpages","hooks":{"expressCreateServer":"ep_etherpad-lite/node/hooks/express/specialpages:expressCreateServer"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/specialpages"},{"name":"socketio","hooks":{"expressCreateServer":"ep_etherpad-lite/node/hooks/express/socketio:expressCreateServer"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/socketio"},{"name":"apicalls","hooks":{"expressCreateServer":"ep_etherpad-lite/node/hooks/express/apicalls:expressCreateServer"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/apicalls"},{"name":"webaccess","hooks":{"expressConfigure":"ep_etherpad-lite/node/hooks/express/webaccess:expressConfigure"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/webaccess"},{"name":"swagger","hooks":{"expressCreateServer":"ep_etherpad-lite/node/hooks/express/swagger:expressCreateServer"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/swagger"}],"package":{"name":"etherpad-lite","version":"1.0.0","description":"my customer etherpad","main":"index.js","scripts":{"test":"echo \"Error: no test specified\" && exit 1"},"author":"","license":"ISC","invalid":true,"realName":"ep_etherpad-lite","path":"/Volumes/RamDisk/goEtherpad/etherpad4go/etherpad-lite/node_modules/ep_etherpad-lite","realPath":"/Volumes/RamDisk/goEtherpad/etherpad4go/etherpad-lite/src","link":"/Volumes/RamDisk/goEtherpad/etherpad4go/etherpad-lite/src","depth":1}}},"parts":[{"name":"swagger","hooks":{"expressCreateServer":"ep_etherpad-lite/node/hooks/express/swagger:expressCreateServer"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/swagger"},{"name":"webaccess","hooks":{"expressConfigure":"ep_etherpad-lite/node/hooks/express/webaccess:expressConfigure"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/webaccess"},{"name":"apicalls","hooks":{"expressCreateServer":"ep_etherpad-lite/node/hooks/express/apicalls:expressCreateServer"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/apicalls"},{"name":"socketio","hooks":{"expressCreateServer":"ep_etherpad-lite/node/hooks/express/socketio:expressCreateServer"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/socketio"},{"name":"specialpages","hooks":{"expressCreateServer":"ep_etherpad-lite/node/hooks/express/specialpages:expressCreateServer"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/specialpages"},{"name":"i18n","hooks":{"expressCreateServer":"ep_etherpad-lite/node/hooks/i18n:expressCreateServer"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/i18n"},{"name":"static","hooks":{"expressCreateServer":"ep_etherpad-lite/node/hooks/express/static:expressCreateServer"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/static"},{"name":"express","hooks":{"createServer":"ep_etherpad-lite/node/hooks/express:createServer","restartServer":"ep_etherpad-lite/node/hooks/express:restartServer"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/express"}]},"indentationOnNewLine":true,"scrollWhenFocusLineIsOutOfViewport":{"percentage":{"editionAboveViewport":0,"editionBelowViewport":0},"duration":0,"scrollWhenCaretIsInTheLastLineOfViewport":false,"percentageToScrollWhenUserPressesArrowUp":0},"initialChangesets":[]}}]

	clientVarsData := api.ClientVarsDataResp{}
	clientVars := api.WarpMsgResp{
		Type: api.MsgTypeClientVars,
		Data: clientVarsData,
	}
	w, err := conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}
	resp, _ := json.Marshal(clientVars)
	w.Write(resp)
	if err := w.Close(); err != nil {
		return
	}
}
