// Package eyeson provides a client to interact with the eyeson video API to
// start video meetings, create access for participants, control recordings,
// add media like overlay images, play videos, start and stop broadcasts, send
// chat messages or assign participants to various layouts.
package eyeson

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	endpoint  = "https://api.eyeson.team"
	userAgent = "eyeson-go"
)

// Client provides methods to communicate with the eyeson API, starting video
// meetings, adapt configurations and send chat, images, presentations and
// videos to the meeting.
type Client struct {
	apiKey  string
	BaseURL *url.URL

	client  *http.Client
	Rooms   *RoomsService
	Webhook *WebhookService
}

type service struct {
	client *Client
}

// NewClient creates a new client in order to send requests to the eyeson API.
func NewClient(key string) *Client {
	baseURL, _ := url.Parse(endpoint)
	c := &Client{apiKey: key, BaseURL: baseURL, client: http.DefaultClient}
	c.Rooms = &RoomsService{c}
	c.Webhook = &WebhookService{c}
	return c
}

// UserClient provides a client for user requests that use the session access
// key for authorization.
func (c *Client) UserClient() *Client {
	return &Client{BaseURL: c.BaseURL, client: c.client}
}

// NewRequest prepares a request to be sent to the API.
func (c *Client) NewRequest(method, urlStr string, data url.Values) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", c.apiKey)
	}
	req.Header.Set("User-Agent", userAgent)

	return req, nil
}

// Do sends a request to the eyeson API and prepares the result from the
// received response.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = validateResponse(resp)
	if err != nil {
		// body, _ := ioutil.ReadAll(resp.Body)
		// log.Printf("Received body: %s", body)
		return nil, err
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
		if err != nil {
			return nil, err
		}
	}
	return resp, err
}

func validateResponse(resp *http.Response) error {
	c := resp.StatusCode
	switch {
	case c == 200 || c == 201:
		return nil
	case c == 404:
		return errors.New("Not found! Resource does not exist or expired.")
	case c == 401:
		return errors.New("Authorization failed! Check the API key to be valid.")
	case c == 403:
		return errors.New("Bad request! Check your request parameters to be valid.")
	default:
		return errors.New(fmt.Sprintf("Unknown error! Request failed for an unknown error. (%d)", c))
	}
}
