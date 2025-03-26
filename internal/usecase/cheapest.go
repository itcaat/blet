package usecase

import (
	"github.com/itcaat/blet/internal/api"
	"github.com/itcaat/blet/internal/models"
)

func GetCheapestTickets(origin, destination, token string) ([]models.Ticket, error) {
	resp, err := api.GetCheapest(origin, destination, token)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}
