package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var address = flag.String("address", ":8080", "The address to bind on.")

var upgrader = websocket.Upgrader{}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Unable to upgrade connection: %v\n", err)
		http.Error(w, "Connection failed.", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	data := make(map[string]interface{})
	err = conn.ReadJSON(&data)
	if err != nil {
		log.Printf("Error decoding json: %v\n", err)
		return
	}

	log.Printf("Got data: %s", data)

	data = make(map[string]interface{})
	data["hi_there"] = "This is the server"
	err = conn.WriteJSON(data)
	if err != nil {
		log.Printf("Error decoding json: %v\n", err)
		return
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ws", wsHandler)

	log.Printf("Starting on %s\n", *address)
	log.Fatalln(http.ListenAndServe(*address, r))
}
