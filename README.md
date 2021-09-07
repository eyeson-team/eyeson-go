
# eyeson go library

A golang client for the [eyeson video conferencing
API](https://eyeson-team.github.io/api/api-reference/).

## Usage

```golang
// Get your api-key at https://eyeson-team.github.io/api
client := eyeson.NewClient(eyesonApiKey)
room, err := client.Rooms.Join("standup meeting", "mike")
room.Links.Gui // https://app.eyeson.team/?sessionToken URL to eyeson web GUI
err = room.WaitReady()
overlayUrl = "https://eyeson-team.github.io/api/images/eyeson-overlay.png"
// Set a foreground image.
err = room.SetLayer(overlayUrl, eyeson.Foreground)
// Send a chat message.
err = room.Chat("Welcome!")
```

## Development

```sh
make test # run go tests
# run an example program that starts a meeting, adds an overlay and sends
# a chat message.
make run API_KEY=...
```

## Releases

- 1.0.0 Initial Release
