package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
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

type ScheduleCreateRequest struct {
	Name    string   `json:"name"`
	Days    []string `json:"days"`
	Time    string   `json:"time"`
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

	days, err := daysToBitvector(data.Days)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}

	t, err := time.Parse("3:04 PM", strings.ToUpper(data.Time))
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to parse time.", http.StatusBadRequest)
		return
	}
	from := time.Date(0, time.January, 0, 0, 0, 0, 0, time.UTC)
	tsec := int(t.Sub(from).Seconds())

	var drinks []Drink
	db.Where(&Drink{ID: data.Drink}).Find(&drinks)
	if len(drinks) != 1 {
		if len(drinks) < 1 {
			log.Printf("No drinks found with id %d\n", data.Drink)
			http.Error(
				w, fmt.Sprintf("No drinks found with id %d", data.Drink),
				http.StatusBadRequest)
		} else {
			log.Panicf("More than one drink with id %d!\n", data.Drink)
		}
		return
	}

	var machines []Machine
	db.Where(&Machine{ID: data.Machine}).Find(&machines)
	if len(machines) != 1 {
		if len(machines) < 1 {
			log.Printf("No machines found with id %d\n", data.Machine)
			http.Error(
				w, fmt.Sprintf("No machines found with id %d", data.Machine),
				http.StatusBadRequest)
		} else {
			log.Panicf("More than one machine with id %d!\n", data.Machine)
		}
		return
	}

	schedule := Schedule{
		Name:      data.Name,
		Days:      days,
		Enabled:   data.Enabled,
		Time:      tsec,
		DrinkID:   data.Drink,
		MachineID: data.Machine,
	}

	db.Create(&schedule)

	enc := json.NewEncoder(w)
	err = enc.Encode(schedule)
	if err != nil {
		log.Printf("Error while encoding: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}
