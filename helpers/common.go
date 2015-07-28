package helpers

import "fmt"

func PrintRowsAffected(elementName string, currentRow int) string {
	return fmt.Sprintf("Таблица %s. Обработано %d строк.", elementName, currentRow)
}
