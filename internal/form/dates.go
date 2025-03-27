package form

import (
	"github.com/charmbracelet/huh"
)

func AskDates(departDate, returnDate *string) error {

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title("Выберите даты").Description("Автоматически формирует диапазоны дат от 3 дней до и 4 дней после выбранной даты — как для вылета, так и для возвращения."),
			huh.NewInput().Title("Примерная дата вылета (ГГГГ-ММ-ДД)").Value(departDate),
			huh.NewInput().Title("Примерная дата возвращения (ГГГГ-ММ-ДД)").Value(returnDate),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}
	return nil
}
