package get

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

var baseURL = "https://api.samanage.com"

type Comp struct {
	Id   int    `json:"id"`
	Href string `json:"href"`
}

func ComputerByName(name string, key string) (*Comp, error) {
	query_url := baseURL + "/hardwares.json"
	client := resty.New()

	resp, err := client.R().
		SetHeader("X-Samanage-Authorization", "Bearer "+key).
		SetHeader("Accept", "application/json").
		SetQueryParam("name", name).
		Get(query_url)

	if err != nil {
		fmt.Println("error fetching computer")
		return nil, err
	}
	computers := []*Comp{}
	err = json.Unmarshal(resp.Body(), &computers)

	if err != nil {
		fmt.Println("error unmarshalling json")
		return nil, err
	}

	return computers[0], nil
}
