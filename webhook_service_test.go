package eyeson

import (
	"net/http"
	"testing"
)

func TestWebhookRegister(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/webhooks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{"url": "https://example.com/webhook-listener",
			"types": "room_update"})
		w.WriteHeader(201)
	})

	if err := client.Webhook.Register("https://example.com/webhook-listener", WEBHOOK_ROOM); err != nil {
		t.Errorf("Webhook register not successfull: %v", err)
	}
}

func TestWebhookUnregister(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/webhooks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(200)
		w.Write([]byte("{\"id\":\"42\"}"))
	})
	mux.HandleFunc("/webhooks/42", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(204)
	})

	if err := client.Webhook.Unregister(); err != nil {
		t.Errorf("Webhook unregister not successfull: %v", err)
	}
}
