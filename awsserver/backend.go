package main

import (
	"log"
	"net/http"
	"sync"
)

var (
	activeBackendMap     = make(map[uint]chan<- *Drink)
	activeBackendMapLock = sync.RWMutex{}
)

type ClientHandshake struct {
	Message string `json:"message"`
	Name    string `json:"name"`
}

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
	clientInfo := ClientHandshake{}
	err = conn.ReadJSON(&clientInfo)
	if err != nil {
		log.Printf("Error decoding json: %v\n", err)
		return
	}

	log.Printf("Got client info: %s", clientInfo)

	machine := Machine{}
	db.Where(&Machine{Name: clientInfo.Name}).First(&machine)
	if machine.ID == 0 {
		machine.Name = clientInfo.Name
		db.Create(&machine)
	}

	// Give our end of the handshake
	data := make(map[string]interface{})
	data["status"] = "successful"
	err = conn.WriteJSON(data)
	if err != nil {
		log.Printf("Error writing to connection: %v\n", err)
		return
	}

	dataChan := make(chan *Drink)
	activeBackendMapLock.Lock()
	if existing, ok := activeBackendMap[machine.ID]; ok {
		log.Println("Closing existing connection.")
		close(existing)
	}
	activeBackendMap[machine.ID] = dataChan
	activeBackendMapLock.Unlock()

	for drk := range dataChan {
		data := make(map[string]interface{})
		data["name"] = drk.Name
		data["mug_size"] = drk.Size
		addIns := make(map[string]interface{})
		addIns["sugar"] = drk.Sugar
		addIns["creamer"] = drk.Creamer
		addIns["tea_bag"] = drk.TeaBag
		addIns["k_cup"] = drk.KCup
		data["add_ins"] = addIns

		err = conn.WriteJSON(data)
		if err != nil {
			log.Printf("Error writing to connection: %v\n", err)
			break
		}

		log.Printf("Sent request to client:\n%v\n", data)
	}

	activeBackendMapLock.Lock()
	if existing, ok := activeBackendMap[machine.ID]; ok && existing == dataChan {
		log.Println("Closing client connection: cleaning up channel.")
		close(existing)
		delete(activeBackendMap, machine.ID)
	} else {
		log.Println("Closing client connection: connection has been replaced.")
	}
	activeBackendMapLock.Unlock()

}
