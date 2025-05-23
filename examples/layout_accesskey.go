package main

import (
	"fmt"
	"os"

	"github.com/eyeson-team/eyeson-go"
)

func main() {
	userService, _ := eyeson.NewUserServiceFromAccessKey(os.Getenv("ACCESS_KEY"))
	layoutName := "custom-map"
	audioInsertConfig := eyeson.AudioInsert{
		Config: eyeson.Enabled,
		Position: &eyeson.AudioInsertPosition{
			X: 10,
			Y: 10,
		},
	}
	err := userService.SetLayout(eyeson.Auto,
		&eyeson.SetLayoutOptions{
			Users:           []string{""},
			LayoutName:      layoutName,
			VoiceActivation: true,
			LayoutMap: &eyeson.LayoutMap{
				Positions: []eyeson.LayoutPos{
					{X: 10, Y: 10, Width: 100, Height: 200,
						ObjectFit: eyeson.Contain}}},
			AudioInsert: &audioInsertConfig,
		})

	if err != nil {
		fmt.Printf("Failed: %s", err)
	}
}
