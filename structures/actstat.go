package structures

// Статус актуальности ФИАС
type ActualStatus struct {
	XMLName   xml.Name `xml:"ActualStatus"`
	ActStatId int      `xml:"ACTSTATID,attr"`
	Name      string   `xml:"NAME,attr"`
}

type ActualStatuses struct {
	XMLName        xml.Name       `xml:"ActualStatuses"`
	ActualStatuses []ActualStatus `xml:"ActualStatus"`
}
