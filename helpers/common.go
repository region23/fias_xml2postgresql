package helpers

import (
	"bytes"
	"fmt"
	"strconv"
)

func concat(values ...string) string {
	var buffer bytes.Buffer
	for _, s := range values {
		buffer.WriteString(s)
	}
	return buffer.String()
}

func humanizeInt(n int64) string {
	var pretty string

	ns := strconv.FormatInt(n, 10)
	nsl := len(ns)

	for i := nsl; i > 0; i-- {
		if (nsl-i)%3 == 0 {
			pretty = concat(" ", pretty)
		}
		pretty = concat(string(ns[i-1]), pretty)
	}

	return pretty
}

func PrintRowsAffected(elementName string, currentRow int) string {
	return fmt.Sprintf("Таблица %s. Обработано %d строк.", elementName, humanizeInt(currentRow))
}
