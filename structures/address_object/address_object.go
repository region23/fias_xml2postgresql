package address_object

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const dateformat = "2006-01-02"
const tableName = "addrobj"

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
    ao_guid UUID UNIQUE NOT NULL,
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

func Export(db *sqlx.DB) {
	// Создаем таблицу
	dbmap.AddTableWithName(DBObject{}, "addrobj")
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

	xmlFile, err := os.Open("xml/AS_ADDROBJ_20150705_e3a7c988-3be1-456a-a329-ba78c181bb1a.XML")
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

			if inElement == "Object" {
				total++
				var item XmlObject

				// decode a whole chunk of following XML into the
				// variable item which is a ActualStatus (se above)
				decoder.DecodeElement(&item, &se)
				obj, err := xml2db(item)
				if err != nil {
					fmt.Println("Error on mapping XML to DB object: ", err)
					return
				}

				var err error

				obj.UPDATEDATE, err = time.Parse(dateformat, xml.UPDATEDATE)
				if err != nil {
					fmt.Println("Error parse UPDATEDATE: ", err)
					return nil, err
				}

				obj.STARTDATE, err = time.Parse(dateformat, xml.STARTDATE)
				if err != nil {
					fmt.Println("Error parse STARTDATE: ", err)
					return nil, err
				}

				obj.ENDDATE, err = time.Parse(dateformat, xml.ENDDATE)
				if err != nil {
					fmt.Println("Error parse ENDDATE: ", err)
					return nil, err
				}

				err = dbmap.Insert(obj)
				if err != nil {
					fmt.Println("Error on creating table:", err)
					return
				}

				s := strconv.Itoa(total)
				fmt.Printf("\rObject: %s rows", s)
			}
		default:
		}

	}

	fmt.Printf("Total processed items in AddressObjects: %d \n", total)
}
