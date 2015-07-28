package helpers

import (
	"encoding/xml"
	"log"
	"os"
	"sync"
)

// CountElementsInXML возвращает количество узлов в XML-файле
func CountElementsInXML(w *sync.WaitGroup, c chan int, tableName string, countedElement string) {
	w.Add(1)
	defer w.Done()

	var err error

	format := "xml"

	fileName, err2 := SearchFile(tableName, format)
	if err2 != nil {
		log.Fatalln("Error searching file:", err2)
		panic(err)
	}

	pathToFile := format + "/" + fileName

	xmlFile, err := os.Open(pathToFile)
	if err != nil {
		log.Fatalln("Error opening file:", err)
		panic(err)
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
				c <- total
			}
		default:
		}
	}
}
