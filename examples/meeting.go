package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	eyeson "github.com/eyeson-team/eyeson-go"
)

const overlayURL string = "https://eyeson-team.github.io/api/images/eyeson-overlay.png"

func main() {
	userName := flag.String("user", "gopher", "unique user name")
	roomName := flag.String("room", "demo", "room identifier")
	flag.Parse()

	client, err := eyeson.NewClient(os.Getenv("API_KEY"))
	if err != nil {
		log.Fatal(err)
	}
	room, err := client.Rooms.Join(*roomName, *userName, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("Join the room at %q", room.Data.Links.GuestJoin)

	setLogo(room)
}

func setLogo(room *eyeson.UserService) {
	fmt.Println("Wait until the room is ready")
	if err := room.WaitReady(); err != nil {
		fmt.Printf("Cannot determine room ready status: %v", err)
	}
	fmt.Println("Room is ready, set Layer")
	if err := room.SetLayer(overlayURL, eyeson.Foreground, nil); err != nil {
		fmt.Printf("Failed to set overlay url: %v", err)
	}
	fmt.Println("Send welcome message")
	if err := room.Chat("Welcome aboard!"); err != nil {
		fmt.Printf("Failed to send chat message: %v", err)
	}
}
