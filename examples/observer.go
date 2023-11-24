package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	eyeson "github.com/eyeson-team/eyeson-go"
)

func main() {
	userName := flag.String("user", "observer", "Meeting Assistant")
	roomName := flag.String("room", "demo", "room identifier")
	apiEndpoint := flag.String("api-ep", "", "Optional api endpoint")
	flag.Parse()

	options := []eyeson.ClientOption{}
	if len(*apiEndpoint) > 0 {
		options = append(options, eyeson.WithCustomEndpoint(*apiEndpoint))
	}

	client, err := eyeson.NewClient(os.Getenv("API_KEY"), options...)
	if err != nil {
		log.Fatal(err)
	}
	room, err := client.Rooms.Join(*roomName, *userName, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Join the room at %q", room.Data.Links.GuestJoin)

	msgCh, err := client.Observer.Connect(context.Background(), room.Data.Room.ID)
	if err != nil {
		fmt.Println("Failed to connect with observer: %s", err)
	}
	for {
		select {
		case msg, ok := <-msgCh:
			if !ok {
				fmt.Println("Channel closed. Probably disconnected")
				return
			}
			fmt.Println("Received event type: ", msg.GetType())
			switch m := msg.(type) {
			case *eyeson.ParticipantUpdate:
				fmt.Printf("user %s is online %v\n", m.Participant.Name, m.Participant.Online)
			case *eyeson.RoomUpdate:
				fmt.Printf("Room %s is ready %v\n", m.Content.Name, m.Content.Ready)
				if m.Content.Shutdown {
					fmt.Printf("Room %s is shutting down now.\n", m.Content.Name)
					return
				}
			case *eyeson.Chat:
				fmt.Printf("Chat: %s - %s", m.ClientID, m.Content)
			}
		}

	}
}
