package eyeson

import (
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

	_, err = srv.client.Do(req, nil)
	return err
}

// Get provides details about a registered webhook.
func (srv *WebhookService) Get() (*WebhookDetails, error) {
	req, err := srv.client.NewRequest(http.MethodGet, "/webhooks", nil)
	if err != nil {
		return nil, err
	}

	var details WebhookDetails
	_, err = srv.client.Do(req, &details)
	if err != nil {
		return nil, err
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

	_, err = srv.client.Do(req, nil)
	if err != nil {
		return err
	}
	return nil
}
