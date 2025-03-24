package form

import (
	"github.com/charmbracelet/huh"
)

func AskDates() (string, string, error) {
	var departDate, returnDate string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Примерная дата вылета (ГГГГ-ММ-ДД)").Value(&departDate),
			huh.NewInput().Title("Примерная дата возвращения (ГГГГ-ММ-ДД)").Value(&returnDate),
		),
	)

	if err := form.Run(); err != nil {
		return "", "", err
	}
	return departDate, returnDate, nil
}
