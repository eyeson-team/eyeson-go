package main

import (
	"os"

	"github.com/eyeson-team/eyeson-go"
)

func main() {
	userService, _ := eyeson.NewUserServiceFromAccessKey(os.Getenv("ACCESS_KEY"))
	userService.SetLayer("https://docs.eyeson.com/img/examples/bg_p1.png",
		eyeson.Background)
}
