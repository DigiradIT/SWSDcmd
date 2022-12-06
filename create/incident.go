package create

import (
	"errors"
	"fmt"

	"github.com/DigiradIT/SWSDcmd/get"
	"github.com/DigiradIT/SWSDcmd/incidentcsv"
	"github.com/go-resty/resty/v2"
)

var incident_url = "https://api.samanage.com/incidents.json"

type User struct {
	Email string `json:"email"`
}

type Category struct {
	Name string `json:"name"`
}

type Computer struct {
	Id   int    `json:"id"`
	Href string `json:"href"`
}

type Incident struct {
	Id     int    `json:"id"`
	Href   string `json:"href"`
	Number int    `json:"number"`
}

type SubProcRes struct {
	sub Submission
	e   error
}

type Submission struct {
	Description string     `json:"description"`
	Name        string     `json:"name"`
	Requester   User       `json:"requester"`
	Assignee    User       `json:"assignee"`
	Category    Category   `json:"category"`
	Subcategory Category   `json:"subcategory"`
	Assets      []Computer `json:"assets"`
	Incidents   []Incident `json:"incidents"`
}

type Wrapper struct {
	Incident Submission `json:"incident"`
}

func TranslateFromCSV(in incidentcsv.Incident, key string) (Submission, error) {
	var sub Submission
	sub.Name = in.Name
	sub.Description = in.Description
	sub.Requester = User{in.Requester}
	sub.Assignee = User{in.Assignee}
	sub.Category = Category{in.Category}
	sub.Subcategory = Category{in.Subcategory}

	var computer *get.Comp
	var err error
	if in.Computer != "" {
		computer, err = get.ComputerByName(in.Computer, key)
		if err != nil {
			return Submission{}, errors.New("Error fetching computer info.")
		}
	} else {
		computer = nil
	}

	var incident *get.Incident
	err = nil
	if in.Incidents != "" {
		incident, err = get.IncidentById(in.Incidents, key)

		if err != nil {
			return Submission{}, errors.New("Error fetching incident info.")
		}
	} else {
		incident = nil
	}
	if computer != nil {
		sub.Assets = []Computer{{
			Id:   computer.Id,
			Href: computer.Href,
		},
		}
	}

	if incident != nil {
		sub.Incidents = []Incident{{
			Id:     incident.Id,
			Href:   incident.Href,
			Number: incident.Number,
		},
		}
	}

	return sub, nil
}

func Submit(sub Submission, key string) error {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "Application/json").
		SetHeader("X-Samanage-Authorization", "Bearer "+key).
		SetBody(Wrapper{sub}).
		Post("https://api.samanage.com/incidents.json")

	if err != nil {
		return errors.New(resp.String())
	} else {
		return nil
	}
}

func SubmitCSV(incidents []incidentcsv.Incident, key string) []error {
	errs := []error{}
	subs := []Submission{}
	resChan := make(chan SubProcRes)

	for _, s := range incidents {
		go func(i incidentcsv.Incident) {
			r, err := TranslateFromCSV(i, key)
			if err != nil {
				es := fmt.Sprintf("Submisison build failed: %s reason: %s", i.Name, err)
				resChan <- SubProcRes{
					Submission{},
					errors.New(es),
				}
			} else {
				resChan <- SubProcRes{
					r,
					nil,
				}
			}
		}(s)
	}

	for i := 0; i < len(incidents); i++ {
		res := <-resChan
		if res.e != nil {
			errs = append(errs, res.e)
		} else {
			subs = append(subs, res.sub)
		}
	}

	subErrChan := make(chan error)

	for _, s := range subs {
		go func(su Submission) {
			e := Submit(su, key)
			if e != nil {
				es := fmt.Sprintf("Submission to SWSD failed: %s reason: %s", su.Name, e)
				subErrChan <- errors.New(es)
			} else {
				subErrChan <- nil
			}
		}(s)
	}

	for i := 0; i < len(subs); i++ {
		e := <-subErrChan
		if e != nil {
			errs = append(errs, e)
		}
	}

	return errs
}
