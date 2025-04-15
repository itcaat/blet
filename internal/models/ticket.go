package models

import "fmt"

type Ticket struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	Price       int    `json:"price"`
	DepartureAt string `json:"departure_at"`
	ReturnAt    string `json:"return_at"`
	Link        string `json:"link"`
	ShortUrl    string
}

func (t *Ticket) URL() string {
	base := "https://www.aviasales.ru"
	return fmt.Sprintf("%s%s", base, t.Link)
}
