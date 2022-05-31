package eyeson

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

const WEBHOOK_ROOM string = "room_update"
const WEBHOOK_RECORDING string = "recording_update"
const WEBHOOK_SNAPSHOT string = "snapshot_update"

// WebhookDetails provide configuration details.
type WebhookDetails struct {
	Id                string    `json:"id"`
	Url               string    `json:"url"`
	Types             []string  `json:"types"`
	LastRequestSentAt time.Time `json:"last_request_sent_at"`
	LastResponseCode  string    `json:"last_response_code"`
}

// Webhook holds available attributes to a room.
type Webhook struct {
	Type      string `json:"type"`
	Recording struct {
		Id        string `json:"id"`
		Duration  int    `json:"duration"`
		CreatedAt int    `json:"created_at"`
		Links     struct {
			Download string `json:"download"`
		} `json:"links"`
		Room struct {
			Id string `json:"id"`
		}
	} `json:"recording,omitempty"`
	Room struct {
		Id        string    `json:"id"`
		Name      string    `json:"name"`
		StartedAt time.Time `json:"started_at"`
		Shutdown  bool      `json:"shutdown"`
	} `json:"room"`
	Snapshot struct {
		Id        string    `json:"id"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
		Links     struct {
			Download string `json:"download"`
		} `json:"links"`
		Room struct {
			Id string `json:"id"`
		}
	} `json:"snapshot,omitempty"`
}

func NewWebhook(apiKey string, r *http.Request) (*Webhook, error) {
	var webhook Webhook
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	h := hmac.New(sha256.New, []byte(apiKey))
	h.Write(raw)
	if hex.EncodeToString(h.Sum(nil)) != r.Header.Get("X-Eyeson-Signature") {
		return nil, errors.New("Webhook signature does not match")
	}
	if err = json.Unmarshal(raw, &webhook); err != nil {
		return nil, err
	}
	return &webhook, nil
}
