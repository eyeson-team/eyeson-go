package eyeson

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestRoomsService_Join(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{"user[name]": "mike@eyeson.team"})
		fmt.Fprint(w, `{"access_key":"token","links":{"gui": "https://app.eyeson.team/?token"}}`)
	})

	room, err := client.Rooms.Join("", "mike@eyeson.team", nil)
	if err != nil {
		t.Errorf("RoomsService Join not successfull, got %v", err)
	}

	want := &RoomResponse{AccessKey: "token", Links: RoomLinks{Gui: "https://app.eyeson.team/?token"}}
	if !reflect.DeepEqual(room.Data, want) {
		t.Errorf("RoomsService Join body = %v, want %v", room, want)
	}
}

func TestRoomsService_GuestJoin(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/guests/guest-token", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{"name": "guest@eyeson.team"})
		fmt.Fprint(w, `{"access_key":"token","links":{"gui": "https://app.eyeson.team/?token"}}`)
	})

	room, err := client.Rooms.GuestJoin("guest-token", "", "guest@eyeson.team", "")
	if err != nil {
		t.Errorf("RoomsService GuestJoin not successfull, got %v", err)
	}

	want := &RoomResponse{AccessKey: "token", Links: RoomLinks{Gui: "https://app.eyeson.team/?token"}}
	if !reflect.DeepEqual(room.Data, want) {
		t.Errorf("RoomsService GuestJoin body = %v, want %v", room, want)
	}
}

func TestRoomsService_Shutdown(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/rooms/seven", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(204)
	})

	err := client.Rooms.Shutdown("seven")
	if err != nil {
		t.Errorf("RoomsService Shutdown not successfull, got %v", err)
	}
}

func TestRoomsService_ForwardSource(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/rooms/room-id/forward/source", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{"forward_id": "fw-id", "user_id": "u-1",
			"type": "audio,video", "url": "https://dest.com"})
		w.WriteHeader(201)
	})

	err := client.Rooms.ForwardSource("room-id", "fw-id", "u-1", []MediaType{Audio, Video}, "https://dest.com")
	if err != nil {
		t.Errorf("RoomsService ForwardSource not successfull, got %v", err)
	}
}

func TestRoomsService_DeleteForward(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/rooms/room-id/forward/fw-id", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(204)
	})

	err := client.Rooms.DeleteForward("room-id", "fw-id")
	if err != nil {
		t.Errorf("RoomsService DeleteForward not successfull, got %v", err)
	}
}

func TestRoomsService_GetSnapshot(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/snapshots/snapshot-id", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":"snapshot-id","links":{"download": "https://fs.eyeson.com/meetings/snapshot-id"}}`)
	})

	snapshot, err := client.Rooms.GetSnapshot("snapshot-id")
	if err != nil {
		t.Errorf("RoomsService GetSnapshot not successfull, got %v", err)
	}
	downloadLink := "https://fs.eyeson.com/meetings/snapshot-id"
	want := &Snapshot{ID: "snapshot-id", Links: Links{Download: &downloadLink}}
	if !reflect.DeepEqual(snapshot, want) {
		t.Errorf("RoomsService GetSnapshot body = %v, want %v", snapshot, want)
	}
}

func TestRoomsService_GetSnapshots(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	since := "2025-06-10T10:00:00Z"
	mux.HandleFunc("/rooms/room-id/snapshots", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testQueryValues(t, r, values{"since": since})
		fmt.Fprint(w, `[{"id":"snapshot-id","links":{"download": "https://fs.eyeson.com/meetings/snapshot-id"}}]`)
	})
	sinceTime, _ := time.Parse(time.RFC3339, since)
	options := &GetSnaphostsOptions{
		Since: &sinceTime,
	}
	snapshots, err := client.Rooms.GetSnapshots("room-id", options)
	if err != nil {
		t.Errorf("RoomsService GetSnapshots not successfull, got %v", err)
	}
	if snapshots == nil || len(*snapshots) != 1 {
    	t.Errorf("RoomsService GetSnapshots body = %v", snapshots)
	}
	first := (*snapshots)[0]
	if  first.ID != "snapshot-id" {
		t.Errorf("RoomsService GetSnapshots body = %v", snapshots)
	}
	// downloadLink := "https://fs.eyeson.com/meetings/snapshot-id"
	// wantSnapshot := &Snapshot{ID: "snapshot-id", Links: Links{Download: &downloadLink}}
	// if !reflect.DeepEqual((*snapshots)[0], wantSnapshot) {
	// 	t.Errorf("RoomsService GetSnapshots body = %v, want %v", snapshots, wantSnapshot)
	// }
}
