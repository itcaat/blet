package api

import "github.com/itcaat/blet/internal/models"

type TravelpayoutsAPI interface {
	GetSpecialOffers(origin string) (models.SpecialOffersResponse, error)
	GetCheapest(origin, destination, oneWay string) (models.PriceForDatesResponse, error)
	GetWeekPrices(origin, destination, depart, back string) (models.WeekMatrixResponse, error)
	GetShortUrl(url string) (models.ShortLinksResponse, error)
}
