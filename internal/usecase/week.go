package usecase

import (
	"github.com/itcaat/blet/internal/api"
	"github.com/itcaat/blet/internal/models"
)

func GetWeekMatrix(client api.TravelpayoutsAPI, origin, destination, depart, back string) ([]models.WeekMatrixFlight, error) {
	resp, err := client.GetWeekPrices(origin, destination, depart, back)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}
