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

type WeekMatrixFlight struct {
	Destination     string `json:"destination"`
	DepartDate      string `json:"depart_date"`
	ReturnDate      string `json:"return_date"`
	Value           int    `json:"value"`
	NumberOfChanges int    `json:"number_of_changes"`
	NumberOfStops   int    `json:"number_of_stops"`
}

type WeekMatrixResponse struct {
	Success bool               `json:"success"`
	Data    []WeekMatrixFlight `json:"data"`
}

type SpecialOffers struct {
	Destination string `json:"destination"`
	DepartDate  string `json:"departure_at"`
	Price       int    `json:"price"`
	Link        string `json:"link"`
}

func (t *SpecialOffers) URL() string {
	base := "https://www.aviasales.ru"
	return fmt.Sprintf("%s%s", base, t.Link)
}

type SpecialOffersResponse struct {
	Success bool            `json:"success"`
	Data    []SpecialOffers `json:"data"`
}
