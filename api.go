package okta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Client to access okta
type Client struct {
	client   *http.Client
	org      string
	Url      string
	ApiToken string
}

// errorResponse is an error wrapper for the okta response
type errorResponse struct {
	HTTPCode int
	Response ErrorResponse
	Endpoint string
}

type OTPResponse struct {
	ExpiresAt  time.Time `json:"expiresAt"`
	Status     string    `json:"status"`
	RelayState string    `json:"relayState"`
	Embedded   struct {
		User struct {
			ID              string    `json:"id"`
			PasswordChanged time.Time `json:"passwordChanged"`
			Profile         struct {
				Login     string `json:"login"`
				FirstName string `json:"firstName"`
				LastName  string `json:"lastName"`
				Locale    string `json:"locale"`
				TimeZone  string `json:"timeZone"`
			} `json:"profile"`
		} `json:"user"`
	} `json:"_embedded"`
}

func (e *errorResponse) Error() string {
	return fmt.Sprintf("Error hitting api endpoint %s %s", e.Endpoint, e.Response.ErrorCode)
}

// NewClient object for calling okta
func NewClient(org string) *Client {
	client := Client{
		client: &http.Client{},
		org:    org,
		Url:    "okta.com",
	}

	return &client
}

// Authenticate with okta using username and password
func (c *Client) Authenticate(username, password string) (*AuthnResponse, error) {
	var request = &AuthnRequest{
		Username: username,
		Password: password,
	}

	var response = &AuthnResponse{}
	err := c.call("authn", "POST", request, response)
	return response, err
}

// VerifyOTP will validate the supplied TOTP code
func (r *AuthnResponse) VerifyOTP(otp string, verifyResponse *OTPResponse) error {
	request := map[string]string{"passCode": otp, "stateToken": r.StateToken}
	data, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", r.Embedded.Factors[0].Links.Verify.Href, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Accept", `application/json`)
	req.Header.Add("Content-Type", `application/json`)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == http.StatusOK {
		err := json.Unmarshal(body, verifyResponse)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		var errors ErrorResponse
		_ = json.Unmarshal(body, &errors)
		return &errorResponse{
			HTTPCode: resp.StatusCode,
			Response: errors,
		}
	}

	return nil
}

// Session takes a session token and always fails
func (c *Client) Session(sessionToken string) (*SessionResponse, error) {
	var request = &SessionRequest{
		SessionToken: sessionToken,
	}

	var response = &SessionResponse{}
	err := c.call("sessions", "POST", request, response)
	return response, err
}

// User takes a user id and returns data about that user
func (c *Client) User(userID string) (*User, error) {

	var response = &User{}
	err := c.call("users/"+userID, "GET", nil, response)
	return response, err
}

// Groups takes a user id and returns the groups the user belongs to
func (c *Client) Groups(userID string) (*Groups, error) {

	var response = &Groups{}
	err := c.call("users/"+userID+"/groups", "GET", nil, response)
	return response, err
}

func (c *Client) call(endpoint, method string, request, response interface{}) error {
	data, _ := json.Marshal(request)

	var url = "https://" + c.org + "." + c.Url + "/api/v1/" + endpoint
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Accept", `application/json`)
	req.Header.Add("Content-Type", `application/json`)
	if c.ApiToken != "" {
		req.Header.Add("Authorization", "SSWS "+c.ApiToken)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == http.StatusOK {
		err := json.Unmarshal(body, &response)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		var errors ErrorResponse
		err = json.Unmarshal(body, &errors)

		return &errorResponse{
			HTTPCode: resp.StatusCode,
			Response: errors,
			Endpoint: url,
		}
	}

	return nil
}
