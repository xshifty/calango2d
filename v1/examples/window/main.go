package main

import (
	"os"

	"github.com/xshifty/calango/v1/engine"
)

func main() {
	game := engine.New("Calango Engine - Example Window", 1920, 1080, 30)
	game.SetFrameRate(game.GetDesktopFrameRate())

	game.AddScene(engine.NewScene(
		"main",
		func(s *engine.Scene) error {
			return nil
		},
		func(s *engine.Scene) error {
			return nil
		},
	))

	os.Exit(game.Run())
}
