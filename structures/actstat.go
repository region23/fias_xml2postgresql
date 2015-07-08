package structures

import (
	"encoding/xml"
	"fmt"
	"os"
)

// Статус актуальности ФИАС
type ActualStatus struct {
	XMLName   xml.Name `xml:"ActualStatus"`
	ActStatId int      `xml:"ACTSTATID,attr"`
	Name      string   `xml:"NAME,attr"`
}

// type ActualStatuses struct {
// 	XMLName        xml.Name       `xml:"ActualStatuses"`
// 	ActualStatuses []ActualStatus `xml:"ActualStatus"`
// }

func (item ActualStatus) String() string {
	return fmt.Sprintf("\t ActStatId : %d - Name : %s \n", item.ActStatId, item.Name)
}

func ExportActualStatus() {
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
				var item ActualStatus
				// decode a whole chunk of following XML into the
				// variable item which is a ActualStatus (se above)
				decoder.DecodeElement(&item, &se)
				fmt.Printf(item)
			}
		default:
		}

	}

	//fmt.Printf("Total articles: %d \n", total)
}
