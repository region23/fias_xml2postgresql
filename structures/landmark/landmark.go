package landmark

import "encoding/xml"

//
// const dateformat = "2006-01-02"
// const tableName = "as_landmark"
// const elementName = "Landmark"

// Описание мест расположения  имущественных объектов
type XmlObject struct {
	XMLName    xml.Name `xml:"Landmark" db:"as_landmark"`
	LOCATION   string   `xml:"LOCATION,attr" db:"location"`
	POSTALCODE *string  `xml:"POSTALCODE,attr,omitempty" db:"postal_code"`
	IFNSFL     int      `xml:"IFNSFL,attr,omitempty" db:"ifns_fl"`
	TERRIFNSFL int      `xml:"TERRIFNSFL,attr,omitempty" db:"terr_ifns_fl"`
	IFNSUL     int      `xml:"IFNSUL,attr,omitempty" db:"ifns_ul"`
	TERRIFNSUL int      `xml:"TERRIFNSUL,attr,omitempty" db:"terr_ifns_ul"`
	OKATO      *string  `xml:"OKATO,attr,omitempty" db:"okato"`
	OKTMO      *string  `xml:"OKTMO,attr,omitempty" db:"oktmo"`
	UPDATEDATE string   `xml:"UPDATEDATE,attr" db:"update_date"`
	LANDID     string   `xml:"LANDID,attr" db:"land_id"`
	LANDGUID   string   `xml:"LANDGUID,attr" db:"land_guid"`
	AOGUID     string   `xml:"AOGUID,attr" db:"ao_guid"`
	STARTDATE  string   `xml:"STARTDATE,attr" db:"start_date"`
	ENDDATE    string   `xml:"ENDDATE,attr" db:"end_date"`
	NORMDOC    *string  `xml:"NORMDOC,attr,omitempty" db:"norm_doc"`
}

func Schema(tableName string) string {
	return `CREATE TABLE ` + tableName + ` (
    location VARCHAR(500) NOT NULL,
    postal_code VARCHAR(6),
    ifns_fl INT,
    terr_ifns_fl INT,
    ifns_ul INT,
    terr_ifns_ul INT,
    okato VARCHAR(11),
    oktmo VARCHAR(11),
    update_date TIMESTAMP NOT NULL,
    land_id UUID NOT NULL,
    land_guid UUID NOT NULL,
    ao_guid UUID NOT NULL,
    start_date TIMESTAMP NOT NULL,
		end_date TIMESTAMP NOT NULL,
		norm_doc UUID,
		PRIMARY KEY (land_id));`
}

/*
func Export(w *sync.WaitGroup, c chan string, db *sqlx.DB, format *string) {

	defer w.Done()

	helpers.DropAndCreateTable(schema, tableName, db)

	var format2 string
	format2 = *format
	fileName, err2 := helpers.SearchFile(tableName, format2)
	if err2 != nil {
		log.Println("Error searching file:", err2)
		return
	}

	pathToFile := format2 + "/" + fileName

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

				query := `INSERT INTO ` + tableName + ` (location,
					postal_code,
					ifns_fl,
					terr_ifns_fl,
					ifns_ul,
					terr_ifns_ul,
					okato,
					oktmo,
					update_date,
          land_id,
          land_guid,
          ao_guid,
          start_date,
					end_date,
					norm_doc
					) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
						$11, $12, $13, $14, $15)`

				_, err = db.Exec(query,
					item.LOCATION,
					item.POSTALCODE,
					item.IFNSFL,
					item.TERRIFNSFL,
					item.IFNSUL,
					item.TERRIFNSUL,
					item.OKATO,
					item.OKTMO,
					item.UPDATEDATE,
					item.LANDID,
					item.LANDGUID,
					item.AOGUID,
					item.STARTDATE,
					item.ENDDATE,
					item.NORMDOC)

				if err != nil {
					log.Fatal(err)
				}

				c <- helpers.PrintRowsAffected(elementName, total)
			}
		default:
		}
	}
}
*/
