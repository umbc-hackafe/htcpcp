package main

import (
	"flag"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
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

	createDummyMachine = flag.Bool(
		"create-dummy", false,
		"Whether to create a dummy machine for the frontend to access.")
)

var (
	upgrader = websocket.Upgrader{}
	db       gorm.DB

	activeBackendMap     = make(map[uint]chan<- *Drink)
	activeBackendMapLock = sync.RWMutex{}
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the http connection to a websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Unable to upgrade connection: %v\n", err)
		http.Error(w, "Connection failed.", http.StatusInternalServerError)
		return
	}
	// If successful, set to close after returning
	defer conn.Close()

	// Test Code
	data := make(map[string]interface{})
	err = conn.ReadJSON(&data)
	if err != nil {
		log.Printf("Error decoding json: %v\n", err)
		return
	}

	log.Printf("Got data: %s", data)

	// Give our end of the handshake
	data = make(map[string]interface{})
	data["status"] = "successful"
	err = conn.WriteJSON(data)
	if err != nil {
		log.Printf("Error writing to connection: %v\n", err)
		return
	}

	data = make(map[string]interface{})
	data["mug_size"] = 8
	data["add_ins"] = make(map[string]interface{})
	data["name"] = "coffelattachino"
	err = conn.WriteJSON(data)
	if err != nil {
		log.Printf("Error writing to connection: %v\n", err)
		return
	}

	for {
	}
}

func main() {
	flag.Parse()

	// Open the database connection
	var err error
	db, err = gorm.Open(*sqlDriver, *sqlConnectionString)
	if err != nil {
		log.Fatalln(err)
	}

	// Run migrate to make sure the tables exist
	db.AutoMigrate(&Schedule{})
	db.AutoMigrate(&Drink{})
	db.AutoMigrate(&Machine{})

	if *createDummyMachine {
		m := Machine{Name: "Dummy"}
		db.Create(&m)
		log.Printf(
			"Created dummy machine with id %d and name %s",
			m.ID, m.Name)
	}

	// Create the base router
	rootRouter := http.NewServeMux()
	// Serve the static files (js, html, css)
	rootRouter.Handle("/", http.FileServer(http.Dir(*staticFilesPath)))
	// Handle the websocket connections
	rootRouter.HandleFunc("/ws", wsHandler)

	// Sub-router for the REST api
	apiRouter := mux.NewRouter()
	apiRouter.Methods("POST").
		Path("/api/update/schedule").
		HandlerFunc(createSchedule)
	apiRouter.Methods("POST").
		Path("/api/update/drink").
		HandlerFunc(createDrink)

	apiRouter.HandleFunc("/api/get/schedules", getSchedules)
	apiRouter.HandleFunc("/api/get/drinks", getDrinks)
	apiRouter.HandleFunc("/api/get/machines", getMachines)

	apiRouter.HandleFunc("/brew", brew)
	apiRouter.HandleFunc("/brew/{device:[0-9]+}/{machine:[0-9]+}", brewPath)

	rootRouter.Handle("/api/", apiRouter)

	// Start the server
	log.Printf("Starting on %s\n", *address)
	log.Fatalln(http.ListenAndServe(*address, rootRouter))
}
