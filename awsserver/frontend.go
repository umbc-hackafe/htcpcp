package main

import (
	"flag"
	"log"
	"net/http"

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
)

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
			"Created dummy machine with id %d and name %s\n",
			m.ID, m.Name)
	}

	b := Drink{
		Name: "__bell__",
		Sugar: 1,
	}
	db.Create(&b)
	log.Printf("Created __bell__ drink with id %d\n", b.ID)

	l1 := Drink{
		Name: "__light__",
		Sugar: 1,
	}
	db.Create(&l1)
	log.Printf("Created __light__ (on) drink with id %d\n", l1.ID)

	l2 := Drink{
		Name: "__light__",
		Sugar: 0,
	}
	db.Create(&l2)
	log.Printf("Created __light__ (off) drink with id %d\n", l2.ID)

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

	apiRouter.HandleFunc("/api/brew", brew)
	apiRouter.HandleFunc("/api/brew/{drink:[0-9]+}/{machine:[0-9]+}", brewPath)

	rootRouter.Handle("/api/", apiRouter)

	// Start the server
	log.Printf("Starting on %s\n", *address)
	log.Fatalln(http.ListenAndServe(*address, rootRouter))
}
