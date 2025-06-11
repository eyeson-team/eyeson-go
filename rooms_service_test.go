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
	downloadLink := "https://fs.eyeson.com/meetings/snapshot-id"
	wantSnapshot := &Snapshot{ID: "snapshot-id", Links: Links{Download: &downloadLink}}
	first := &(*snapshots)[0]
	if !reflect.DeepEqual(first, wantSnapshot) {
		t.Errorf("RoomsService GetSnapshots body = %v, want %v", first, wantSnapshot)
	}
}

func TestRoomsService_DeleteSnapshot(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/snapshots/snapshot-id", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(204)
	})

	err := client.Rooms.DeleteSnapshot("snapshot-id")
	if err != nil {
		t.Errorf("RoomsService DeleteSnapshot not successfull, got %v", err)
	}
}

func TestRoomsService_GetRecording(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/recordings/recording-id", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":"recording-id","links":{"download": "https://fs.eyeson.com/meetings/snapshot-id"}}`)
	})

	recording, err := client.Rooms.GetRecording("recording-id")
	if err != nil {
		t.Errorf("RoomsService GetRecording not successfull, got %v", err)
	}
	downloadLink := "https://fs.eyeson.com/meetings/snapshot-id"
	want := &Recording{ID: "recording-id", Links: Links{Download: &downloadLink}}
	if !reflect.DeepEqual(recording, want) {
		t.Errorf("RoomsService GetRecording body = %v, want %v", recording, want)
	}
}

func TestRoomsService_GetRecordings(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	since := "2025-06-10T10:00:00Z"
	mux.HandleFunc("/rooms/room-id/recordings", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testQueryValues(t, r, values{"since": since})
		fmt.Fprint(w, `[{"id":"recording-id","links":{"download": "https://fs.eyeson.com/meetings/recording-id"}}]`)
	})
	sinceTime, _ := time.Parse(time.RFC3339, since)
	options := &GetRecordingsOptions{
		Since: &sinceTime,
	}
	recordings, err := client.Rooms.GetRecordings("room-id", options)
	if err != nil {
		t.Errorf("RoomsService GetRecordings not successfull, got %v", err)
	}
	if recordings == nil || len(*recordings) != 1 {
    	t.Errorf("RoomsService GetRecordings body = %v", recordings)
	}
	downloadLink := "https://fs.eyeson.com/meetings/recording-id"
	wantRecording := &Recording{ID: "recording-id", Links: Links{Download: &downloadLink}}
	first := &(*recordings)[0]
	if !reflect.DeepEqual(first, wantRecording) {
		t.Errorf("RoomsService GetRecordings body = %v, want %v", first, wantRecording)
	}
}

func TestRoomsService_DeleteRecording(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/recordings/recording-id", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(204)
	})

	err := client.Rooms.DeleteRecording("recording-id")
	if err != nil {
		t.Errorf("RoomsService DeleteRecording not successfull, got %v", err)
	}
}

func TestRoomsService_GetCurrentMeetings(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":"room-id","name":"test","ready":true,"started_at":"2025-06-10T10:00:00Z","shutdown":false,"guest_token":"abc"}]`)
	})
	rooms, err := client.Rooms.GetCurrentMeetings()
	if err != nil {
		t.Errorf("RoomsService GetCurrentMeetings not successfull, got %v", err)
	}
	if rooms == nil || len(*rooms) != 1 {
    	t.Errorf("RoomsService GetCurrentMeetings body = %v", rooms)
	}
	first := &(*rooms)[0]
	want := &RoomInfo{ID: "room-id", Name: "test", Ready: true, StartedAt: "2025-06-10T10:00:00Z", Shutdown: false, GuestToken: "abc"}
	if !reflect.DeepEqual(first, want) {
		t.Errorf("RoomsService GetCurrentMeetings body = %v, want %v", first, want)
	}
}

func TestRoomsService_GetRoomUsers(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/rooms/room-id/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testQueryValues(t, r, values{"online": "true"})
		fmt.Fprint(w, `[{"id":"user-id","room_id":"room-id","name":"test","online":true}]`)
	})
	online := true
	users, err := client.Rooms.GetRoomUsers("room-id", &online)
	if err != nil {
		t.Errorf("RoomsService GetRoomUsers not successfull, got %v", err)
	}
	if users == nil || len(*users) != 1 {
    	t.Errorf("RoomsService GetRoomUsers body = %v", users)
	}
	first := &(*users)[0]
	want := &Participant{ID: "user-id", RoomID: "room-id", Name: "test", Online: true}
	if !reflect.DeepEqual(first, want) {
		t.Errorf("RoomsService GetRoomUsers body = %v, want %v", first, want)
	}
}
