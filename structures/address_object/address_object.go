package address_object

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
const tableName = "addrobj"
const elementName = "Object"

// Классификатор адресообразующих элементов
type XmlObject struct {
	XMLName    xml.Name `xml:"Object"`
	AOGUID     string   `xml:"AOGUID,attr"`
	FORMALNAME string   `xml:"FORMALNAME,attr"`
	REGIONCODE string   `xml:"REGIONCODE,attr"`
	AUTOCODE   string   `xml:"AUTOCODE,attr"`
	AREACODE   string   `xml:"AREACODE,attr"`
	CITYCODE   string   `xml:"CITYCODE,attr"`
	CTARCODE   string   `xml:"CTARCODE,attr"`
	PLACECODE  string   `xml:"PLACECODE,attr"`
	STREETCODE string   `xml:"STREETCODE,attr"`
	EXTRCODE   string   `xml:"EXTRCODE,attr"`
	SEXTCODE   string   `xml:"SEXTCODE,attr"`
	OFFNAME    string   `xml:"OFFNAME,attr"`
	POSTALCODE string   `xml:"POSTALCODE,attr"`
	IFNSFL     string   `xml:"IFNSFL,attr"`
	TERRIFNSFL string   `xml:"TERRIFNSFL,attr"`
	IFNSUL     string   `xml:"IFNSUL,attr"`
	TERRIFNSUL string   `xml:"TERRIFNSUL,attr"`
	OKATO      string   `xml:"OKATO,attr"`
	OKTMO      string   `xml:"OKTMO,attr"`
	UPDATEDATE string   `xml:"UPDATEDATE,attr"`
	SHORTNAME  string   `xml:"SHORTNAME,attr"`
	AOLEVEL    int      `xml:"AOLEVEL,attr"`
	PARENTGUID string   `xml:"PARENTGUID,attr"`
	AOID       string   `xml:"AOID,attr"`
	PREVID     string   `xml:"PREVID,attr"`
	NEXTID     string   `xml:"NEXTID,attr"`
	CODE       string   `xml:"CODE,attr"`
	PLAINCODE  string   `xml:"PLAINCODE,attr"`
	ACTSTATUS  int      `xml:"ACTSTATUS,attr"`
	CENTSTATUS int      `xml:"CENTSTATUS,attr"`
	OPERSTATUS int      `xml:"OPERSTATUS,attr"`
	CURRSTATUS int      `xml:"CURRSTATUS,attr"`
	STARTDATE  string   `xml:"STARTDATE,attr"`
	ENDDATE    string   `xml:"ENDDATE,attr"`
	NORMDOC    string   `xml:"NORMDOC,attr"`
	LIVESTATUS bool     `xml:"LIVESTATUS,attr"`
}

const schema = `CREATE TABLE ` + tableName + ` (
    ao_guid UUID NOT NULL,
    formal_name VARCHAR(120) NOT NULL,
		region_code VARCHAR(2) NOT NULL,
		auto_code VARCHAR(1) NOT NULL,
		area_code VARCHAR(3) NOT NULL,
		city_code VARCHAR(3) NOT NULL,
		ctar_code VARCHAR(3) NOT NULL,
		place_code VARCHAR(3) NOT NULL,
		street_code VARCHAR(4),
		extr_code VARCHAR(4) NOT NULL,
		sext_code VARCHAR(3) NOT NULL,
		off_name VARCHAR(120),
		postal_code VARCHAR(6),
		ifns_fl VARCHAR(4),
		terr_ifns_fl VARCHAR(4),
		ifns_ul VARCHAR(4),
		terr_ifns_ul VARCHAR(4),
		okato VARCHAR(11),
		oktmo VARCHAR(8),
		update_date TIMESTAMP NOT NULL,
		short_name VARCHAR(10) NOT NULL,
		ao_level INT NOT NULL,
		parent_guid UUID,
		ao_id UUID NOT NULL,
		prev_id UUID,
		next_id UUID,
		code VARCHAR(17),
		plain_code VARCHAR(15),
		act_status INT NOT NULL,
		cent_status INT NOT NULL,
		oper_status INT NOT NULL,
		curr_status INT NOT NULL,
		start_date TIMESTAMP NOT NULL,
		end_date TIMESTAMP NOT NULL,
		norm_doc UUID,
		live_status INT NOT NULL,
		PRIMARY KEY (ao_id));`

func Export(db *sqlx.DB, format *string) {
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
	// countedElements, err := helpers.CountElementsInXML(pathToFile, elementName)
	// if err != nil {
	// 	fmt.Println("Error counting elements in XML file:", err)
	// 	return
	// }
	// fmt.Println("Необходимо обработать элементов: ", countedElements)

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
				decoder.DecodeElement(&item, &se)

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

				query := `INSERT INTO ` + tableName + ` (ao_guid,
					formal_name,
					region_code,
					auto_code,
					area_code,
					city_code,
					ctar_code,
					place_code,
					street_code,
					extr_code,
					sext_code,
					off_name,
					postal_code,
					ifns_fl,
					terr_ifns_fl,
					ifns_ul,
					terr_ifns_ul,
					okato,
					oktmo,
					update_date,
					short_name,
					ao_level,
					parent_guid,
					ao_id,
					prev_id,
					next_id,
					code,
					plain_code,
					act_status,
					cent_status,
					oper_status,
					curr_status,
					start_date,
					end_date,
					norm_doc,
					live_status
					) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
						$11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
						$21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
						$31, $32, $33, $34, $35, $36)`
				db.MustExec(query,
					item.AOGUID,
					item.FORMALNAME,
					item.REGIONCODE,
					item.AUTOCODE,
					item.AREACODE,
					item.CITYCODE,
					item.CTARCODE,
					item.PLACECODE,
					item.STREETCODE,
					item.EXTRCODE,
					item.SEXTCODE,
					item.OFFNAME,
					item.POSTALCODE,
					item.IFNSFL,
					item.TERRIFNSFL,
					item.IFNSUL,
					item.TERRIFNSUL,
					item.OKATO,
					item.OKTMO,
					updDate,
					item.SHORTNAME,
					item.AOLEVEL,
					item.PARENTGUID,
					item.AOID,
					item.PREVID,
					item.NEXTID,
					item.CODE,
					item.PLAINCODE,
					item.ACTSTATUS,
					item.CENTSTATUS,
					item.OPERSTATUS,
					item.CURRSTATUS,
					startDate,
					endDate,
					item.NORMDOC,
					item.LIVESTATUS)
				if err != nil {
					fmt.Println("Error on adding row:", err)
					return
				}

				s := strconv.Itoa(total)
				fmt.Printf("\r"+elementName+": %s rows", s)
			}
		default:
		}

	}

	fmt.Printf("Total processed items in "+elementName+": %d \n", total)
}
