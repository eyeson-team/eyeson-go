// Package eyeson provides a client to interact with the eyeson video API to
// start video meetings, create access for participants, control recordings,
// add media like overlay images, play videos, start and stop broadcasts, send
// chat messages or assign participants to various layouts.
package eyeson

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
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

	client       *http.Client
	Rooms        *RoomsService
	Webhook      *WebhookService
	Observer     *ObserverService
	customCAFile string
}

type service struct {
	client *Client
}

// ClientOption interface to specify options for client
type ClientOption func(*Client)

// WithCustomCAFile Set a custom CA file to be used instead of the
// system-poool CAs.
func WithCustomCAFile(customCAFile string) ClientOption {
	return func(c *Client) {
		c.customCAFile = customCAFile
	}
}

// WithCustomEndpoint Set an endpoint which differs from the official
// api endpoint.
func WithCustomEndpoint(endpoint string) ClientOption {
	return func(c *Client) {
		c.BaseURL, _ = url.Parse(endpoint)
	}
}

// NewClient creates a new client in order to send requests to the eyeson API.
func NewClient(key string, options ...ClientOption) (*Client, error) {
	baseURL, _ := url.Parse(endpoint)
	c := &Client{apiKey: key, BaseURL: baseURL, client: http.DefaultClient}

	for _, opt := range options {
		opt(c)
	}

	// load customCAFile here if set
	if len(c.customCAFile) > 0 {

		// Load CA cert
		caCert, err := ioutil.ReadFile(c.customCAFile)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("Failed to append CA")
		}
		tlsConfig := &tls.Config{
			RootCAs: caCertPool,
		}
		tr := &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    5 * time.Second,
			DisableCompression: true,
			TLSClientConfig:    tlsConfig,
		}
		c.client = &http.Client{
			Transport: tr,
		}
	}

	c.Rooms = &RoomsService{c}
	c.Webhook = &WebhookService{c}
	c.Observer = &ObserverService{c}
	return c, nil
}

// UserClient provides a client for user requests that use the session access
// key for authorization.
func (c *Client) UserClient() *Client {
	return &Client{BaseURL: c.BaseURL, client: c.client}
}

// NewRequest prepares a request to be sent to the API.
func (c *Client) NewRequest(method, urlStr string, data url.Values) (*http.Request, error) {
	u := c.BaseURL.JoinPath(urlStr)

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

// NewPlainRequest create a request with bytes and content-type
func (c *Client) NewPlainRequest(method, urlStr string, data *bytes.Buffer, contentType string) (*http.Request, error) {
	u := c.BaseURL.JoinPath(urlStr)

	req, err := http.NewRequest(method, u.String(), data)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
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
	case c == 200 || c == 201 || c == 204:
		return nil
	case c == 404:
		return errors.New("Not found! Resource does not exist or expired")
	case c == 401:
		return errors.New("Authorization failed! Check the API key to be valid")
	case c == 403:
		return errors.New("Bad request! Check your request parameters to be valid")
	default:
		return fmt.Errorf("Unknown error! Request failed for an unknown error (%d)", c)
	}
}
