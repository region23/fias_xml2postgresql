package house

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pavlik/fias_xml2postgresql/helpers"
)

const dateformat = "2006-01-02"

// Сведения по номерам домов улиц городов и населенных пунктов, номера земельных участков и т.п
type XmlObject struct {
	XMLName    xml.Name `xml:"House"`
	POSTALCODE *string  `xml:"POSTALCODE,attr,omitempty"`
	IFNSFL     int      `xml:"IFNSFL,attr,omitempty"`
	TERRIFNSFL int      `xml:"TERRIFNSFL,attr,omitempty"`
	IFNSUL     int      `xml:"IFNSUL,attr,omitempty"`
	TERRIFNSUL int      `xml:"TERRIFNSUL,attr,omitempty"`
	OKATO      *string  `xml:"OKATO,attr,omitempty"`
	OKTMO      *string  `xml:"OKTMO,attr,omitempty"`
	UPDATEDATE string   `xml:"UPDATEDATE,attr"`
	HOUSENUM   *string  `xml:"HOUSENUM,attr,omitempty"`
	ESTSTATUS  int      `xml:"ESTSTATUS,attr"`
	BUILDNUM   *string  `xml:"BUILDNUM,attr,omitempty"`
	STRUCNUM   *string  `xml:"STRUCNUM,attr,omitempty"`
	STRSTATUS  int      `xml:"STRSTATUS,attr"`
	HOUSEID    string   `xml:"HOUSEID,attr"`
	HOUSEGUID  string   `xml:"HOUSEGUID,attr"`
	AOGUID     string   `xml:"AOGUID,attr"`
	STARTDATE  string   `xml:"STARTDATE,attr"`
	ENDDATE    string   `xml:"ENDDATE,attr"`
	STATSTATUS int      `xml:"STATSTATUS,attr"`
	NORMDOC    *string  `xml:"NORMDOC,attr,omitempty"`
	COUNTER    int      `xml:"COUNTER,attr"`
}

// схема таблицы в БД

const tableName = "as_house_"
const elementName = "House"

const schema = `CREATE TABLE ` + tableName + ` (
    house_id UUID NOT NULL,
    postal_code VARCHAR(6),
		ifns_fl INT,
		terr_ifns_fl INT,
		ifns_ul INT,
		terr_ifns_ul INT,
		okato VARCHAR(11),
		oktmo VARCHAR(11),
		update_date TIMESTAMP NOT NULL,
		house_num VARCHAR(20),
		est_status INT NOT NULL,
		build_num VARCHAR(20),
		struc_num VARCHAR(20),
		str_status INT,
		house_guid UUID NOT NULL,
		ao_guid UUID NOT NULL,
		start_date TIMESTAMP NOT NULL,
		end_date TIMESTAMP NOT NULL,
		stat_status INT NOT NULL,
		norm_doc UUID,
		counter INT NOT NULL,
		PRIMARY KEY (house_id));`

func Export(c chan string, db *sqlx.DB, format *string) {
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
	//fmt.Println("Подсчет строк")
	// _, err := helpers.CountElementsInXML(pathToFile, elementName)
	// if err != nil {
	// 	fmt.Println("Error counting elements in XML file:", err)
	// 	return
	// }
	//fmt.Println("\nВ ", elementName, " содержится ", countedElements, " строк")

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

				//fmt.Println(item, "\n\n")

				var err error
				var updDate, startDate, endDate time.Time

				updDate, err = time.Parse(dateformat, item.UPDATEDATE)
				if err != nil {
					fmt.Println("Error parse UPDATEDATE: ", err)
					return
				}

				startDate, err = time.Parse(dateformat, item.STARTDATE)
				if err != nil {
					fmt.Println("Error parse STARTDATE: ", err)
					return
				}

				endDate, err = time.Parse(dateformat, item.ENDDATE)
				if err != nil {
					fmt.Println("Error parse ENDDATE: ", err)
					return
				}

				query := `INSERT INTO ` + tableName + ` (house_guid,
					postal_code,
					ifns_fl,
					terr_ifns_fl,
					ifns_ul,
					terr_ifns_ul,
					okato,
					oktmo,
					update_date,
					house_num,
					est_status,
					build_num,
					struc_num,
					str_status,
					house_id,
					ao_guid,
					start_date,
					end_date,
					stat_status,
					norm_doc,
					counter
					) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
						$11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
						$21)`

				db.MustExec(query,
					item.HOUSEGUID,
					item.POSTALCODE,
					item.IFNSFL,
					item.TERRIFNSFL,
					item.IFNSUL,
					item.TERRIFNSUL,
					item.OKATO,
					item.OKTMO,
					updDate,
					item.HOUSENUM,
					item.ESTSTATUS,
					item.BUILDNUM,
					item.STRUCNUM,
					item.STRSTATUS,
					item.HOUSEID,
					item.AOGUID,
					startDate,
					endDate,
					item.STATSTATUS,
					item.NORMDOC,
					item.COUNTER)

				s := strconv.Itoa(total)
				c <- elementName + " " + s + " rows affected"
			}
		default:
		}

	}

	//fmt.Printf("Total processed items in AddressObjects: %d \n", total)
}
