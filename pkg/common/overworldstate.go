package common

import (
	"errors"
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
)

type GameState interface {
	GetInputs(g *Game) error
	Update(g *Game) error
	Draw(g *Game, screen *ebiten.Image)
	ChangeTo(g *Game)
	ChangeFrom(g *Game)
}

type OverworldState struct {
	PlayerNameTags map[int]*ebiten.Image
	tileMap TileMap
	tileset *ebiten.Image
}

func NewOverworldState(playerId int) *OverworldState {
	ows := &OverworldState{}
	var err error
	loadPlayerImgs()
	ows.PlayerNameTags = make(map[int]*ebiten.Image)
	ows.PlayerNameTags[playerId] = NewNameTag(NameIndexMap[playerId])

	ows.tileset, _, err = ebitenutil.NewImageFromFile(ResourceDir + "textures/tileset1.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

  return ows
}

func (ows *OverworldState) SetPlayerTag(id int) {
	ows.PlayerNameTags[id] = NewNameTag(NameIndexMap[id])
}

func holdingSprint() bool {
	return ebiten.IsKeyPressed(ebiten.KeyShift) || ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton1)
}

func (o *OverworldState) GetInputs(g *Game) error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errors.New("")	//TODO Gotta be a better way to do this
	}

	if !g.Player.isWalking && holdingSprint() {
		g.Player.isRunning = true
	} else {
		g.Player.isRunning = false
	}

	if movingUp() {
		g.Player.TryStep(Up, g)
	} else if movingDown() {
		g.Player.TryStep(Down, g)
	} else if movingRight() {
		g.Player.TryStep(Right, g)
	} else if movingLeft() {
		g.Player.TryStep(Left, g)
	} else {
		g.Player.TryStep(Static, g)
	}

	return nil
}

func (o *OverworldState) Update(g *Game) error {
	g.Player.Update(g)

	/*
	if g.Client.Active {
		g.Client.WritePlayer(&g.Player)
	}
	*/

	return nil
}

func (o *OverworldState) Draw(g *Game, screen *ebiten.Image) {
	o.tileMap.Draw(&g.Rend, o.tileset)
	g.DrawPlayer(&g.Player)

	/*
	if g.Client.Active {
		g.Client.playerMap.mutex.Lock()
		for _, player := range g.Client.playerMap.players {
			if player.Location == g.Player.Location {
				g.DrawPlayer(&player)
			}
		}
		g.Client.playerMap.mutex.Unlock()
	}
	*/

	g.CenterRendererOnPlayer()
	g.Rend.Display(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf(
`player.frames: %d`,
g.Player.frames))
}

func (o *OverworldState) ChangeTo(g *Game) {
	g.Audio.audioPlayer.Play()
}

func (o *OverworldState) ChangeFrom(g *Game) {
	g.Audio.audioPlayer.Rewind()
	g.Audio.audioPlayer.Pause()
}

func (g *Game) CenterRendererOnPlayer() {
	g.Rend.LookAt(
		g.Player.Gx - WindowWidth / 2 + TileSize / 2,
		g.Player.Gy - WindowHeight / 2 + TileSize / 2,
	)
}
