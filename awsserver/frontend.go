package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var (
	address         = flag.String("address", ":8080", "The address to bind on.")
	staticFilesPath = flag.String(
		"static-files-path", "../frontend/", "The static files to use.")

	sqlConnectionString = flag.String(
		"sql-connect-string", ":memory:",
		"Connection mode string to use to connect to the datablase")
	sqlDriver = flag.String("sql-driver", "sqlite3", "Sql driver to use")
)

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
	flag.Parse()

	db, err := gorm.Open(*sqlDriver, *sqlConnectionString)
	if err != nil {
		log.Fatalln(err)
	}

	_ = db

	rootRouter := http.NewServeMux()
	rootRouter.Handle("/", http.FileServer(http.Dir(*staticFilesPath)))
	rootRouter.HandleFunc("/ws", wsHandler)

	log.Printf("Starting on %s\n", *address)
	log.Fatalln(http.ListenAndServe(*address, rootRouter))
}
