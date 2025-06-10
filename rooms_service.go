package eyeson

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
	mediaTypesStrings := []string{}
	for _, m := range mediaTypes {
		mediaTypesStrings = append(mediaTypesStrings, string(m))
	}
	data.Set("type", strings.Join(mediaTypesStrings, ","))

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
func (srv *RoomsService) GetSnapshot(snapshotID string) (*Snapshot, error) {
	path := "/snapshots/" + snapshotID
	req, err := srv.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var snapshot Snapshot
	resp, err := srv.client.Do(req, &snapshot)
	if err != nil {
		return nil, err
	}
	return &snapshot, validateResponse(resp)
}

// GetSnaphostsOptions options supporting pagination and filtering.
type GetSnaphostsOptions struct {
	Page      *int
	StartedAt *time.Time
	Since     *time.Time
	Until     *time.Time
}

// GetSnapshots retrieves a a list of snapshots for a room.
func (srv *RoomsService) GetSnapshots(ID string, options *GetSnaphostsOptions) (*[]Snapshot, error) {
	data := url.Values{}
	if options != nil {
		if options.Page != nil {
			data.Set("page", strconv.Itoa(*options.Page))
		}
		if options.StartedAt != nil {
			data.Set("started_at", options.StartedAt.Format(time.RFC3339))
		}
		if options.Since != nil {
			data.Set("since", options.Since.Format(time.RFC3339))
		}
		if options.Until != nil {
			data.Set("until", options.Until.Format(time.RFC3339))
		}
	}
	path := "/rooms/" + ID + "/snapshots"
	req, err := srv.client.NewRequest(http.MethodGet, path, data)
	if err != nil {
		return nil, err
	}
	var snapshots []Snapshot
	resp, err := srv.client.Do(req, &snapshots)
	if err != nil {
		return nil, err
	}
	return &snapshots, validateResponse(resp)
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

// GetRecording retrieves a recording.
func (srv *RoomsService) GetRecording(recordingID string) (*Recording, error) {
	path := "/recordings/" + recordingID
	req, err := srv.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var recording Recording
	resp, err := srv.client.Do(req, &recording)
	if err != nil {
		return nil, err
	}
	return &recording, validateResponse(resp)
}

// GetRecordingsOptions options supporting pagination and filtering.
type GetRecordingsOptions struct {
	Page      *int
	StartedAt *time.Time
	Since     *time.Time
	Until     *time.Time
}

// GetSnapshots retrieves a a list of recordings for a room.
func (srv *RoomsService) GetRecordings(ID string, options *GetRecordingsOptions) (*[]Recording, error) {
	data := url.Values{}
	if options != nil {
		if options.Page != nil {
			data.Set("page", strconv.Itoa(*options.Page))
		}
		if options.StartedAt != nil {
			data.Set("started_at", options.StartedAt.Format(time.RFC3339))
		}
		if options.Since != nil {
			data.Set("since", options.Since.Format(time.RFC3339))
		}
		if options.Until != nil {
			data.Set("until", options.Until.Format(time.RFC3339))
		}
	}
	path := "/rooms/" + ID + "/recordings"
	req, err := srv.client.NewRequest(http.MethodGet, path, data)
	if err != nil {
		return nil, err
	}
	var recordings []Recording
	resp, err := srv.client.Do(req, &recordings)
	if err != nil {
		return nil, err
	}
	return &recordings, validateResponse(resp)
}

// DeleteRecording deletes a recording.
func (u *RoomsService) DeleteRecording(recordingID string) error {
	path := "/recordings/" + recordingID
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

func (srv *RoomsService) GetCurrentMeetings() (*[]RoomInfo, error) {
	req, err := srv.client.NewRequest(http.MethodGet, "/rooms", nil)
	if err != nil {
		return nil, err
	}

	var rooms []RoomInfo
	resp, err := srv.client.Do(req, &rooms)
	if err != nil {
		return nil, err
	}
	return &rooms, validateResponse(resp)
}

func (srv *RoomsService) GetRoomUsers(ID string, online *bool) (*[]Participant, error) {
	path := "/rooms/" + ID + "/users"
	if online != nil {
		path += "online=" + strconv.FormatBool(*online)
	}
	req, err := srv.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var users []Participant
	resp, err := srv.client.Do(req, &users)
	if err != nil {
		return nil, err
	}
	return &users, validateResponse(resp)
}
