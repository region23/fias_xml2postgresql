package helpers

import (
	"encoding/xml"
	"fmt"
	"os"
)

// CountElementsInXML возвращает количество узлов в XML-файле
func CountElementsInXML(pathToFile string, countedElement string) (int, error) {
	var err error

	xmlFile, err := os.Open(pathToFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return 0, err
	}

	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	total := 0
	var inElement string
	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		// Inspect the type of the token just read.
		switch se := t.(type) {
		case xml.StartElement:
			// If we just read a StartElement token
			inElement = se.Name.Local

			if inElement == countedElement {
				total++
			}
		default:
		}

	}

	return total, nil

}
