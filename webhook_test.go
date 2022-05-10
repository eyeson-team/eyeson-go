package eyeson

import (
	"encoding/json"
	"os"
	"testing"
)

func TestWebhookUnmarshalRoom(t *testing.T) {
	sample, err := os.ReadFile("./fixtures/webhook_room_update.json")
	if err != nil {
		t.Errorf("Failed to read fixture file: %v", err)
	}
	var webhook Webhook
	if err = json.Unmarshal(sample, &webhook); err != nil {
		t.Errorf("Failed to decode sample: %v", err)
	}
	if webhook.Type != "room_update" {
		t.Errorf("Expected type room_update, got %v", webhook.Type)
	}
	if webhook.Room.Id != "demo" {
		t.Errorf("Expected room identifier demo, got %v", webhook.Room.Id)
	}
	if webhook.Room.StartedAt.Weekday().String() != "Wednesday" {
		t.Errorf("Expected meeting from Wednesday, got %v", webhook.Room.StartedAt.Weekday())
	}
	if webhook.Room.Shutdown != false {
		t.Errorf("Expected meeting to be active, got %v", webhook.Room.Shutdown)
	}
}

func TestWebhookUnmarshalRecording(t *testing.T) {
	sample, err := os.ReadFile("./fixtures/webhook_recording_update.json")
	if err != nil {
		t.Errorf("Failed to read fixture file: %v", err)
	}
	var webhook Webhook
	if err = json.Unmarshal(sample, &webhook); err != nil {
		t.Errorf("Failed to decode sample: %v", err)
	}
	if webhook.Type != "recording_update" {
		t.Errorf("Expected type recording_update, got %v", webhook.Type)
	}
	if webhook.Recording.Duration != 2 {
		t.Errorf("Expected recording duration of two seconds, got %v", webhook.Recording.Duration)
	}
	if webhook.Recording.Room.Id != "demo" {
		t.Errorf("Expected room identifier demo, got %v", webhook.Recording.Room.Id)
	}
}
