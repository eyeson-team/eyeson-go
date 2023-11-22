package eyeson

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	ac "github.com/bgentry/actioncable-go"
)

// ObserverService Service to listen and control a room.
type ObserverService service

type MsgInterface interface {
	GetType() string
}

type ObserverMessage struct {
	Type string `json:"type"`
}

func (msg *ObserverMessage) GetType() string {
	return msg.Type
}

type Options struct {
	ShowNames bool `json:"show_names"`
}

type RoomUpdateData struct {
	ID           string                  `json:"id"`
	Name         string                  `json:"name"`
	Ready        bool                    `json:"ready"`
	StartedAt    time.Time               `json:"started_at"`
	Shutdown     bool                    `json:"shutdown"`
	GuestToken   string                  `json:"guest_token"`
	Options      Options                 `json:"options"`
	Participants []ParticipantUpdateData `json:"participants"`
	/*
		Presentation string    `json:"presentation"`
		Broadcasts   []string  `json:"broadcasts"`
		Recording    string    `json:"recording"`
	*/
}

type RoomUpdate struct {
	ObserverMessage
	Content RoomUpdateData `json:"content"`
}

type ParticipantUpdateData struct {
	ID     string `json:"id"`
	RoomID string `json:"room_id"`
	Name   string `json:"name"`
	Guest  bool   `json:"guest"`
	Online bool   `json:"online"`
	//Avatar
}

type ParticipantUpdate struct {
	ObserverMessage
	Content ParticipantUpdateData `json:"content"`
}

var eventTypes = map[string]func() MsgInterface{
	"room_update":        func() MsgInterface { return &RoomUpdate{} },
	"participant_update": func() MsgInterface { return &ParticipantUpdate{} },
}

func (os *ObserverService) Connect(ctx context.Context, roomID string) (<-chan MsgInterface, error) {
	headerFunc := func() (*http.Header, error) {
		header := &http.Header{}
		header.Set("Authorization", os.client.apiKey)
		return header, nil
	}

	baseURL := os.client.BaseURL
	if baseURL == nil {
		return nil, fmt.Errorf("Client-BaseURL not specified")
	}

	wsURL := fmt.Sprintf("%s/rt?room_id=%s", baseURL, roomID)
	if strings.HasPrefix(wsURL, "https") {
		wsURL = strings.Replace(wsURL, "https", "wss", 1)
	} else if strings.HasPrefix(wsURL, "http") {
		wsURL = strings.Replace(wsURL, "http", "ws", 1)
	}

	cable := ac.NewClient(wsURL, headerFunc)
	ch, err := cable.Subscribe("RoomChannel")
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe: %s", err)
	}

	msgChan := make(chan MsgInterface, 1)

	go func() {
		for {
			select {
			case <-ctx.Done():
				// try shutting down the channel
				cable.Close()
				return
			case ev, ok := <-ch:
				if !ok {
					return // fmt.Errorf("Read from channel failed. Probably closed")
				}
				if ev.Err != nil {
					return //fmt.Errorf("Channel returned err: %s", ev.Err)
				}
				//fmt.Println("Message: ", ev.Event.Type, ev.Event.Data, string(ev.Event.Message))

				if len(ev.Event.Message) == 0 {
					continue
				}

				var msgBase ObserverMessage
				err := json.Unmarshal(ev.Event.Message, &msgBase)
				if err != nil {
					log.Printf("Failed to unmarshal [%s].\n", err)
					continue
				}
				log.Printf("message received. [%s].\n", ev.Event.Message)
				msgInitFunc, ok := eventTypes[msgBase.Type]
				if !ok {
					log.Printf("Message-type %s not supported.", msgBase.Type)
					continue
				}
				interf := msgInitFunc()
				err = json.Unmarshal(ev.Event.Message, interf)
				if err != nil {
					log.Println("Failed to unmarshal.")
					continue
				}

				msgChan <- interf
			}
		}
	}()

	return msgChan, nil
}
