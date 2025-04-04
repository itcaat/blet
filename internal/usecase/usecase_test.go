package usecase

import (
	"testing"

	"github.com/itcaat/blet/internal/models"
	"github.com/stretchr/testify/require"
)

type mockAPIClient struct{}

func (m *mockAPIClient) GetSpecialOffers(origin string) (models.SpecialOffersResponse, error) {
	return models.SpecialOffersResponse{
		Success: true,
		Data: []models.SpecialOffers{
			{Destination: "LED", DepartDate: "2025-07-01", Price: 4999, Link: "/search/led"},
			{Destination: "KUF", DepartDate: "2025-07-05", Price: 3500, Link: "/search/kuf"},
		},
	}, nil
}

func (m *mockAPIClient) GetCheapest(origin, destination, oneWay string) (models.PriceForDatesResponse, error) {
	return models.PriceForDatesResponse{
		Success: true,
		Data: []models.Ticket{
			{Origin: origin, Destination: destination, Price: 1000, DepartureAt: "2025-06-01", ReturnAt: "2025-06-10", Link: "/search/cheapest"},
		},
		Currency: "RUB",
	}, nil
}

func (m *mockAPIClient) GetWeekPrices(origin, destination, depart, back string) (models.WeekMatrixResponse, error) {
	return models.WeekMatrixResponse{
		Success: true,
		Data: []models.WeekMatrixFlight{
			{Destination: destination, DepartDate: depart, ReturnDate: back, Value: 1500, NumberOfStops: 1},
		},
	}, nil
}

func (m *mockAPIClient) GetShortUrl(url string) (models.ShortLinksResponse, error) {
	return models.ShortLinksResponse{
		Status: "success",
		Result: models.ShortLinksResult{
			Links: []models.ShortLink{{Url: url, PartnerUrl: "https://short.url/mock"}},
		},
	}, nil
}

func TestGetSpecialOffers(t *testing.T) {
	client := &mockAPIClient{}
	offers, err := GetSpecialOffers(client, "MOW")
	require.NoError(t, err)
	require.Len(t, offers, 2)
	require.Equal(t, "LED", offers[0].Destination)
}

func TestGetCheapestTickets(t *testing.T) {
	client := &mockAPIClient{}
	tickets, err := GetCheapestTickets(client, "MOW", "LED", "true")
	require.NoError(t, err)
	require.Len(t, tickets, 1)
	require.Equal(t, "LED", tickets[0].Destination)
	require.Equal(t, 1000, tickets[0].Price)
}

func TestGetWeekMatrix(t *testing.T) {
	client := &mockAPIClient{}
	flights, err := GetWeekMatrix(client, "MOW", "LED", "2025-06-01", "2025-06-10")
	require.NoError(t, err)
	require.Len(t, flights, 1)
	require.Equal(t, "LED", flights[0].Destination)
	require.Equal(t, 1500, flights[0].Value)
}
