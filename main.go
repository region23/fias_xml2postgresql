package main

import (
	"encoding/xml"
	"fmt"
	//"io/ioutil"
	_ "github.com/pavlik/fias_xml2postgres/structures"
	"math"
	"os"
)

// Классификатор адресообразующих элементов
// type Object struct {
//        XMLName xml.Name `xml:"Object"`
//        AOGUID string `xml:"AOGUID, attr"`
//        FORMALNAME string `xml:"FORMALNAME, attr"`
// }
//

func (actstat ActualStatus) String() string {
	return fmt.Sprintf("\t ActStatId : %d - Name : %s \n", actstat.ActStatId, actstat.Name)
}

func main() {
	xmlFile, err := os.Open("ACTSTAT.xml")
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
			// ...and its name is "Object"
			if inElement == "ActualStatus" {
				var as ActualStatus
				// decode a whole chunk of following XML into the
				// variable p which is a Page (se above)
				decoder.DecodeElement(&as, &se)
				//obj := xml2db(p)
				// fmt.Println(obj)
				//dbmap.Insert(obj)
				total++
				if math.Mod(float64(total), 1000) == 0 {
					fmt.Println(total)
				}
				// fmt.Println(p)

				// Do some stuff with the page.
				// p.Title = CanonicalizeTitle(p.Title)
				// m := filter.MatchString(p.Title)
			}
		default:
		}
	}

	fmt.Printf("Total objects: %d \n", total)
	// XMLdata, _ := ioutil.ReadAll(xmlFile)

	// var actstats ActualStatuses
	// xml.Unmarshal(XMLdata, &actstats)

	//fmt.Println(actstats.ActualStatuses[1])
}
