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
	if _, err := srv.client.Do(req, &room); err != nil {
		return nil, err
	}
	return &UserService{Data: room, client: srv.client.UserClient()}, nil
}
