package eyeson

import (
	"fmt"
	"net/http"
	"testing"
)

func TestUserService_Chat(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{"user[name]": "mike@eyeson.team"})
		fmt.Fprint(w, `{"access_key":"token","links":{"gui": "https://app.eyeson.team/?token"}}`)
	})

	content := "/me sending a message"
	mux.HandleFunc("/rooms/token/messages", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{"type": "chat", "content": content})
		fmt.Fprint(w, `{}`)
	})

	user, err := client.Rooms.Join("", "mike@eyeson.team", nil)
	if err != nil {
		t.Errorf("RoomsService Join not successfull, got %v", err)
	}

	if err = user.Chat(content); err != nil {
		t.Errorf("UserService could not send chat message, got %v", err)
	}
}

func TestUserService_StartRecording(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{"user[name]": "mike@eyeson.team"})
		fmt.Fprint(w, `{"access_key":"token","links":{"gui": "https://app.eyeson.team/?token"}}`)
	})

	mux.HandleFunc("/rooms/token/recording", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{}`)
	})

	user, err := client.Rooms.Join("", "mike@eyeson.team", nil)
	if err != nil {
		t.Errorf("RoomsService Join not successfull, got %v", err)
	}

	if err = user.StartRecording(); err != nil {
		t.Errorf("UserService could not start a recording, got %v", err)
	}
}

func TestUserService_StopRecording(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{"user[name]": "mike@eyeson.team"})
		fmt.Fprint(w, `{"access_key":"token","links":{"gui": "https://app.eyeson.team/?token"}}`)
	})

	mux.HandleFunc("/rooms/token/recording", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{}`)
	})

	user, err := client.Rooms.Join("", "mike@eyeson.team", nil)
	if err != nil {
		t.Errorf("RoomsService Join not successfull, got %v", err)
	}

	if err = user.StopRecording(); err != nil {
		t.Errorf("UserService could not stop a recording, got %v", err)
	}
}

func TestUserService_StartBroadcast(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{"user[name]": "mike@eyeson.team"})
		fmt.Fprint(w, `{"access_key":"token","links":{"gui": "https://app.eyeson.team/?token"}}`)
	})

	streamUrl := "youtube/stream@url"
	mux.HandleFunc("/rooms/token/broadcasts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{"stream_url": streamUrl})
		fmt.Fprint(w, `{}`)
	})

	user, err := client.Rooms.Join("", "mike@eyeson.team", nil)
	if err != nil {
		t.Errorf("RoomsService Join not successfull, got %v", err)
	}

	if err = user.StartBroadcast(streamUrl); err != nil {
		t.Errorf("UserService could not start a broadcast, got %v", err)
	}
}

func TestUserService_StopBroadcast(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{"user[name]": "mike@eyeson.team"})
		fmt.Fprint(w, `{"access_key":"token","links":{"gui": "https://app.eyeson.team/?token"}}`)
	})

	mux.HandleFunc("/rooms/token/broadcasts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{}`)
	})

	user, err := client.Rooms.Join("", "mike@eyeson.team", nil)
	if err != nil {
		t.Errorf("RoomsService Join not successfull, got %v", err)
	}

	if err = user.StopBroadcast(); err != nil {
		t.Errorf("UserService could not stop a broadcast, got %v", err)
	}
}

func TestUserService_SetLayout(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{"user[name]": "mike@eyeson.team"})
		fmt.Fprint(w, `{"access_key":"token","links":{"gui": "https://app.eyeson.team/?token"}}`)
	})

	mux.HandleFunc("/rooms/token/layout", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{"layout": "custom", "users[0]": "first", "users[1]": "second",
			"voice_activation": "false", "show_names": "true"})
		fmt.Fprint(w, `{}`)
	})

	user, err := client.Rooms.Join("", "mike@eyeson.team", nil)
	if err != nil {
		t.Errorf("RoomsService Join not successfull, got %v", err)
	}

	users := []string{"first", "second"}
	if err = user.SetLayout("custom", users, false, true); err != nil {
		t.Errorf("UserService could not set layout, got %v", err)
	}
}

func TestUserService_SetLayer(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{"user[name]": "mike@eyeson.team"})
		fmt.Fprint(w, `{"access_key":"token","links":{"gui": "https://app.eyeson.team/?token"}}`)
	})

	imgUrl := "https://eyeson.com/overlay.png"
	mux.HandleFunc("/rooms/token/layers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{"url": imgUrl, "z-index": "1"})
		fmt.Fprint(w, `{}`)
	})

	user, err := client.Rooms.Join("", "mike@eyeson.team", nil)
	if err != nil {
		t.Errorf("RoomsService Join not successfull, got %v", err)
	}

	if err = user.SetLayer(imgUrl, Foreground); err != nil {
		t.Errorf("UserService could not set layer, got %v", err)
	}
}

func TestUserService_ClearLayer(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{"user[name]": "mike@eyeson.team"})
		fmt.Fprint(w, `{"access_key":"token","links":{"gui": "https://app.eyeson.team/?token"}}`)
	})

	mux.HandleFunc("/rooms/token/layers/-1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{}`)
	})

	user, err := client.Rooms.Join("", "mike@eyeson.team", nil)
	if err != nil {
		t.Errorf("RoomsService Join not successfull, got %v", err)
	}

	if err = user.ClearLayer(Background); err != nil {
		t.Errorf("UserService could not clear background layer, got %v", err)
	}
}

func TestUserService_StartPlayback(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{"user[name]": "mike@eyeson.team"})
		fmt.Fprint(w, `{"access_key":"token","links":{"gui": "https://app.eyeson.team/?token"}}`)
	})

	vidUrl := "https://eyeson.com/playback.mp4"
	mux.HandleFunc("/rooms/token/playbacks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{"playback[url]": vidUrl, "playback[replacement_id]": "first"})
		fmt.Fprint(w, `{}`)
	})

	user, err := client.Rooms.Join("", "mike@eyeson.team", nil)
	if err != nil {
		t.Errorf("RoomsService Join not successfull, got %v", err)
	}

	if err = user.StartPlayback(vidUrl, "first"); err != nil {
		t.Errorf("UserService could not start playback, got %v", err)
	}
}

func TestUserService_StopMeeting(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{"user[name]": "mike@eyeson.team"})
		fmt.Fprint(w, `{"access_key":"token","links":{"gui": "https://app.eyeson.team/?token"}}`)
	})

	mux.HandleFunc("/rooms/token", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{}`)
	})

	user, err := client.Rooms.Join("", "mike@eyeson.team", nil)
	if err != nil {
		t.Errorf("RoomsService Join not successfull, got %v", err)
	}

	if err = user.StopMeeting(); err != nil {
		t.Errorf("UserService could not stop a meeting, got %v", err)
	}
}
