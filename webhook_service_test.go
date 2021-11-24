package eyeson

import (
	"fmt"
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
		fmt.Fprint(w, "")
	})

	if err := client.Webhook.Register("https://example.com/webhook-listener", WEBHOOK_ROOM); err != nil {
		t.Errorf("Webhook register nnot successfull, got error %v", err)
	}
}

func TestWebhookUnregister(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/webhooks/42", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, "")
	})

	if err := client.Webhook.Unregister("42"); err != nil {
		t.Errorf("Webhook unregister nnot successfull, got error %v", err)
	}
}
