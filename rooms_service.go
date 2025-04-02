package eyeson

import (
	"net/http"
	"net/url"
	"strconv"
	"time"
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

// ForwardSource starts forwarding the userID media to the specified url.
func (srv *RoomsService) ForwardSource(id string, forwardID string, userID string,
	mediaTypes []MediaType, destURL string) error {
	data := url.Values{}
	data.Set("forward_id", forwardID)
	data.Set("user_id", userID)
	data.Set("url", destURL)
	for _, m := range mediaTypes {
		data.Set("type", string(m))
	}

	req, err := srv.client.NewRequest(http.MethodPost, "/rooms/"+id+"/forward/source", data)
	if err != nil {
		return err
	}
	_, err = srv.client.Do(req, nil)
	if err != nil {
		return err
	}
	return nil
}

// DeleteForward deletes a forward by its forwardID
func (srv *RoomsService) DeleteForward(id string, forwardID string) error {
	req, err := srv.client.NewRequest(http.MethodDelete, "/rooms/"+id+"/forward/"+forwardID, nil)
	if err != nil {
		return err
	}
	_, err = srv.client.Do(req, nil)
	if err != nil {
		return err
	}
	return nil
}

// GetSnapshot retrieves a snapshot.
func (u *RoomsService) GetSnapshot(snapshotID string) (*Snapshot, error) {
	path := "/snapshots/" + snapshotID
	req, err := u.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var snapshot Snapshot
	resp, err := u.client.Do(req, &snapshot)
	if err != nil {
		return nil, err
	}
	return &snapshot, validateResponse(resp)
}

type GetSnaphostsOptions struct {
	page      int
	startedAt time.Time
}

// GetSnapshots retrieves a a list of snapshots for a room.
func (u *RoomsService) GetSnapshots(ID string, options *GetSnaphostsOptions) (*[]Snapshot, error) {
	data := url.Values{}
	if options != nil {
		data.Set("page", strconv.Itoa(options.page))
		data.Set("started_at", options.startedAt.Format(time.RFC3339))
	}
	path := "/rooms/" + ID + "/snapshots"
	req, err := u.client.NewRequest(http.MethodGet, path, data)
	if err != nil {
		return nil, err
	}
	var snapshot []Snapshot
	resp, err := u.client.Do(req, &snapshot)
	if err != nil {
		return nil, err
	}
	return &snapshot, validateResponse(resp)
}

// DeleteSnapshot deletes a snapshot.
func (u *RoomsService) DeleteSnapshot(snapshotID string) error {
	path := "/snapshots/" + snapshotID
	req, err := u.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	resp, err := u.client.Do(req, nil)
	if err != nil {
		return err
	}
	return validateResponse(resp)
}
