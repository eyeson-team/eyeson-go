package eyeson

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// WebhookService provides method Register and Unregister a webhook.
type WebhookService service

// Register will assign an endpoint URL to the current ApiKey.
func (srv *WebhookService) Register(endpoint, types string) error {
	data := url.Values{}
	data.Set("url", endpoint)
	data.Set("types", types)
	req, err := srv.client.NewRequest(http.MethodPost, "/webhooks", data)
	if err != nil {
		return err
	}

	res, err := srv.client.Do(req, nil)
	if res.StatusCode != http.StatusCreated {
		return errors.New(fmt.Sprintf("Bad API status code 201, got %d", res.StatusCode))
	}
	return err
}

// Get provides details about a registered webhook.
func (srv *WebhookService) Get() (*WebhookDetails, error) {
	req, err := srv.client.NewRequest(http.MethodGet, "/webhooks", nil)
	if err != nil {
		return nil, err
	}

	var details WebhookDetails
	res, err := srv.client.Do(req, &details)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Bad API status code 200, got %d", res.StatusCode))
	}
	return &details, err
}

// Unregister will clear the current webhook.
func (srv *WebhookService) Unregister() error {
	w, err := srv.Get()
	if err != nil {
		return err
	}
	req, err := srv.client.NewRequest(http.MethodDelete, "/webhooks/"+w.Id, nil)
	if err != nil {
		return err
	}

	srv.client.Do(req, nil)
	// res, err := srv.client.Do(req, nil)
	// if err != nil {
	// 	return err
	// }
	// if res.StatusCode != http.StatusNoContent {
	// 	return errors.New(fmt.Sprintf("Bad API status code 204, got %d", res.StatusCode))
	// }
	return nil
}
