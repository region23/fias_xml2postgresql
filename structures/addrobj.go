package structures

import "encoding/xml"

// Классификатор адресообразующих элементов
type Object struct {
	XMLName    xml.Name `xml:"Object"`
	AOGUID     string   `xml:"AOGUID, attr"`
	FORMALNAME string   `xml:"FORMALNAME, attr"`
}

type AddressObjects struct {
	XMLName        xml.Name `xml:"AddressObjects"`
	AddressObjects []Object `xml:"Object"`
}
