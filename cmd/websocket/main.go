package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"

	poker "github.com/Michael2008S/etherpad4go"
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/y0ssar1an/q"
	//socketio "github.com/googollee/go-socket.io"
)

var addr = flag.String("addr", ":8800", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// http.ServeFile(w, r, "home.html")
	http.ServeFile(w, r, "./template/index.html")
}

func servePad(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	http.ServeFile(w, r, "./template/pad.html")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")
	err = ws.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}
	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	reader(ws)
}
func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		log.Println("message_type:")
		log.Println(messageType, string(p))
		q.Q(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}

func main() {
	flag.Parse()
	hub := poker.NewHub()
	go hub.Run()


	r := mux.NewRouter()
	r.HandleFunc("/", serveHome)
	r.HandleFunc("/p/{name}", servePad)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./template/static"))))
	r.HandleFunc("/socket.io", func(w http.ResponseWriter, r *http.Request) {
		poker.ServeWs(hub, w, r)
	})

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		poker.ServeWs(hub, w, r)
	})

	//r.HandleFunc("/ws", wsEndpoint)

	http.Handle("/test", r)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:9011"},
		AllowCredentials: true,
	})
	handler := c.Handler(r)
	err := http.ListenAndServe(*addr, handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
