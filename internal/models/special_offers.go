package models

type SpecialOffersResponse struct {
	Success bool     `json:"success"`
	Data    []Ticket `json:"data"`
}
