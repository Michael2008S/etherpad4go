package main

import (
	"flag"
	"log"
	"net/http"

	poker "github.com/Michael2008S/etherpad4go"
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/googollee/go-socket.io"
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

func main() {
	flag.Parse()
	hub := poker.NewHub()
	go hub.Run()

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", serveHome)
	r.HandleFunc("/p/{name}", servePad)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./template/static"))))
	//r.HandleFunc("/socket.io", func(w http.ResponseWriter, r *http.Request) {
	//	poker.ServeWs(hub, w, r)
	//})

	go server.Serve()
	defer server.Close()
	http.Handle("/socket.io/", server)

	http.Handle("/", r)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:9011"},
		AllowCredentials: true,
	})
	handler := c.Handler(r)
	err = http.ListenAndServe(*addr, handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
