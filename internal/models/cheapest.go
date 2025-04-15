package models

type PriceForDatesResponse struct {
	Success  bool     `json:"success"`
	Data     []Ticket `json:"data"`
	Currency string   `json:"currency"`
}
