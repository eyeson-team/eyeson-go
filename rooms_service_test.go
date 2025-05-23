package eyeson

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
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
