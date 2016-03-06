package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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

	schedule.Name = data.Name[:100]
	schedule.Days = days
	schedule.Enabled = data.Enabled
	schedule.Time = data.Time
	schedule.DrinkID = data.Drink
	schedule.MachineID = data.Machine

	db.Save(&schedule)

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

	drink.Name = data.Name[:100]
	drink.Size = data.Size
	drink.Sugar = data.Sugar
	drink.Creamer = data.Creamer
	drink.TeaBag = data.TeaBag[:100]
	drink.KCup = data.KCup[:100]

	db.Save(&drink)

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
