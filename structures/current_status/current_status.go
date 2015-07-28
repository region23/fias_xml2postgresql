package current_status

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pavlik/fias_xml2postgresql/helpers"
)

const dateformat = "2006-01-02"

// Статус актуальности КЛАДР 4.0
type XmlObject struct {
	XMLName    xml.Name `xml:"CurrentStatus"`
	CURENTSTID int      `xml:"CURENTSTID,attr"`
	NAME       string   `xml:"NAME,attr"`
}

// схема таблицы в БД

const tableName = "as_curentst"
const elementName = "CurrentStatus"

const schema = `CREATE TABLE ` + tableName + ` (
    curent_st_id INT UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
		PRIMARY KEY (curent_st_id));`

func Export(w *sync.WaitGroup, c chan string, db *sqlx.DB, format *string) {
	w.Add(1)
	defer w.Done()
	helpers.DropAndCreateTable(schema, tableName, db)

	var format2 string
	format2 = *format
	fileName, err2 := helpers.SearchFile(tableName, format2)
	if err2 != nil {
		fmt.Println("Error searching file:", err2)
		return
	}

	pathToFile := format2 + "/" + fileName

	// Подсчитываем, сколько элементов нужно обработать
	//_, err := helpers.CountElementsInXML(pathToFile, elementName)
	// if err != nil {
	// 	fmt.Println("Error counting elements in XML file:", err)
	// 	return
	// }
	// fmt.Println("\nВ ", elementName, " содержится ", countedElements, " строк")

	xmlFile, err := os.Open(pathToFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
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

			if inElement == elementName {
				total++
				var item XmlObject

				// decode a whole chunk of following XML into the
				// variable item which is a ActualStatus (se above)
				err = decoder.DecodeElement(&item, &se)
				if err != nil {
					fmt.Println("Error in decode element:", err)
					return
				}
				query := "INSERT INTO " + tableName + " (curent_st_id, name) VALUES ($1, $2)"
				db.MustExec(query, item.CURENTSTID, item.NAME)

				s := strconv.Itoa(total)

				c <- elementName + " " + s + " rows affected"

			}
		default:
		}

	}

	//fmt.Printf("Total processed items in CurrentStatus: %d \n", total)
}
