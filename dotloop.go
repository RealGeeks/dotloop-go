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
	ReqBody, Method, URL string
	ResBody              string
}

func (e *ErrInvalid) Error() string {
	return fmt.Sprintf("dotloop: %s %s %s returned 400 with body %s", e.Method, e.URL, e.ReqBody, e.ResBody)
}

// ErrInvalidToken is returned on status 401 when access token is invalid
type ErrInvalidToken struct {
	Msg string
}

func (e *ErrInvalidToken) Error() string {
	return fmt.Sprintf("dotloop: %s", e.Msg)
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
		return &ErrInvalid{Method: "POST", URL: dl.url("loop-it"), ReqBody: string(reqbody), ResBody: string(resbody)}
	}
	if ok, err := isInvalidToken(res.StatusCode, resbody); ok {
		return err
	}
	if res.StatusCode != 201 {
		return fmt.Errorf("dotloop: %v - %v", res.StatusCode, string(resbody))
	}
	return nil
}

func isInvalidToken(code int, body []byte) (ok bool, err error) {
	if code != 401 {
		return false, nil
	}
	var data map[string]string
	if err := json.Unmarshal(body, &data); err != nil {
		return false, nil
	}
	if data["error"] != "invalid_token" {
		return false, nil
	}
	return true, &ErrInvalidToken{Msg: "dotloop: " + data["error_description"]}
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
