package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	eyeson "github.com/eyeson-team/eyeson-go"
)

// This example demonstrates the observer capabilities. It starts
// a room and connects on the observer socket.
// If a new participant gets online, a giphy is played and a
// welcome-chat message is sent.
func main() {
	userName := flag.String("user", "observer", "Meeting Assistant")
	roomName := flag.String("room", "demo", "room identifier")
	apiEndpoint := flag.String("api-ep", "", "Optional api endpoint")
	flag.Parse()

	options := []eyeson.ClientOption{}
	if len(*apiEndpoint) > 0 {
		options = append(options, eyeson.WithCustomEndpoint(*apiEndpoint))
	}

	if len(os.Getenv("API_KEY")) == 0 {
		fmt.Println("Error: Please set the environment variable API_KEY")
		return
	}

	client, err := eyeson.NewClient(os.Getenv("API_KEY"), options...)
	if err != nil {
		log.Fatal(err)
	}
	room, err := client.Rooms.Join(*roomName, *userName, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nJoin the room at %q\n\n", room.Data.Links.GuestJoin)

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

				if m.Participant.Online {
					// play a hello after 1s for this new participant
					go func() {
						time.Sleep(1 * time.Second)
						helloVideo := "https://media4.giphy.com/media/3pZipqyo1sqHDfJGtz/giphy.mp4"
						if err = room.StartPlayback(helloVideo, ""); err != nil {
							fmt.Println("Failed to start playback: ", err)
						}
					}()

					// add a chat message
					if err = room.Chat(fmt.Sprintf("Please welcome user %s to this meeting", m.Participant.Name)); err != nil {
						fmt.Println("Failed to send chat: ", err)
					}
				}

			case *eyeson.RoomUpdate:
				fmt.Printf("Room %s is ready %v\n", m.Content.Name, m.Content.Ready)
				if m.Content.Shutdown {
					fmt.Printf("Room %s is shutting down now.\n", m.Content.Name)
					return
				}
			case *eyeson.Chat:
				fmt.Printf("Chat: %s - %s\n", m.ClientID, m.Content)
			}
		}

	}
}
