package main

import (
	"fmt"
	"os"

	"github.com/eyeson-team/eyeson-go"
)

func main() {
	client, _ := eyeson.NewClient(os.Getenv("API_KEY"))
	room, _ := client.Rooms.Join("test-room", "test-user-name",
		map[string]string{"options[widescreen]": "true",
			"options[sfu_mode]": "disabled"})
	room.WaitReady()
	fmt.Println("Join via:", room.Data.Links.Gui)
	room.SetLayer("https://docs.eyeson.com/img/examples/bg_p1.png",
		eyeson.Background)
}
