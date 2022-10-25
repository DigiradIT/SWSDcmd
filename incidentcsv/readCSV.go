package incidentcsv

import (
	"os"

	"github.com/gocarina/gocsv"
)

type Incident struct {
	Description string `csv:"description" json:"description"`
	Name        string `csv:"name" json:"name"`
	Requester   string `csv:"requester" json:"requester"`
	Assignee    string `csv:"assignee" json:"assignee"`
	Category    string `csv:"category" json:"category"`
	Subcategory string `csv:"subcategory" json:"subcategory"`
	Computer    string `csv:"computer" json:"computer"`
	Incidents   string `csv:"incidents" json:"incidents"`
}

func ReadCSV(path string) []Incident {
	csvf, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	incidents := []Incident{}

	if err := gocsv.UnmarshalBytes(csvf, &incidents); err != nil {
		panic(err)
	}

	return incidents
}
