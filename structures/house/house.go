package house

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-gorp/gorp"
	_ "github.com/lib/pq"
)

const dateformat = "2006-01-02"

// Сведения по номерам домов улиц городов и населенных пунктов, номера земельных участков и т.п
type XmlObject struct {
	XMLName    xml.Name `xml:"House"`
	POSTALCODE string   `xml:"POSTALCODE,attr"`
	IFNSFL     string   `xml:"IFNSFL,attr"`
	TERRIFNSFL string   `xml:"TERRIFNSFL,attr"`
	IFNSUL     string   `xml:"IFNSUL,attr"`
	TERRIFNSUL string   `xml:"TERRIFNSUL,attr"`
	OKATO      string   `xml:"OKATO,attr"`
	OKTMO      string   `xml:"OKTMO,attr"`
	UPDATEDATE string   `xml:"UPDATEDATE,attr"`
	HOUSENUM   string   `xml:"HOUSENUM,attr"`
	ESTSTATUS  int      `xml:"ESTSTATUS,attr"`
	BUILDNUM   string   `xml:"BUILDNUM,attr"`
	STRUCNUM   string   `xml:"STRUCNUM,attr"`
	STRSTATUS  int      `xml:"STRSTATUS,attr"`
	HOUSEID    string   `xml:"HOUSEID,attr"`
	HOUSEGUID  string   `xml:"HOUSEGUID,attr"`
	AOGUID     string   `xml:"AOGUID,attr"`
	STARTDATE  string   `xml:"STARTDATE,attr"`
	ENDDATE    string   `xml:"ENDDATE,attr"`
	STATSTATUS int      `xml:"STATSTATUS,attr"`
	NORMDOC    string   `xml:"NORMDOC,attr"`
	COUNTER    string   `xml:"COUNTER,attr"`
}

type DBObject struct {
	POSTALCODE string    `db:"postal_code"`
	IFNSFL     string    `db:"ifnsfl"`
	REGIONCODE string    `db:"region_code"`
	AUTOCODE   string    `db:"auto_code"`
	AREACODE   string    `db:"area_code"`
	CITYCODE   string    `db:"city_code"`
	CTARCODE   string    `db:"ctar_code"`
	PLACECODE  string    `db:"place_code"`
	STREETCODE string    `db:"street_code"`
	EXTRCODE   string    `db:"extr_code"`
	SEXTCODE   string    `db:"sext_code"`
	OFFNAME    string    `db:"off_name"`
	POSTALCODE string    `db:"postal_code"`
	IFNSFL     string    `db:"ifnsfl"`
	TERRIFNSFL string    `db:"terrifnsfl"`
	IFNSUL     string    `db:"ifnsul"`
	TERRIFNSUL string    `db:"terrifnsul"`
	OKATO      string    `db:"okato"`
	OKTMO      string    `db:"oktmo"`
	UPDATEDATE time.Time `db:"update_date"`
	SHORTNAME  string    `db:"short_name"`
	AOLEVEL    int       `db:"ao_level"`
	PARENTGUID string    `db:"parent_guid"`
	AOID       string    `db:"ao_id, primarykey"`
	PREVID     string    `db:"prev_id"`
	NEXTID     string    `db:"next_id"`
	CODE       string    `db:"code"`
	PLAINCODE  string    `db:"plain_code"`
	ACTSTATUS  int       `db:"act_status"`
	CENTSTATUS int       `db:"cent_status"`
	OPERSTATUS int       `db:"oper_status"`
	CURRSTATUS int       `db:"curr_status"`
	STARTDATE  time.Time `db:"start_date"`
	ENDDATE    time.Time `db:"end_date"`
	NORMDOC    string    `db:"norm_doc"`
	LIVESTATUS bool      `db:"live_status"`
}

func xml2db(xml XmlObject) (*DBObject, error) {
	obj := &DBObject{
		AOGUID:     xml.AOGUID,
		FORMALNAME: xml.FORMALNAME,
		REGIONCODE: xml.REGIONCODE,
		AUTOCODE:   xml.AUTOCODE,
		AREACODE:   xml.AREACODE,
		CITYCODE:   xml.CITYCODE,
		CTARCODE:   xml.CTARCODE,
		PLACECODE:  xml.PLACECODE,
		STREETCODE: xml.STREETCODE,
		EXTRCODE:   xml.EXTRCODE,
		SEXTCODE:   xml.SEXTCODE,
		OFFNAME:    xml.OFFNAME,
		POSTALCODE: xml.POSTALCODE,
		IFNSFL:     xml.IFNSFL,
		TERRIFNSFL: xml.TERRIFNSFL,
		IFNSUL:     xml.IFNSUL,
		TERRIFNSUL: xml.TERRIFNSUL,
		OKATO:      xml.OKATO,
		OKTMO:      xml.OKTMO,
		SHORTNAME:  xml.SHORTNAME,
		AOLEVEL:    xml.AOLEVEL,
		PARENTGUID: xml.PARENTGUID,
		AOID:       xml.AOID,
		PREVID:     xml.PREVID,
		NEXTID:     xml.NEXTID,
		CODE:       xml.CODE,
		PLAINCODE:  xml.PLAINCODE,
		ACTSTATUS:  xml.ACTSTATUS,
		CENTSTATUS: xml.CENTSTATUS,
		OPERSTATUS: xml.OPERSTATUS,
		CURRSTATUS: xml.CURRSTATUS,
		NORMDOC:    xml.NORMDOC,
		LIVESTATUS: xml.LIVESTATUS}

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

	return obj, nil
}

func Export(dbmap *gorp.DbMap) {
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
