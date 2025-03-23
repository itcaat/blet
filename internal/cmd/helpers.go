package cmd

import "fmt"

func FormatAviasalesLink(path string) string {
	return fmt.Sprintf("https://www.aviasales.ru%s", path)
}
