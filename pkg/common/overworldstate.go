package common

import (
	"errors"
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type GameState interface {
	GetInputs(g *Game) error
	Update(g *Game) error
	Draw(g *Game, screen *ebiten.Image)
}

type OverworldState struct {
	PlayerTextures []*ebiten.Image
	PlayerNameTags map[int]*ebiten.Image
	tileMap TileMap
	tileset *ebiten.Image
}

func gamepadUp() bool {
	return ebiten.GamepadAxis(0, 1) < -0.1 || ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton11)
}

func gamepadDown() bool {
	return ebiten.GamepadAxis(0, 1) > 0.1 || ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton13)
}

func gamepadLeft() bool {
	return ebiten.GamepadAxis(0, 0) < -0.1 || ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton14)
}

func gamepadRight() bool {
	return ebiten.GamepadAxis(0, 0) > 0.1 || ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton12)
}

func movingUp() bool {
	return ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyK) || ebiten.IsKeyPressed(ebiten.KeyW) || gamepadUp()
}

func movingDown() bool {
	return ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyJ) || ebiten.IsKeyPressed(ebiten.KeyS) || gamepadDown()
}

func movingLeft() bool {
	return ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyH) || ebiten.IsKeyPressed(ebiten.KeyA) || gamepadLeft()
}

func movingRight() bool {
	return ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyL) || ebiten.IsKeyPressed(ebiten.KeyD) || gamepadRight()
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

func (g *Game) CenterRendererOnPlayer() {
	g.Rend.LookAt(
		g.Player.Gx - WindowWidth / 2 + TileSize / 2,
		g.Player.Gy - WindowHeight / 2 + TileSize / 2,
	)
}
