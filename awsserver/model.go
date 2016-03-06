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

type Schedule struct {
	ID        uint    `gorm:"primary_key",json:"id"`
	Name      string  `sql:"not null,type:varchar(100)",json:"name"`
	Days      uint8   `sql:"not null",json:"days"`
	Enabled   bool    `sql:"not null",json:"enabled"`
	Time      int     `sql:"not null",json:"time"`
	DrinkID   uint    `json:"drink_id"`
	Drink     Drink   `sql:"not null",json:"-"`
	MachineID uint    `json:"machine_id"`
	Machine   Machine `sql:"not null",json:"-"`
}

type Drink struct {
	ID          uint   `gorm:"primary_key"`
	Name        string `sql:"not null,type:varchar(100)"`
	Description string `sql:"type:varcahr(2048)"`
	Size        int    `sql:"not null"`
}

type Machine struct {
	ID        uint   `gorm:"primary_key"`
	Name      string `sql:"not null,type:varchar(100)"`
	Connected bool   `sql:"not null"`
}
