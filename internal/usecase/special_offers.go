package usecase

import (
	"github.com/itcaat/blet/internal/api"
	"github.com/itcaat/blet/internal/models"
)

func GetSpecialOffers(client api.TravelpayoutsAPI, origin string) ([]models.Ticket, error) {
	resp, err := client.GetSpecialOffers(origin)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}
