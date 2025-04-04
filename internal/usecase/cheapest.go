package usecase

import (
	"github.com/itcaat/blet/internal/api"
	"github.com/itcaat/blet/internal/models"
)

func GetCheapestTickets(client api.TravelpayoutsAPI, origin, destination, one_way string) ([]models.Ticket, error) {
	resp, err := client.GetCheapest(origin, destination, one_way)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}
