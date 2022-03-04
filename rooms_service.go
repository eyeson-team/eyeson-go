package eyeson

import (
	"net/http"
	"net/url"
)

// RoomsService provides method Join to start and join a room.
type RoomsService service

// Join starts and joins a video meeting. The user string represents a unique
// identifier, to join as participant. If the same room identifier is provided
// the participants will join the same meeting. The room identifier can be
// omitted, the eyeson-api will therefor create a new room for every user
// joining.
func (srv *RoomsService) Join(id string, user string, options map[string]string) (*UserService, error) {
	data := url.Values{}
	if id != "" {
		data.Set("id", id)
	}
	data.Set("user[name]", user)
	for k, v := range options {
		data.Set(k, v)
	}
	req, err := srv.client.NewRequest(http.MethodPost, "/rooms", data)
	if err != nil {
		return nil, err
	}

	var room *RoomResponse
	if _, err = srv.client.Do(req, &room); err != nil {
		return nil, err
	}
	return &UserService{Data: room, client: srv.client.UserClient()}, nil
}

// GuestJoin creates a new guest user for an active meeting.
func (srv *RoomsService) GuestJoin(guestToken, id, name, avatar string) (*UserService, error) {
	data := url.Values{}
	data.Set("name", name)
	if id != "" {
		data.Set("id", id)
	}
	if avatar != "" {
		data.Set("avatar", avatar)
	}
	req, err := srv.client.NewRequest(http.MethodPost, "/guests/"+guestToken, data)
	if err != nil {
		return nil, err
	}

	var room *RoomResponse
	if _, err = srv.client.Do(req, &room); err != nil {
		return nil, err
	}
	return &UserService{Data: room, client: srv.client.UserClient()}, nil
}

// Shutdown force stops a running meeting.
func (srv *RoomsService) Shutdown(id string) error {
	req, err := srv.client.NewRequest(http.MethodDelete, "/rooms/"+id, nil)
	if err != nil {
		return err
	}
	_, err = srv.client.Do(req, nil)
	if err != nil {
		return err
	}
	return nil
}
