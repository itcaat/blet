package models

import "fmt"

type Ticket struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	Price       int    `json:"price"`
	DepartureAt string `json:"departure_at"`
	Link        string `json:"link"`
}

func (t *Ticket) URL() string {
	base := "https://www.aviasales.ru"
	return fmt.Sprintf("%s%s", base, t.Link)
}

type PriceForDatesResponse struct {
	Success  bool     `json:"success"`
	Data     []Ticket `json:"data"`
	Currency string   `json:"currency"`
}
