package usecase

import (
	"github.com/itcaat/blet/internal/api"
	"github.com/itcaat/blet/internal/models"
)

func GetWeekMatrix(origin, destination, token string) ([]models.WeekMatrixFlight, error) {
	resp, err := api.GetWeekPrices(origin, destination, token)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}
