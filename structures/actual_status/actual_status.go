package actual_status

import (
	"encoding/xml"

	_ "github.com/lib/pq"
)

// Статус актуальности ФИАС
type XmlObject struct {
	XMLName   xml.Name `xml:"ActualStatus" db:"as_actstat"`
	ActStatId int      `xml:"ACTSTATID,attr" db:"act_stat_id"`
	Name      string   `xml:"NAME,attr" db:"name"`
}

// схема таблицы в БД

// const tableName = "as_actstat"
// const elementName = "ActualStatus"

func Schema(tableName string) string {
	return `CREATE TABLE ` + tableName + ` (
    act_stat_id INT UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
		PRIMARY KEY (act_stat_id));`
}

// func (item XmlObject) String() string {
// 	return fmt.Sprintf("\t ActStatId : %d - Name : %s \n", item.ActStatId, item.Name)
// }

//
// func Export(w *sync.WaitGroup, c chan string, db *sqlx.DB, format *string) {
//
// 	defer w.Done()
//
// 	helpers.DropAndCreateTable(schema, tableName, db)
//
// 	var format2 string
// 	format2 = *format
// 	fileName, err2 := helpers.SearchFile(tableName, format2)
// 	if err2 != nil {
// 		fmt.Println("Error searching file:", err2)
// 		return
// 	}
//
// 	pathToFile := format2 + "/" + fileName
//
// 	// Подсчитываем, сколько элементов нужно обработать
// 	//_, err := helpers.CountElementsInXML(pathToFile, elementName)
// 	// if err != nil {
// 	// 	fmt.Println("Error counting elements in XML file:", err)
// 	// 	return
// 	// }
// 	// fmt.Println("\nВ ", elementName, " содержится ", countedElements, " строк")
//
// 	xmlFile, err := os.Open(pathToFile)
// 	if err != nil {
// 		fmt.Println("Error opening file:", err)
// 		return
// 	}
//
// 	defer xmlFile.Close()
//
// 	decoder := xml.NewDecoder(xmlFile)
// 	total := 0
// 	var inElement string
// 	for {
// 		// Read tokens from the XML document in a stream.
// 		t, _ := decoder.Token()
// 		if t == nil {
// 			break
// 		}
// 		// Inspect the type of the token just read.
// 		switch se := t.(type) {
// 		case xml.StartElement:
// 			// If we just read a StartElement token
// 			inElement = se.Name.Local
//
// 			if inElement == elementName {
// 				total++
// 				var item XmlObject
// 				// decode a whole chunk of following XML into the
// 				// variable item which is a ActualStatus (se above)
// 				err = decoder.DecodeElement(&item, &se)
// 				if err != nil {
// 					fmt.Println("Error in decode element:", err)
// 					return
// 				}
// 				query := "INSERT INTO " + tableName + " (act_stat_id, name) VALUES ($1, $2)"
// 				db.MustExec(query, item.ActStatId, item.Name)
//
// 				c <- helpers.PrintRowsAffected(elementName, total)
// 			}
// 		default:
// 		}
//
// 	}
//
// 	//fmt.Printf("\nВсего в "+elementName+" обработано %d строк\n", total)
// }
