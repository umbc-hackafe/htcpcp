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

type Schedule struct {
	ID      int
	Name    string `sql:"type:varchar(100)"`
	Days    uint8  `sql:"not null"`
	Enabled bool   `sql:"not null"`
	Time    int    `sql:"not null"`
	Drink   Drink
}

type Drink struct {
	ID int
}
