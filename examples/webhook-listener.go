package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	eyeson "github.com/eyeson-team/eyeson-go"
)

var port = flag.Int("port", 8042, "listener HTTP port")

func main() {
	flag.Parse()

	url := os.Args[len(os.Args)-1]

	client := eyeson.NewClient(os.Getenv("API_KEY"))
	fmt.Println("Register webhook for endpoint", url)
	err := client.Webhook.Register(url, eyeson.WEBHOOK_ROOM)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var data eyeson.Webhook
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			log.Println("Could not parse request: ", err)
		}
		log.Println("Received new webhook for Room", data.Room.Name)
		if err := logRoomUpdate(&data); err != nil {
			log.Println("Could not store data,", err)
		}
		w.WriteHeader(204)
	})

	srv := &http.Server{Addr: fmt.Sprintf(":%d", *port), Handler: mux}
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Printf("Listen for connections on port %d", *port)
		if err = srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()
	<-stop

	log.Println("Shutting down...")
	log.Println("Unregister webhook")
	if err = client.Webhook.Unregister(); err != nil {
		log.Fatal("Failed to unregister webhook: ", err)
	}
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
