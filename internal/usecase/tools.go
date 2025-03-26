package usecase

import (
	"github.com/itcaat/blet/internal/api"
	"github.com/itcaat/blet/internal/models"
)

func GetShortUrl(url, token string) ([]models.ShortLink, error) {
	resp, err := api.GetShortUrl(url, token)
	if err != nil {
		return nil, err
	}
	return resp.Result.Links, nil
}
