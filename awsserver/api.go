package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

var maximumRequestSize = flag.Int(
	"max-request-size", 1024*1024, "The maximum body size for incoming requests.")

// Iterate over the given strings, check that they are valid days, and turn them into a
// bit vector if they are.
func daysToBitvector(days []string) (uint8, error) {
	m := uint8(0)
	for _, v := range days {
		vUp := strings.ToUpper(v)
		if d, ok := DaysMap[vUp]; ok {
			m |= d
		} else {
			return 0, fmt.Errorf("No such day: %s", v)
		}
	}
	return m, nil
}

func bitvectorToDays(days uint8) []string {
	d := []string{}
	for i := uint(0); i < 7; i += 1 {
		if (days & (1 << i)) != 0 {
			d = append(d, DaysReverseMap[1<<i])
		}
	}
	return d
}

type ScheduleCreateRequest struct {
	ID      uint     `json:"id"`
	Name    string   `json:"name"`
	Days    []string `json:"days"`
	Time    int      `json:"time"`
	Enabled bool     `json:"bool"`
	Drink   uint     `json:"drink"`
	Machine uint     `json:"machine"`
}

func createSchedule(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(io.LimitReader(r.Body, int64(*maximumRequestSize)))

	data := ScheduleCreateRequest{}

	err := d.Decode(&data)
	if err != nil {
		log.Printf("Error decoding request: %v\n", err)
		http.Error(w, "Unable to decode request", http.StatusBadRequest)
		return
	}

	// Convert the days list to the bit vector
	days, err := daysToBitvector(data.Days)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}

	if data.Time < 0 {
		data.Time = 0
	} else if data.Time > 24*60*60 {
		data.Time = 24 * 60 * 60
	}

	drink := Drink{}
	db.First(&drink, data.Drink)
	data.Drink = drink.ID // 0 if missing, otherwise data.Drink

	machine := Machine{}
	db.First(&machine, data.Machine)
	data.Machine = machine.ID // 0 if missing, otherwise data.Machine

	schedule := Schedule{}
	db.First(&schedule, data.ID)
	new := schedule.ID == 0

	schedule.Name = data.Name[:100]
	schedule.Days = days
	schedule.Enabled = data.Enabled
	schedule.Time = data.Time
	schedule.DrinkID = data.Drink
	schedule.MachineID = data.Machine

	db.Save(&schedule)
	if new {
		log.Printf(
			"Created new schedule with id %d and name %s", schedule.ID, schedule.Name)
	} else {
		log.Printf("Updated schedule %d (name: %s)", schedule.ID, schedule.Name)
	}

	data.ID = schedule.ID
	data.Name = schedule.Name
	data.Days = bitvectorToDays(schedule.Days)
	data.Enabled = schedule.Enabled
	data.Time = schedule.Time
	data.Drink = schedule.DrinkID
	data.Machine = schedule.MachineID

	enc := json.NewEncoder(w)
	err = enc.Encode(data)
	if err != nil {
		log.Printf("Error while encoding: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}

func getSchedules(w http.ResponseWriter, r *http.Request) {
	schedules := []Schedule{}
	db.Find(&schedules)

	sch := make([]ScheduleCreateRequest, len(schedules))
	for i, schedule := range schedules {
		sch[i].ID = schedule.ID
		sch[i].Name = schedule.Name
		sch[i].Days = bitvectorToDays(schedule.Days)
		sch[i].Enabled = schedule.Enabled
		sch[i].Time = schedule.Time
		sch[i].Drink = schedule.DrinkID
		sch[i].Machine = schedule.MachineID
	}
	wrap := map[string]interface{}{
		"schedules": sch,
	}

	enc := json.NewEncoder(w)
	err := enc.Encode(wrap)
	if err != nil {
		log.Printf("Error while encoding: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	log.Printf("Got %d schedules.", len(schedules))
}

type DrinkCreateRequest struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Size    uint8  `json:"size"`
	Sugar   uint8  `json:"sugar"`
	Creamer uint8  `json:"creamer"`
	TeaBag  string `json:"tea_bag"`
	KCup    string `json:"k_cup"`
}

func createDrink(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(io.LimitReader(r.Body, int64(*maximumRequestSize)))

	data := DrinkCreateRequest{}

	err := d.Decode(&data)
	if err != nil {
		log.Printf("Error decoding request: %v\n", err)
		http.Error(w, "Unable to decode request", http.StatusBadRequest)
		return
	}

	drink := Drink{}
	db.First(&drink, data.ID)
	new := drink.ID == 0

	drink.Name = data.Name[:100]
	drink.Size = data.Size
	drink.Sugar = data.Sugar
	drink.Creamer = data.Creamer
	drink.TeaBag = data.TeaBag[:100]
	drink.KCup = data.KCup[:100]

	db.Save(&drink)
	if new {
		log.Printf("Created a new drink with ID %d and name %s\n", drink.ID, drink.Name)
	} else {
		log.Printf("Updated drink %d (name %s)\n", drink.ID, drink.Name)
	}

	data.ID = drink.ID
	data.Name = drink.Name
	data.Size = drink.Size
	data.Sugar = drink.Sugar
	data.Creamer = drink.Creamer
	data.TeaBag = drink.TeaBag
	data.KCup = drink.KCup

	enc := json.NewEncoder(w)
	err = enc.Encode(data)
	if err != nil {
		log.Printf("Error while encoding: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}

func getDrinks(w http.ResponseWriter, r *http.Request) {
	drinks := []Drink{}
	db.Find(&drinks)

	dri := make([]DrinkCreateRequest, len(drinks))
	for i, drink := range drinks {
		dri[i].ID = drink.ID
		dri[i].Name = drink.Name
		dri[i].Size = drink.Size
		dri[i].Sugar = drink.Sugar
		dri[i].Creamer = drink.Creamer
		dri[i].TeaBag = drink.TeaBag
		dri[i].KCup = drink.KCup
	}

	wrap := map[string]interface{}{
		"drinks": dri,
	}

	enc := json.NewEncoder(w)
	err := enc.Encode(wrap)
	if err != nil {
		log.Printf("Error while encoding: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	log.Printf("Got %d drinks.\n", len(drinks))
}

func getMachines(w http.ResponseWriter, r *http.Request) {
	machines := []Machine{}
	db.Find(&machines)

	wrap := map[string]interface{}{
		"machines": machines,
	}

	enc := json.NewEncoder(w)
	err := enc.Encode(wrap)
	if err != nil {
		log.Printf("Error while encoding: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	log.Printf("Got %d machines.\n", len(machines))
}

type BrewRequest struct {
	Drink   uint `json:"drink"`
	Machine uint `json:"machine"`
}

func brew(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(io.LimitReader(r.Body, int64(*maximumRequestSize)))

	req := BrewRequest{}

	err := d.Decode(&req)
	if err != nil {
		log.Printf("Error decoding request: %v\n", err)
		http.Error(w, "Unable to decode request", http.StatusBadRequest)
		return
	}

	enqueueBrew(&req, w)
}

func brewPath(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	drink, err := strconv.ParseUint(vars["drink"], 10, 64)
	if err != nil {
		log.Panicf(
			"Illegal string-integer converstion of vars[\"drink\"] = %s.\n",
			vars["drink"])
	}

	machine, err := strconv.ParseUint(vars["machine"], 10, 64)
	if err != nil {
		log.Panicf(
			"Illegal string-integer converstion of vars[\"machine\"] = %s.\n",
			vars["machine"])
	}

	enqueueBrew(&BrewRequest{
		Drink:   uint(drink),
		Machine: uint(machine),
	}, w)
}

func enqueueBrew(req *BrewRequest, w http.ResponseWriter) {
	drink := Drink{}
	db.First(&drink, req.Drink)

	if drink.ID == 0 {
		log.Printf("Brew request: unknown drink %d\n", req.Drink)
		http.Error(
			w, fmt.Sprintf("Drink %d not found.", req.Drink), http.StatusBadRequest)
		return
	}

	machine := Machine{}
	db.First(&machine, req.Drink)

	if machine.ID == 0 {
		log.Printf("Brew request: unknown machine %d\n", req.Machine)
		http.Error(
			w, fmt.Sprintf("Machine %d not found.", req.Machine), http.StatusBadRequest)
		return
	}

	activeBackendMapLock.RLock()
	defer activeBackendMapLock.RUnlock()
	backend, ok := activeBackendMap[machine.ID]

	if !ok {
		log.Printf(
			"Ignoring a request to make drink %d (name: %s) on machine %d (name: %s) "+
				"because machine %d is not available.\n",
			drink.ID, drink.Name, machine.ID, machine.Name, machine.ID)
		return
	}
	backend <- &drink
	log.Printf(
		"Enqueued a request to make drink %d (name: %s) on machine %d "+
			"(name: %s).\n",
		drink.ID, drink.Name, machine.ID, machine.Name)
}
