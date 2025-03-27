package form

import (
	"github.com/charmbracelet/huh"
)

func AskDates(departDate, returnDate *string) error {

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Примерная дата вылета (ГГГГ-ММ-ДД)").Value(departDate),
			huh.NewInput().Title("Примерная дата возвращения (ГГГГ-ММ-ДД)").Value(returnDate),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}
	return nil
}
