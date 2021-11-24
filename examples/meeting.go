package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	eyeson "github.com/eyeson-team/eyeson-go"
)

const OverlayUrl string = "https://eyeson-team.github.io/api/images/eyeson-overlay.png"

func main() {
	userName := flag.String("user", "gopher", "unique user name")
	roomName := flag.String("room", "demo", "room identifier")
	flag.Parse()

	client := eyeson.NewClient(os.Getenv("API_KEY"))
	room, err := client.Rooms.Join(*roomName, *userName, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Join the room at %q", room.Data.Links.GuestJoin)

	setLogo(room)
}

func setLogo(room *eyeson.UserService) {
	fmt.Println("Wait until the room is ready")
	if err := room.WaitReady(); err != nil {
		fmt.Printf("Cannot determine room ready status: %v", err)
	}
	fmt.Println("Room is ready, set Layer")
	if err := room.SetLayer(OverlayUrl, eyeson.Foreground); err != nil {
		fmt.Printf("Failed to set overlay url: %v", err)
	}
	fmt.Println("Send welcome message")
	if err := room.Chat("Welcome aboard!"); err != nil {
		fmt.Printf("Failed to send chat message: %v", err)
	}
}
