package normative_document

/*
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

	helpers.DropAndCreateTable(schema, tableName, db)

	var format2 string
	format2 = *format
	fileName, err2 := helpers.SearchFile(tableName, format2)
	if err2 != nil {
		logger.Fatalln("Error searching file:", err2)
	}

	pathToFile := format2 + "/" + fileName

	xmlFile, err := os.Open(pathToFile)
	if err != nil {
		logger.Fatalln("Error opening file:", err)
	}

	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	var inElement string
	total := 0
	i := 0

	txn, err := db.Begin()
	if err != nil {
		logger.Fatalln(err)
	}

	query := pq.CopyIn(tableName,
		"norm_doc_id",
		"doc_name",
		"doc_date",
		"doc_num",
		"doc_type",
		"doc_img_id")

	stmt, err := txn.Prepare(query)
	if err != nil {
		logger.Fatalln(err)
	}

	for {
		if i == 50000 {
			i = 0

			_, err = stmt.Exec()
			if err != nil {
				logger.Fatalln(err)
			}

			err = stmt.Close()
			if err != nil {
				logger.Fatalln(err)
			}

			err = txn.Commit()
			if err != nil {
				logger.Fatalln(err)
			}

			//c <- helpers.PrintRowsAffected(elementName, total)

			txn, err = db.Begin()
			if err != nil {
				logger.Fatalln(err)
			}

			stmt, err = txn.Prepare(query)
			if err != nil {
				logger.Fatalln(err)
			}
		}
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()

		// Если достигли конца xml-файла
		if t == nil {
			if i > 0 {
				_, err = stmt.Exec()
				if err != nil {
					logger.Fatalln(err)
				}

				err = stmt.Close()
				if err != nil {
					logger.Fatalln(err)
				}

				err = txn.Commit()
				if err != nil {
					logger.Fatalln(err)
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
					logger.Fatalln("Error in decode element:", err)
				}

				_, err = stmt.Exec(item.NORMDOCID,
					item.DOCNAME,
					item.DOCDATE,
					item.DOCNUM,
					item.DOCTYPE,
					item.DOCIMGID)

				if err != nil {
					logger.Fatalln(err)
				}
				c <- helpers.PrintRowsAffected(elementName, total)
				i++
			}
		default:
		}

	}
}
*/
