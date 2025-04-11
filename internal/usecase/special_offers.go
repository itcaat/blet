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

	// Convert resp.Data to []*models.Ticket
	tickets := make([]*models.Ticket, len(resp.Data))
	for i := range resp.Data {
		tickets[i] = &resp.Data[i]
	}
	if err = client.GetShortUrlArray(tickets); err != nil {
		return nil, err
	}
	return resp.Data, nil

}
