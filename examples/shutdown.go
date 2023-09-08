package main

import (
	"flag"
	"fmt"
	"os"

	eyeson "github.com/eyeson-team/eyeson-go"
)

func main() {
	id := flag.String("id", "", "meeting identifier to shutdown")
	flag.Parse()

	client, err := eyeson.NewClient(os.Getenv("API_KEY"))
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	err = client.Rooms.Shutdown(*id)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
