package main

import (
	"fmt"
	"os"

	"github.com/eyeson-team/eyeson-go"
)

func main() {
	userService, _ := eyeson.NewUserServiceFromAccessKey(os.Getenv("ACCESS_KEY"))
	layoutName := "custom-map"
	err := userService.SetLayout(eyeson.AUTO, []string{""}, false, false, &layoutName,
		&eyeson.LayoutMap{
			Positions: []eyeson.LayoutPos{
				{X: 10, Y: 10, Width: 100, Height: 200,
					ObjectFit: eyeson.CONTAIN}}},
		nil)
	if err != nil {
		fmt.Printf("Failed: %s", err)
	}
}
