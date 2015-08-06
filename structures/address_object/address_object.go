package address_object

import (
	"encoding/xml"

	_ "github.com/lib/pq"
)

// const dateformat = "2006-01-02"
// const tableName = "as_addrobj"
// const elementName = "Object"

// Классификатор адресообразующих элементов
type XmlObject struct {
	XMLName    xml.Name `xml:"Object" db:"as_addrobj"`
	AOGUID     string   `xml:"AOGUID,attr" db:"ao_guid"`
	FORMALNAME string   `xml:"FORMALNAME,attr" db:"formal_name"`
	REGIONCODE int      `xml:"REGIONCODE,attr" db:"region_code"`
	AUTOCODE   int      `xml:"AUTOCODE,attr" db:"auto_code"`
	AREACODE   int      `xml:"AREACODE,attr" db:"area_code"`
	CITYCODE   int      `xml:"CITYCODE,attr" db:"city_code"`
	CTARCODE   int      `xml:"CTARCODE,attr" db:"ctar_code"`
	PLACECODE  int      `xml:"PLACECODE,attr" db:"place_code"`
	STREETCODE int      `xml:"STREETCODE,attr,omitempty" db:"street_code"`
	EXTRCODE   int      `xml:"EXTRCODE,attr" db:"extr_code"`
	SEXTCODE   int      `xml:"SEXTCODE,attr" db:"sext_code"`
	OFFNAME    *string  `xml:"OFFNAME,attr,omitempty" db:"off_name"`
	POSTALCODE *string  `xml:"POSTALCODE,attr,omitempty" db:"postal_code"`
	IFNSFL     int      `xml:"IFNSFL,attr,omitempty" db:"ifns_fl"`
	TERRIFNSFL int      `xml:"TERRIFNSFL,attr,omitempty" db:"terr_ifns_fl"`
	IFNSUL     int      `xml:"IFNSUL,attr,omitempty" db:"ifns_ul"`
	TERRIFNSUL int      `xml:"TERRIFNSUL,attr,omitempty" db:"terr_ifns_ul"`
	OKATO      *string  `xml:"OKATO,attr,omitempty" db:"okato"`
	OKTMO      *string  `xml:"OKTMO,attr,omitempty" db:"oktmo"`
	UPDATEDATE string   `xml:"UPDATEDATE,attr" db:"update_date"`
	SHORTNAME  string   `xml:"SHORTNAME,attr" db:"short_name"`
	AOLEVEL    int      `xml:"AOLEVEL,attr" db:"ao_level"`
	PARENTGUID *string  `xml:"PARENTGUID,attr,omitempty" db:"parent_guid"`
	AOID       string   `xml:"AOID,attr" db:"ao_id"`
	PREVID     *string  `xml:"PREVID,attr,omitempty" db:"prev_id"`
	NEXTID     *string  `xml:"NEXTID,attr,omitempty" db:"next_id"`
	CODE       *string  `xml:"CODE,attr,omitempty" db:"code"`
	PLAINCODE  *string  `xml:"PLAINCODE,attr,omitempty" db:"plain_code"`
	ACTSTATUS  bool     `xml:"ACTSTATUS,attr" db:"act_status"`
	CENTSTATUS bool     `xml:"CENTSTATUS,attr" db:"cent_status"`
	OPERSTATUS int      `xml:"OPERSTATUS,attr" db:"oper_status"`
	CURRSTATUS int      `xml:"CURRSTATUS,attr" db:"curr_status"`
	STARTDATE  string   `xml:"STARTDATE,attr" db:"start_date"`
	ENDDATE    string   `xml:"ENDDATE,attr" db:"end_date"`
	NORMDOC    *string  `xml:"NORMDOC,attr,omitempty" db:"norm_doc"`
	LIVESTATUS bool     `xml:"LIVESTATUS,attr" db:"live_status"`
}

func Schema(tableName string) string {
	return `CREATE TABLE ` + tableName + ` (
    ao_guid UUID NOT NULL,
    formal_name VARCHAR(120) NOT NULL,
		region_code INT NOT NULL,
		auto_code INT NOT NULL,
		area_code INT NOT NULL,
		city_code INT NOT NULL,
		ctar_code INT NOT NULL,
		place_code INT NOT NULL,
		street_code INT,
		extr_code INT NOT NULL,
		sext_code INT NOT NULL,
		off_name VARCHAR(120),
		postal_code VARCHAR(6),
		ifns_fl INT,
		terr_ifns_fl INT,
		ifns_ul INT,
		terr_ifns_ul INT,
		okato VARCHAR(11),
		oktmo VARCHAR(11),
		update_date TIMESTAMP NOT NULL,
		short_name VARCHAR(10) NOT NULL,
		ao_level INT NOT NULL,
		parent_guid UUID,
		ao_id UUID NOT NULL,
		prev_id UUID,
		next_id UUID,
		code VARCHAR(17),
		plain_code VARCHAR(15),
		act_status BOOL NOT NULL,
		cent_status BOOL NOT NULL,
		oper_status INT NOT NULL,
		curr_status INT NOT NULL,
		start_date TIMESTAMP NOT NULL,
		end_date TIMESTAMP NOT NULL,
		norm_doc UUID,
		live_status BOOL NOT NULL,
		PRIMARY KEY (ao_id));`
}

/*
func Export(w *sync.WaitGroup, c chan string, db *sqlx.DB, format *string) {

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

				//fmt.Println(item, "\n\n")

				var err error

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

				_, err = db.Exec(query,
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
					item.UPDATEDATE,
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
					item.STARTDATE,
					item.ENDDATE,
					item.NORMDOC,
					item.LIVESTATUS)

				if err != nil {
					log.Fatal(err)
				}

				c <- helpers.PrintRowsAffected(elementName, total)
			}
		default:
		}

	}

	//fmt.Printf("\nВсего в "+elementName+" обработано %d строк\n", total)
}
*/
