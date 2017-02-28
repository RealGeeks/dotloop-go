// Package dotloop is a client for the dotloop API
//
// https://dotloop.github.io/public-api/
package dotloop

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// DefaultURL for the API v2
const DefaultURL = "https://api-gateway.dotloop.com/public/v2/"

// ErrInvalid is returned on status code 400. Usually a validation error
// like missing a required field
type ErrInvalid struct {
	Body string
}

func (e *ErrInvalid) Error() string {
	return fmt.Sprintf("dotloop: invalid request %v", e.Body)
}

// Dotloop is the client to API v2
//
// https://dotloop.github.io/public-api/
type Dotloop struct {
	Token string       // OAuth access token
	URL   string       // (optional) defaults to DefaultURL
	HTTP  *http.Client // (optional) http client to perform requests
}

// LoopIt creates a new loop
//
// https://dotloop.github.io/public-api/#loop-it
func (dl *Dotloop) LoopIt(loop Loop) error {
	reqbody, err := json.Marshal(loop)
	if err != nil {
		return fmt.Errorf("dotloop: encoding body (%v)", err)
	}
	req, err := http.NewRequest("POST", dl.url("loop-it"), bytes.NewReader(reqbody))
	if err != nil {
		return fmt.Errorf("dotloop: building request (%v)", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+dl.Token)
	res, err := dl.http().Do(req)
	if err != nil {
		return fmt.Errorf("dotloop: making request (%v)", err)
	}
	defer res.Body.Close()
	resbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("dotloop: reading response (%v)", err)
	}
	if res.StatusCode == 400 {
		return &ErrInvalid{Body: string(resbody)}
	}
	if res.StatusCode != 201 {
		return fmt.Errorf("dotloop: %v - %v", res.StatusCode, string(resbody))
	}
	return nil
}

func (dl *Dotloop) http() *http.Client {
	if dl.HTTP == nil {
		dl.HTTP = &http.Client{Timeout: 3 * time.Second}
	}
	return dl.HTTP
}

func (dl *Dotloop) url(path string) string {
	u := dl.URL
	if u == "" {
		u = DefaultURL
	}
	return u + path
}

type Loop struct {
	Name            string        `json:"name"`
	ProfileID       int           `json:"profile_id,omitempty"`
	TemplateID      int           `json:"templateId,omitempty"`
	TransactionType string        `json:"transactionType"`
	Status          string        `json:"status"`
	Participants    []Participant `json:"participants,omitempty"`
	StreetName      string        `json:"streetName,omitempty"`
	StreetNumber    string        `json:"streetNumber,omitempty"`
	Unit            string        `json:"unit,omitempty"`
	City            string        `json:"city,omitempty"`
	State           string        `json:"state,omitempty"`
	ZipCode         string        `json:"zipCode,omitempty"`
	County          string        `json:"county,omitempty"`
	Country         string        `json:"country,omitempty"`
	MLSPropertyID   string        `json:"mlsPropertId,omitempty"`
	MLSAgentID      string        `json:"mlsAgentId,omitempty"`
}

type Participant struct {
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
