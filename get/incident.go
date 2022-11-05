package get

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Incident struct {
	Id     int    `json:"id"`
	Number int    `json:"number"`
	Href   string `json:"href"`
}

func IncidentById(id string, key string) (*Incident, error) {
	query_url := baseURL + "/incidents/"
	client := resty.New()

	resp, err := client.R().
		SetHeader("X-Samanage-Authorization", "Bearer "+key).
		SetHeader("Accept", "application/json").
		Get(query_url + id + ".json")

	if err != nil {
		fmt.Println("error fetching computer")
		return nil, err
	}
	var incident *Incident
	err = json.Unmarshal(resp.Body(), &incident)

	if err != nil {
		fmt.Println("error unmarshalling json")
		return nil, err
	}

	return incident, nil
}
