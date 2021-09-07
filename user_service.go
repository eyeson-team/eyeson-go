package eyeson

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Timeout provides the maximum number of seconds WaitReady will wait for a
// meeting and user to be ready.
const Timeout int = 180

// Background provides the z-index to represent a background image.
const Background int = -1

// Foreground provides the z-index to represent a foreground image.
const Foreground int = 1

// RoomsService provides method Join to start and join a room.
type UserService struct {
	client *Client
	Data   *RoomResponse
}

// WaitReady waits until a meeting has successfully been started. It has a
// fixed polling interval of one second. WaitReady responds with an error on
// timeout or any communication problems.
func (u *UserService) WaitReady() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Timeout)*time.Second)
	defer cancel()

	res := make(chan error)
	go func() {
		for u.Data.Ready == false {
			time.Sleep(1 * time.Second)
			if err := u.updateRoomData(); err != nil {
				res <- err
				break
			}
			if u.Data.Room.Shutdown {
				res <- errors.New("Meeting has been shutdown.")
				break
			}
		}
		close(res)
	}()
	for {
		select {
		case err := <-res:
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (u *UserService) updateRoomData() error {
	path := "/rooms/" + u.Data.AccessKey
	req, err := u.client.NewRequest(http.MethodGet, path, url.Values{})
	if err != nil {
		return err
	}
	resp, err := u.client.Do(req, &u.Data)
	if err != nil {
		return err
	}
	return validateResponse(resp)
}

// Chat sends a chat message.
func (u *UserService) Chat(content string) error {
	data := url.Values{}
	data.Set("type", "chat")
	data.Set("content", content)
	path := "/rooms/" + u.Data.AccessKey + "/messages"
	req, err := u.client.NewRequest(http.MethodPost, path, data)
	if err != nil {
		return err
	}
	resp, err := u.client.Do(req, nil)
	if err != nil {
		return err
	}
	return validateResponse(resp)
}

// StartRecording starts a recording.
func (u *UserService) StartRecording() error {
	path := "/rooms/" + u.Data.AccessKey + "/recording"
	req, err := u.client.NewRequest(http.MethodPost, path, nil)
	if err != nil {
		return err
	}
	resp, err := u.client.Do(req, nil)
	if err != nil {
		return err
	}
	return validateResponse(resp)
}

// StopRecording stops a recording.
func (u *UserService) StopRecording() error {
	path := "/rooms/" + u.Data.AccessKey + "/recording"
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

// StartBroadcast starts a broadcast to the given stream url given by a
// streaming service like YouTube, Vimeo, and others.
func (u *UserService) StartBroadcast(streamUrl string) error {
	data := url.Values{}
	data.Set("stream_url", streamUrl)
	path := "/rooms/" + u.Data.AccessKey + "/broadcasts"
	req, err := u.client.NewRequest(http.MethodPost, path, data)
	if err != nil {
		return err
	}
	resp, err := u.client.Do(req, nil)
	if err != nil {
		return err
	}
	return validateResponse(resp)
}

// StopBroadcast stops a broadcast.
func (u *UserService) StopBroadcast() error {
	path := "/rooms/" + u.Data.AccessKey + "/broadcasts"
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

// SetLayout sets a participant podium layout where the layout is either
// "custom" or "auto". The users list is of user-ids or empty strings for empty
// participant positions. The flag voiceActivation replaces participants
// actively by voice detection. The flag showNames show or hides participant
// name overlays.
func (u *UserService) SetLayout(layout string, users []string, voiceActivation, showNames bool) error {
	data := url.Values{}
	if layout == "custom" {
		data.Set("layout", "custom")
	} else {
		data.Set("layout", "auto")
	}
	for i, userId := range users {
		data.Set("users["+strconv.Itoa(i)+"]", userId)
	}
	if voiceActivation {
		data.Set("voice_activation", "true")
	} else {
		data.Set("voice_activation", "false")
	}
	if showNames {
		data.Set("show_names", "true")
	} else {
		data.Set("show_names", "false")
	}
	path := "/rooms/" + u.Data.AccessKey + "/layout"
	req, err := u.client.NewRequest(http.MethodPost, path, data)
	if err != nil {
		return err
	}
	resp, err := u.client.Do(req, nil)
	if err != nil {
		return err
	}
	return validateResponse(resp)
}

// SetLayer sets a layer image using the given public available URL pointing to
// an image file. The z-index should be set using the constants Foreground or
// Background.
func (u *UserService) SetLayer(imgUrl string, zIndex int) error {
	data := url.Values{}
	data.Set("url", imgUrl)
	if zIndex == 1 {
		data.Set("z-index", "1")
	} else {
		data.Set("z-index", "-1")
	}
	path := "/rooms/" + u.Data.AccessKey + "/layers"
	req, err := u.client.NewRequest(http.MethodPost, path, data)
	if err != nil {
		return err
	}
	resp, err := u.client.Do(req, nil)
	if err != nil {
		return err
	}
	return validateResponse(resp)
}

// ClearLayer clears a layer given by the z-index that should be set using
// the constants Foreground or Background.
func (u *UserService) ClearLayer(zIndex int) error {
	path := "/rooms/" + u.Data.AccessKey + "/layers/"
	if zIndex == 1 {
		path += "1"
	} else {
		path += "-1"
	}
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

// StartPlayback starts a playback using the given public available URL to a
// video file. The given user id marks the position of the participant that
// is going to be replaced while the playback is shown.
func (u *UserService) StartPlayback(playbackUrl string, userId string) error {
	data := url.Values{}
	data.Set("playback[url]", playbackUrl)
	data.Set("playback[replacement_id]", userId)
	path := "/rooms/" + u.Data.AccessKey + "/playbacks"
	req, err := u.client.NewRequest(http.MethodPost, path, data)
	if err != nil {
		return err
	}
	resp, err := u.client.Do(req, nil)
	if err != nil {
		return err
	}
	return validateResponse(resp)
}

// StopMeeting stops a meeting for all participants.
func (u *UserService) StopMeeting() error {
	path := "/rooms/" + u.Data.AccessKey
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
