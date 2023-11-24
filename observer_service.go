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

// EventInterface interface for all event-messages.
type EventInterface interface {
	GetType() string
}

// EventBase Base for all events. Has only the type field.
type EventBase struct {
	Type string `json:"type"`
}

// GetType retrieve the type. Implements the EventInterface
func (msg *EventBase) GetType() string {
	return msg.Type
}

// Options Struct containing list of options of the room.
type Options struct {
	ShowNames bool `json:"show_names"`
}

// EventRoom represents the room.
type EventRoom struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Ready        bool          `json:"ready"`
	StartedAt    time.Time     `json:"started_at"`
	Shutdown     bool          `json:"shutdown"`
	GuestToken   string        `json:"guest_token"`
	Options      Options       `json:"options"`
	Participants []Participant `json:"participants"`
	Broadcasts   []Broadcast   `json:"broadcasts"`
}

// RoomUpdate event is sent if any of the room properties is changed.
type RoomUpdate struct {
	EventBase
	Content EventRoom `json:"content"`
}

// Participant holds informationof a participating user containing its online status.
type Participant struct {
	ID     string `json:"id"`
	RoomID string `json:"room_id"`
	Name   string `json:"name"`
	Guest  bool   `json:"guest"`
	Online bool   `json:"online"`
	Avatar string `json:"avatar"`
}

// ParticipantUpdate event is sent whenever the list of participants changes.
type ParticipantUpdate struct {
	EventBase
	Participant Participant `json:"participant"`
}

// PodiumPosition defines an area on the podium wich belongs to the specified user identified
// by its user-id.
type PodiumPosition struct {
	UserID string  `json:"user_id"`
	PlayID *string `json:"play_id"`
	Width  int     `json:"width"`
	Height int     `json:"height"`
	Left   int     `json:"left"`
	Top    int     `json:"top"`
	ZIndex int     `json:"z-index"`
}

// PodiumUpdate event is sent whenever the podium layout or the positioning of participants
// changes.
type PodiumUpdate struct {
	EventBase
	Podium []PodiumPosition `json:"podium"`
}

// EventUser represents an user within the event messages.
type EventUser struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Guest    bool      `json:"guest"`
	Avatar   string    `json:"avatar"`
	JoinedAt time.Time `json:"joined_at"`
}

// Broadcast contains information of a live-stream broadcast.
type Broadcast struct {
	ID        string    `json:"id"`
	Platform  string    `json:"platform"`
	PlayerURL string    `json:"player_url"`
	User      EventUser `json:"user"`
	Room      EventRoom `json:"room"`
}

// BroadcastUpdate event is sent whenever the live-stream broadcasts changes.
type BroadcastUpdate struct {
	EventBase
	Broadcasts []Broadcast `json:"broadcasts"`
}

// Links holds a list of links for the corresponding resource.
type Links struct {
	Self     *string `json:"self"`
	Download *string `json:"download"`
}

// Recording holds information on a recording.
type Recording struct {
	ID string `json:"id"`
	// Unix timestamp
	CreatedAt int       `json:"created_at"`
	Duration  int       `json:"duration"`
	Links     Links     `json:"links"`
	User      EventUser `json:"user"`
	Room      EventRoom `json:"room"`
}

// RecordingUpdate event is sent whenver recording is started or stopped.
type RecordingUpdate struct {
	EventBase
	Recording Recording `json:"recording"`
}

// OptionsUpdate event is sent whenever an option parameter is modified
// or added via the rest-interface.
type OptionsUpdate struct {
	EventBase
	Options Options `json:"options"`
}

// Snapshot represents a snapshot
type Snapshot struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Links     Links     `json:"links"`
	Creator   EventUser `json:"creator"`
	CreatedAt time.Time `json:"created_at"`
	Room      EventRoom `json:"room"`
}

// SnapshotUpdate fired whenever a new snapshot ist taken.
type SnapshotUpdate struct {
	EventBase
	Snapshots []Snapshot `json:"snapshots"`
}

// Chat contains a chat message.
type Chat struct {
	EventBase
	Content   string    `json:"content"`
	ClientID  string    `json:"cid"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Playback represents a playback i.e. media inject into the confserver.
type Playback struct {
	URL    string `json:"url"`
	PlayID string `json:"play_id"`
	Audio  bool   `json:"audio"`
}

// PlaybackUpdate Sent when a playback was started.
type PlaybackUpdate struct {
	EventBase
	Playing Playback `json:"playing"`
}

var eventTypes = map[string]func() EventInterface{
	"room_update":        func() EventInterface { return &RoomUpdate{} },
	"participant_update": func() EventInterface { return &ParticipantUpdate{} },
	"podium_update":      func() EventInterface { return &PodiumUpdate{} },
	"recording_update":   func() EventInterface { return &RecordingUpdate{} },
	"broadcasts_update":  func() EventInterface { return &BroadcastUpdate{} },
	"options_update":     func() EventInterface { return &OptionsUpdate{} },
	"snapshots_update":   func() EventInterface { return &SnapshotUpdate{} },
	"playback_update":    func() EventInterface { return &PlaybackUpdate{} },
	"chat":               func() EventInterface { return &Chat{} },
}

// Connect connects the observer and returns an eventInterface channel on success.
func (os *ObserverService) Connect(ctx context.Context, roomID string) (<-chan EventInterface, error) {
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

	msgChan := make(chan EventInterface, 1)

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

				var msgBase EventBase
				err := json.Unmarshal(ev.Event.Message, &msgBase)
				if err != nil {
					log.Printf("Failed to unmarshal [%s].\n", err)
					continue
				}
				//log.Printf("message received. [%s].\n", ev.Event.Message)
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
