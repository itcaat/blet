package usecase

import (
	"github.com/itcaat/blet/internal/api"
	"github.com/itcaat/blet/internal/models"
)

func GetSpecialOffers(origin, token string) ([]models.SpecialOffers, error) {
	resp, err := api.GetSpecialOffers(origin, token)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}
