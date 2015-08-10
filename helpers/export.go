package helpers

import (
	"encoding/xml"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type xmlObjectName struct {
	tableName   string
	elementName string
}

func extractXMLObjectName(xmlObject interface{}) xmlObjectName {
	v := reflect.ValueOf(xmlObject).Elem()
	field := v.Type().Field(0)
	obj := xmlObjectName{tableName: field.Tag.Get("db"), elementName: field.Tag.Get("xml")}
	return obj
}

func extractFeilds(xmlObject interface{}) []string {
	v := reflect.ValueOf(xmlObject).Elem()
	fields := make([]string, v.NumField()-1)

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		if field.Type.String() != "xml.Name" {
			fields[i-1] = field.Tag.Get("db")
		}
	}
	return fields
}

func extractValues(xmlObject interface{}) []interface{} {
	//fmt.Println(xmlObject)
	s := reflect.ValueOf(xmlObject).Elem()
	values := make([]interface{}, s.NumField()-1)

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if f.Type().Name() != "xml.Name" {
			if f.Kind() == reflect.String {
				values[i-1] = f.String()
			} else if f.Kind() == reflect.Int {
				values[i-1] = f.Int()
			} else if f.Kind() == reflect.Bool {
				values[i-1] = f.Bool()
			}
		}
	}

	return values
}

// ExportBulk экспортирует данные из xml-файла в таблицу указанную в описании xml-структуры
func ExportBulk(schema func(tableName string) string, xmlObject interface{}, w *sync.WaitGroup, c chan string, db *sqlx.DB, format *string, logger *log.Logger) {

	defer w.Done()

	const dateformat = "2006-01-02"

	objName := extractXMLObjectName(xmlObject)
	fields := extractFeilds(xmlObject)
	searchFileName := objName.tableName
	objName.tableName = strings.TrimSuffix(objName.tableName, "_")

	DropAndCreateTable(schema(objName.tableName), objName.tableName, db)

	var format2 string
	format2 = *format
	fileName, err2 := SearchFile(searchFileName, format2)
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

	query := pq.CopyIn(objName.tableName, fields...)

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

			txn, err = db.Begin()
			if err != nil {
				logger.Fatalln(err)
			}

			stmt, err = txn.Prepare(query)
			if err != nil {
				logger.Fatalln(err)
			}
		}

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

			if inElement == objName.elementName {
				total++
				//var item xmlObject

				// decode a whole chunk of following XML into the
				// variable item which is a ActualStatus (se above)
				err = decoder.DecodeElement(&xmlObject, &se)
				if err != nil {
					logger.Fatalln("Error in decode element:", err)
				}

				values := extractValues(xmlObject)
				//logger.Println(values)
				_, err = stmt.Exec(values...)

				if err != nil {
					logger.Fatalln(err)
				}
				c <- PrintRowsAffected(objName.elementName, total)
				i++
			}
		default:
		}

	}

	// Если закончили обработку всех записей, ставим вначале строки восклицательный знак.
	// Признак того, что все записи в данном файле обработаны
	c <- "!" + PrintRowsAffected(objName.elementName, total)
}
