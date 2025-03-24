package usecase

import (
	"github.com/itcaat/blet/internal/api"
	"github.com/itcaat/blet/internal/models"
)

func GetWeekMatrix(origin, destination, depart, back, token string) ([]models.WeekMatrixFlight, error) {
	resp, err := api.GetWeekPrices(origin, destination, depart, back, token)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}
