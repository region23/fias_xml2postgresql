package normative_document

import (
	"encoding/xml"
	"log"
	"os"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pavlik/fias_xml2postgresql/helpers"
)

const dateformat = "2006-01-02"
const tableName = "as_normdoc"
const elementName = "NormativeDocument"

// Сведения по нормативному документу,
// являющемуся основанием присвоения адресному элементу наименования
type XmlObject struct {
	XMLName   xml.Name `xml:"NormativeDocument"`
	NORMDOCID string   `xml:"NORMDOCID,attr"`
	DOCNAME   *string  `xml:"DOCNAME,attr,omitempty"`
	DOCDATE   *string  `xml:"DOCDATE,attr,omitempty"`
	DOCNUM    *string  `xml:"DOCNUM,attr,omitempty"`
	DOCTYPE   int      `xml:"DOCTYPE,attr"`
	DOCIMGID  int      `xml:"DOCIMGID,attr,omitempty"`
}

const schema = `CREATE TABLE ` + tableName + ` (
    norm_doc_id UUID NOT NULL,
    doc_name VARCHAR(1000),
    doc_date TIMESTAMP,
    doc_num VARCHAR(20),
    doc_type INT NOT NULL,
    doc_img_id INT,
		PRIMARY KEY (norm_doc_id));`

func Export(w *sync.WaitGroup, c chan string, db *sqlx.DB, format *string) {
	w.Add(1)
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

	// Подсчитываем, сколько элементов нужно обработать
	//fmt.Println("Подсчет строк")
	// _, err := helpers.CountElementsInXML(pathToFile, elementName)
	// if err != nil {
	// 	fmt.Println("Error counting elements in XML file:", err)
	// 	return
	// }
	//fmt.Println("\nВ ", elementName, " содержится ", countedElements, " строк")

	xmlFile, err := os.Open(pathToFile)
	if err != nil {
		log.Println("Error opening file:", err)
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
					log.Println("Error in decode element:", err)
					return
				}

				var err error

				query := `INSERT INTO ` + tableName + ` (norm_doc_id,
					doc_name,
          doc_date,
          doc_num,
          doc_type,
          doc_img_id
					) VALUES ($1, $2, $3, $4, $5, $6)`

				_, err = db.Exec(query,
					item.NORMDOCID,
					item.DOCNAME,
					item.DOCDATE,
					item.DOCNUM,
					item.DOCTYPE,
					item.DOCIMGID)

				if err != nil {
					log.Fatal(err)
				}

				c <- helpers.PrintRowsAffected(elementName, total)
			}
		default:
		}

	}
}
