package graph

import (
	"time"
)

type Crime struct {
	Date      time.Time `json:"date"`
	ID        uint64    `json:"id"`
	Lat       float64   `json:"lat"`
	Lng       float64   `json:"lng"`
	Type      string    `json:"type"`
	Transport string    `json:"transport"`
	Weapon    string    `json:"weapon"`
	Victim    Victim    `json:"victim"`
}

type Victim struct {
	Sex string `json:"sex"`
	Age int    `json:"age"`
}

type Crimes []Crime

type CrimeInMemoryProvider interface {
	Fetch() map[uint64]Crime
}
