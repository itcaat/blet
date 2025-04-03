package models

import "fmt"

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
