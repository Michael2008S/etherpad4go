package main

import (
	"flag"
	"log"
	"net/http"

	poker "github.com/Michael2008S/etherpad4go"
	"github.com/gorilla/mux"
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
	r := mux.NewRouter()
	r.HandleFunc("/", serveHome)
	r.HandleFunc("/p/{name}", servePad)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./template/static"))))
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		poker.ServeWs(hub, w, r)
	})
	http.Handle("/", r)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
