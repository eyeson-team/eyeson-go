package eyeson

import (
	"time"
)

const WEBHOOK_ROOM string = "room_update"
const WEBHOOK_RECORDING string = "recording_update"

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
	} `json:"recording,omitempty"`
	Room struct {
		Id        string    `json:"id"`
		Name      string    `json:"name"`
		StartedAt time.Time `json:"started_at"`
		Shutdown  bool      `json:"shutdown"`
	} `json:"room"`
}
