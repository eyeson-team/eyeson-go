package main

import (
	"fmt"
	"os"
	"time"

	"github.com/eyeson-team/eyeson-go"
)

func main() {
	client, _ := eyeson.NewClient(os.Getenv("API_KEY"))
	room, _ := client.Rooms.Join("test-room", "test-user-name",
		map[string]string{"options[widescreen]": "true",
			"options[sfu_mode]": "disabled"})
	room.WaitReady()

	if err := room.CreateSnapshot(); err != nil {
		fmt.Println("Failed to create snapshot: ", err)
	}
	time.Sleep(2 * time.Second)
	if err := room.CreateSnapshot(); err != nil {
		fmt.Println("Failed to create snapshot: ", err)
	}

	// retrieve all snapshots
	snapshots, err := client.Rooms.GetSnapshots(room.Data.Room.ID, nil)
	if err != nil {
		fmt.Println("Failed to retrieve snapshot: ", err)
	}
	for _, s := range *snapshots {
		fmt.Printf("Snapshot id %s download_url %s\n", s.ID, *s.Links.Download)
	}
}
