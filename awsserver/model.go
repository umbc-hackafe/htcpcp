package main

const (
	SUNDAY    uint8 = 1 << 0
	MONDAY    uint8 = 1 << 1
	TUESDAY   uint8 = 1 << 2
	WEDNESDAY uint8 = 1 << 3
	THURSDAY  uint8 = 1 << 4
	FRIDAY    uint8 = 1 << 5
	SATURDAY  uint8 = 1 << 6
)

var DaysMap = map[string]uint8{
	"SUNDAY":    SUNDAY,
	"MONDAY":    MONDAY,
	"TUESDAY":   TUESDAY,
	"WEDNESDAY": WEDNESDAY,
	"THURSDAY":  THURSDAY,
	"FRIDAY":    FRIDAY,
	"SATURDAY":  SATURDAY,
}

var DaysReverseMap = map[uint8]string{
	SUNDAY:    "SUNDAY",
	MONDAY:    "MONDAY",
	TUESDAY:   "TUESDAY",
	WEDNESDAY: "WEDNESDAY",
	THURSDAY:  "THURSDAY",
	FRIDAY:    "FRIDAY",
	SATURDAY:  "SATURDAY",
}

type Schedule struct {
	ID        uint   `gorm:"primary_key"`
	Name      string `sql:"not null,type:varchar(100)"`
	Days      uint8  `sql:"not null"`
	Enabled   bool   `sql:"not null"`
	Time      int    `sql:"not null"`
	DrinkID   uint
	Drink     Drink
	MachineID uint
	Machine   Machine
}

type Drink struct {
	ID      uint   `gorm:"primary_key"`
	Name    string `sql:"not null,type:varchar(100)"`
	Size    uint8  `sql:"not null"`
	Sugar   uint8  `sql:"not null"`
	Creamer uint8  `sql:"not null"`
	TeaBag  string `sql:"not null,type:varchar(100)"`
	KCup    string `sql:"not null,typ:varchar(100)"`
}

type Machine struct {
	ID   uint   `gorm:"primary_key",json:"id"`
	Name string `sql:"not null,type:varchar(100)",json:"name"`
}
