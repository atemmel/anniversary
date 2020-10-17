package main

import (
	"github.com/hajimehoshi/ebiten"
	"log"

	"github.com/atemmel/anniversary/pkg/common"
)

func main() {
	g := common.CreateGame()
	g.Load("./resources/tilemaps/event.json", 0)
	ebiten.SetWindowSize(common.WindowWidth, common.WindowHeight)
	ebiten.SetWindowTitle("Anniversary")
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetFullscreen(true)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	g.Client = common.CreateClient()
	g.Player.Id = g.Client.Connect()
	if g.Client.Active {
		g.Player.Connected = true
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
