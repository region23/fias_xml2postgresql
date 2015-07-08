package structures

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/go-gorp/gorp"
	_ "github.com/lib/pq"
)

// Статус актуальности ФИАС
type XmlObject struct {
	XMLName   xml.Name `xml:"ActualStatus"`
	ActStatId int      `xml:"ACTSTATID,attr"`
	Name      string   `xml:"NAME,attr"`
}

type DBObject struct {
	ActStatId int    `db:"actstat_id, primarykey"`
	Name      string `db:"name"`
}

func xml2db(xml XmlObject) *DBObject {
	obj := &DBObject{
		Name:      xml.Name,
		ActStatId: xml.ActStatId}
	return obj
}

func (item XmlObject) String() string {
	return fmt.Sprintf("\t ActStatId : %d - Name : %s \n", item.ActStatId, item.Name)
}

func ExportActualStatus(dbmap *gorp.DbMap) {
	// Создаем таблицу
	dbmap.AddTableWithName(DBObject{}, "actstat")
	err := dbmap.DropTableIfExists(DBObject{})
	if err != nil {
		fmt.Println("Error on drop table:", err)
		return
	}
	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		fmt.Println("Error on creating table:", err)
		return
	}

	xmlFile, err := os.Open("xml/AS_ACTSTAT_20150705_c9027b5f-3370-4705-be8a-fa06793614ee.XML")
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
			// ...and its name is "ActualStatus"
			if inElement == "ActualStatus" {
				total++
				var item XmlObject
				// decode a whole chunk of following XML into the
				// variable item which is a ActualStatus (se above)
				decoder.DecodeElement(&item, &se)
				obj := xml2db(item)
				err := dbmap.Insert(obj)
				if err != nil {
					fmt.Println("Error on creating table:", err)
					return
				}

				fmt.Printf(item.String())
			}
		default:
		}

	}

	fmt.Printf("Total items in ActualStatus: %d \n", total)
}
