package form

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
)

func ShowConfirm(confirm *bool, title, affirmative, negative string) {

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(title).
				Affirmative(affirmative).
				Negative(negative).
				Value(confirm),
		),
	)

	if err := form.Run(); err != nil {
		fmt.Println("❌ Ошибка выбора:", err)
		os.Exit(1)
	}

}
