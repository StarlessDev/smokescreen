package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// This fields are present in every
// response of the API
type APIResponse struct {
	Success  bool      `json:"success"`
	Messages []Message `json:"messages"`
	Errors   []Message `json:"errors"`
}

// This struct is used by both messages and errors
type Message struct {
	Code             int    `json:"code"`
	Message          string `json:"message"`
	DocumentationUrl string `json:"documentation_url"`
	Pointer          string `json:"source.pointer"`
}

// Represents what Cloudflare will do
// when they receive an email to the
// generated address
type EmailAction struct {
	Type  string   `json:"type"`
	Value []string `json:"value"`
}

// Represents how Cloudflare will match
// the address of the incoming emails.
type EmailMatcher struct {
	Type  string `json:"type"`
	Field string `json:"field"`
	Value string `json:"value"`
}

// Represents an Email
// The Id and Tag fields are populated only
// when this struct is in an API response.
type Email struct {
	Id       string         `json:"id,omitempty"`
	Tag      string         `json:"tag,omitempty"`
	Name     string         `json:"name"`
	Enabled  bool           `json:"enabled"`
	Priority int            `json:"priority"`
	Actions  []EmailAction  `json:"actions"`
	Matchers []EmailMatcher `json:"matchers"`
}

// Requests payloads
type GenEmailBody struct {
	Email
	Actions  []EmailAction  `json:"actions"`
	Matchers []EmailMatcher `json:"matchers"`
}

// Resposes bodies
type ListEmailResponse struct {
	APIResponse
	Result []Email `json:"result"`
}

type GenEmailResponse struct {
	APIResponse
	Result Email `json:"result"`
}

type RevokeEmailResponse struct {
	APIResponse
	Result string `json:"result"`
}

var HttpClient = &http.Client{
	Timeout: 5 * time.Second,
}

var EmailNamePrefix = "smokescreen-"

func GenerateEmail(id *Identity, name string, email string) (*GenEmailResponse, error) {
	payload := GenEmailBody{
		Email: Email{
			Enabled:  true,
			Name:     name,
			Priority: 0,
		},
		Actions: []EmailAction{
			{
				Type: "forward",
				Value: []string{
					id.Email,
				},
			},
		},
		Matchers: []EmailMatcher{
			{
				Field: "to",
				Type:  "literal",
				Value: email,
			},
		},
	}

	responseBody := new(GenEmailResponse)
	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", getApiURL(id), bytes.NewBuffer(jsonPayload))
	populateHeaders(req, id)

	res, err := HttpClient.Do(req)
	if err != nil {
		return responseBody, err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	if err := json.Unmarshal(body, responseBody); err != nil {
		return responseBody, err
	}

	return responseBody, nil
}

func RevokeEmail(id *Identity, tag string) (*RevokeEmailResponse, error) {
	responseBody := new(RevokeEmailResponse)
	url := getApiURL(id) + "/" + tag
	req, _ := http.NewRequest("DELETE", url, nil)
	populateHeaders(req, id)

	res, err := HttpClient.Do(req)
	if err != nil {
		return responseBody, err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	if err := json.Unmarshal(body, responseBody); err != nil {
		return responseBody, err
	}

	return responseBody, nil
}

func ListEmails(id *Identity) (*ListEmailResponse, error) {
	responseBody := new(ListEmailResponse)
	req, _ := http.NewRequest("GET", getApiURL(id), nil)
	populateHeaders(req, id)

	res, err := HttpClient.Do(req)
	if err != nil {
		return responseBody, err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	if err := json.Unmarshal(body, responseBody); err != nil {
		return responseBody, err
	}

	return responseBody, nil
}

func populateHeaders(request *http.Request, id *Identity) {
	request.Header.Set("User-Agent", "smokescreen-cli")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", id.Token))
}

func getApiURL(id *Identity) string {
	return fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/email/routing/rules", id.ZoneId)
}
