package eyeson

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Timeout provides the maximum number of seconds WaitReady will wait for a
// meeting and user to be ready.
const Timeout int = 180

// Background provides the z-index to represent a background image.
const Background int = -1

// Foreground provides the z-index to represent a foreground image.
const Foreground int = 1

// ImageType provides custom type for specifying image type
type ImageType string

// List of supported image types
const (
	Jpg  ImageType = "jpg"
	Png  ImageType = "png"
	Svg  ImageType = "svg"
	Webp ImageType = "webp"
)

// Layout provides a custom type for specifying layout configuration.
type Layout string

const (
	// Auto Automatically sets layouts according to the number of participants
	Auto Layout = "auto"
	// Custom Maintains manually assigned positions.
	Custom Layout = "custom"
)

// UserService provides methods a user can perform.
type UserService struct {
	client *Client
	Data   *RoomResponse
}

// NewUserServiceFromAccessKey Create a new UserService from an access-key.
func NewUserServiceFromAccessKey(accessKey string, options ...ClientOption) (*UserService, error) {
	client, err := NewClient("", options...)
	if err != nil {
		return nil, err
	}
	u := &UserService{
		client: client,
		Data: &RoomResponse{
			AccessKey: accessKey,
		},
	}
	return u, nil
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
				res <- errors.New("Meeting has been shutdown")
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
func (u *UserService) StartBroadcast(streamURL string) error {
	data := url.Values{}
	data.Set("stream_url", streamURL)
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

// LayoutObjectFit defines how content fits within its container in a layout.
type LayoutObjectFit string

const (
	// Cover scales the content to cover the entire container, potentially cropping some parts.
	Cover LayoutObjectFit = "cover"
	// Contain scales the content to fit within the container while maintaining aspect ratio.
	Contain LayoutObjectFit = "contain"
	// Autofit automatically determines the best fitting method for the content.
	Autofit LayoutObjectFit = "auto"
)

// LayoutPos represents the position and dimensions of a participant in a layout.
type LayoutPos struct {
	// X is the horizontal position coordinate.
	X int
	// Y is the vertical position coordinate.
	Y int
	// Width is the horizontal size of the position.
	Width int
	// Height is the vertical size of the position.
	Height int
	// ObjectFit determines how the participant's video fits within the assigned space.
	ObjectFit LayoutObjectFit
}

// LayoutMap contains the positions of participants in a custom layout configuration.
type LayoutMap struct {
	// Positions is a slice of participant position configurations.
	Positions []LayoutPos
}

func (lmap *LayoutMap) toString() string {
	serialMaps := []string{}
	for _, p := range lmap.Positions {
		serialMaps = append(serialMaps, fmt.Sprintf("[%d, %d, %d, %d, \"%s\"]", p.X, p.Y, p.Width, p.Height, p.ObjectFit))
	}
	return "[" + strings.Join(serialMaps, ",") + "]"
}

// AudioInsertConfig defines the configuration options for audio insertion.
type AudioInsertConfig string

const (
	// Enabled indicates that audio insert is shown all the time.
	Enabled AudioInsertConfig = "enabled"
	// Disabled indicates that audio insert is turned off.
	Disabled AudioInsertConfig = "disabled"
	// AudioOnly indicates that the insert is only shown if the participant is not shown on the podium.
	AudioOnly AudioInsertConfig = "audio_only"
)

// AudioInsertPosition represents the coordinates for positioning an audio insert visual element.
type AudioInsertPosition struct {
	// X is the horizontal position coordinate.
	X int
	// Y is the vertical position coordinate.
	Y int
}

// AudioInsert contains configuration for inserting audio into a meeting.
type AudioInsert struct {
	// Config specifies whether audio insertion is enabled, disabled, or audio-only.
	Config AudioInsertConfig
	// Position defines the visual position of the audio insert when enabled.
	// May be nil if no position is specified or for audio-only inserts.
	Position *AudioInsertPosition
}

// SetLayout sets a participant podium layout where the layout is either
// "custom" or "auto". The users list is of user-ids or empty strings for empty
// participant positions. The flag voiceActivation replaces participants
// actively by voice detection. The flag showNames show or hides participant
// name overlays.
func (u *UserService) SetLayout(layout Layout, users []string, voiceActivation, showNames bool,
	layoutName *string, layoutMap *LayoutMap, audioInsert *AudioInsert) error {
	data := url.Values{}
	if layout == "custom" {
		data.Set("layout", "custom")
	} else {
		data.Set("layout", "auto")
	}
	for _, userID := range users {
		data.Set("users[]", userID)
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
	if layoutName != nil {
		data.Set("name", *layoutName)
	}
	if layoutMap != nil {
		data.Set("map", layoutMap.toString())
	}
	if audioInsert != nil {
		data.Set("audio_insert", string(audioInsert.Config))
		if audioInsert.Position != nil {
			data.Set("audio_insert_position[x]", fmt.Sprint(audioInsert.Position.X))
			data.Set("audio_insert_position[y]", fmt.Sprint(audioInsert.Position.Y))
		}
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
func (u *UserService) SetLayer(imgURL string, zIndex int) error {
	data := url.Values{}
	data.Set("url", imgURL)
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

// SetLayerImage sets a layer image providing
// an image file. The z-index should be set using the constants Foreground or
// Background.
func (u *UserService) SetLayerImage(imgData []byte, imageType ImageType, zIndex int) error {

	body := &bytes.Buffer{}
	// Create a multipart writer
	writer := multipart.NewWriter(body)
	fileName := "layer-img."
	switch imageType {
	case Png, Jpg, Svg, Webp:
		fileName += string(imageType)
	default:
		return fmt.Errorf("Unsupported image type %s", imageType)
	}
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return err
	}
	_, err = io.Copy(part, bytes.NewReader(imgData))
	if err != nil {
		return err
	}

	if zIndex == 1 {
		writer.WriteField("z-index", "1")
	} else {
		writer.WriteField("z-index", "-1")
	}
	writer.Close()
	path := "/rooms/" + u.Data.AccessKey + "/layers"
	req, err := u.client.NewPlainRequest(http.MethodPost, path, body, writer.FormDataContentType())
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
func (u *UserService) StartPlayback(playbackURL string, userID string) error {
	data := url.Values{}
	data.Set("playback[url]", playbackURL)
	data.Set("playback[replacement_id]", userID)
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
