package normative_document_type

import (
	"encoding/xml"
	"log"
	"os"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pavlik/fias_xml2postgresql/helpers"
)

func ExportBulk(w *sync.WaitGroup, c chan string, db *sqlx.DB, format *string, logger *log.Logger) {

	defer w.Done()
	// make sure log.txt exists first
	// use touch command to create if log.txt does not exist
	var logFile *os.File
	var err error
	if _, err1 := os.Stat("log.txt"); err1 == nil {
		logFile, err = os.OpenFile("log.txt", os.O_WRONLY, 0666)
	} else {
		logFile, err = os.Create("log.txt")
	}

	if err != nil {
		panic(err)
	}

	defer logFile.Close()

	log.SetFlags(log.Llongfile)
	// direct all log messages to log.txt
	log.SetOutput(logFile)

	helpers.DropAndCreateTable(schema, tableName, db)

	var format2 string
	format2 = *format
	fileName, err2 := helpers.SearchFile(tableName, format2)
	if err2 != nil {
		log.Println("Error searching file:", err2)
		return
	}

	pathToFile := format2 + "/" + fileName

	xmlFile, err := os.Open(pathToFile)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}

	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	var inElement string
	total := 0
	i := 0

	txn, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	query := pq.CopyIn(tableName, "nd_type_id", "name")

	stmt, err := txn.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}

	for {
		if i == 50000 {
			i = 0

			_, err = stmt.Exec()
			if err != nil {
				log.Fatal(err)
			}

			err = stmt.Close()
			if err != nil {
				log.Fatal(err)
			}

			err = txn.Commit()
			if err != nil {
				log.Fatal(err)
			}

			//c <- helpers.PrintRowsAffected(elementName, total)

			txn, err = db.Begin()
			if err != nil {
				log.Fatal(err)
			}

			stmt, err = txn.Prepare(query)
			if err != nil {
				log.Fatal(err)
			}
		}
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()

		// Если достигли конца xml-файла
		if t == nil {
			if i > 0 {
				_, err = stmt.Exec()
				if err != nil {
					log.Fatal(err)
				}

				err = stmt.Close()
				if err != nil {
					log.Fatal(err)
				}

				err = txn.Commit()
				if err != nil {
					log.Fatal(err)
				}
			}

			//c <- helpers.PrintRowsAffected(elementName, total)

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
					log.Println("Error in decode element:", err)
					return
				}

				_, err = stmt.Exec(item.NDTYPEID, item.NAME)

				if err != nil {
					log.Fatal(err)
				}
				c <- helpers.PrintRowsAffected(elementName, total)
				i++
			}
		default:
		}

	}
}
