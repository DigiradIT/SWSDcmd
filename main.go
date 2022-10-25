package main

import (
	"flag"
	"fmt"

	"github.com/DigiradIT/SWSDcmd/create"
	"github.com/DigiradIT/SWSDcmd/incidentcsv"
)

func main() {

	var (
		key           string
		csvPath       string
		help          bool
		dispCSVFormat bool
	)

	flag.StringVar(&key, "key", "", "API for SWSD REST API")
	flag.StringVar(&csvPath, "csv", "", "Path to CSV file containing incident information")
	flag.BoolVar(&help, "help", false, "Displays usage information.")
	flag.BoolVar(&dispCSVFormat, "csvhelp", false, "Displays formatting help for creating incidents from CSV")

	flag.Parse()

	if help {
		flag.Usage()
		return
	}
	if dispCSVFormat {
		fmt.Println("Required headers: description,name,requester,assignee,category,subcategory,computer,incidents")
		fmt.Println("")
		fmt.Println("description: string field that describes incident.")
		fmt.Println("name: string field that will be the title of the incident.")
		fmt.Println("requester: email address (name@digirad.com) that will be set as requester.")
		fmt.Println("assignee: email address (name@digirad.com) that will be set as assignee.")
		fmt.Println("category: string field that will be set to category.")
		fmt.Println("subcategory: string field that will be set to subcategory.")
		fmt.Println("computer: computer that will be attched to the incident.")
		fmt.Println("incident: incident number that will be attached to incident.  9 digit number can be found in the URL of the incident immediately following /incidents/.")
		return
	}
	if key == "" {
		fmt.Println("API key needed")
		return
	}
	if csvPath == "" {
		fmt.Println("Path to CSV file needed")
		return
	}
	is := incidentcsv.ReadCSV(csvPath)

	errs := create.SubmitCSV(is, key)

	if len(errs) == 0 {
		fmt.Printf("%d incidents created successfully.", len(is))
	} else {
		for _, e := range errs {
			fmt.Println(e)
		}
	}

}
