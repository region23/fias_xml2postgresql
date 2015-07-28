package address_object_type

import (
	"encoding/xml"
	"fmt"
	"os"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pavlik/fias_xml2postgresql/helpers"
)

const dateformat = "2006-01-02"

// Статус действия
type XmlObject struct {
	XMLName  xml.Name `xml:"AddressObjectType"`
	LEVEL    int      `xml:"LEVEL,attr"`
	SCNAME   string   `xml:"SCNAME,attr"`
	SOCRNAME string   `xml:"SOCRNAME,attr"`
	KOD_T_ST string   `xml:"KOD_T_ST,attr"`
}

// схема таблицы в БД

const tableName = "as_socrbase"
const elementName = "AddressObjectType"

const schema = `CREATE TABLE ` + tableName + ` (
    level INT NOT NULL,
    sc_name VARCHAR(20),
    socr_name VARCHAR(60),
    kod_t_st INT UNIQUE NOT NULL,
		PRIMARY KEY (kod_t_st));`

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
				query := `INSERT INTO ` + tableName + ` (
          level,
          sc_name,
          socr_name,
          kod_t_st)
          VALUES (
            $1, $2, $3, $4)`

				db.MustExec(query,
					item.LEVEL,
					item.SCNAME,
					item.SOCRNAME,
					item.KOD_T_ST)

				c <- helpers.PrintRowsAffected(elementName, total)
			}
		default:
		}

	}
}
