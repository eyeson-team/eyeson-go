package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	eyeson "github.com/eyeson-team/eyeson-go"
)

const OverlayUrl string = "https://eyeson-team.github.io/api/images/eyeson-overlay.png"

func main() {
	url := os.Args[len(os.Args)-1]
	fmt.Println("Setup Webhook Listener for Endpoint", url)

	client := eyeson.NewClient(os.Getenv("API_KEY"))
	err := client.Webhook.Register(url, eyeson.WEBHOOK_ROOM)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Webhook.Unregister()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var data eyeson.Webhook
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			log.Println("Could not parse request,", err)
		}
		log.Println("Received new webhook for Room", data.Room.Name)
		if err := logRoomUpdate(&data); err != nil {
			log.Println("Could not store data,", err)
		}
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8042"
	}
	log.Println("Listen for connections on port", port)
	http.ListenAndServe(":"+port, nil)
}

func logRoomUpdate(data *eyeson.Webhook) error {
	filename := "./examples/logs/" + data.Room.Id + ".jsonl"
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = file.WriteString(string(jsonData))
	return err
}
