package common

import (
	"encoding/json"
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"io/ioutil"
	"image"
)

type Game struct {
	As GameState
	Ows *OverworldState
	Is *IntroState
	Gss *GroupedSpinnerState
	Sel *SelectionState
	Player Player
	Client Client
	Rend Renderer
	Audio Audio
}

func CreateGame() *Game {
	g := &Game{}
	g.Player.Id = 1
	g.Audio = NewAudio()
	g.Is = NewIntroState()
	g.Ows = NewOverworldState(g, g.Player.TexId)
	g.Sel = NewSelectionState()
	//g.Gss = NewGroupedSpinnerState()
	g.ChangeState(g.Is)
	//g.ChangeState(g.Gss)
	g.Rend = NewRenderer(WindowWidth, WindowHeight)
	return g
}

func (g *Game) TileIsOccupied(x int, y int, z int) bool {
	if x < 0 || x >= g.Ows.tileMap.Width || y < 0 ||  y >= g.Ows.tileMap.Height {
		return true
	}

	index := y * g.Ows.tileMap.Width + x

	// Out of bounds check
	if z < 0 || z >= len(g.Ows.tileMap.Tiles) {
		return true
	}

	if index >= len(g.Ows.tileMap.Tiles[z]) || index < 0 {
		return true
	}

	if g.Ows.tileMap.Collision[z][index] {
		return true
	}

	for _, p := range g.Client.playerMap.players {
		if p.X == x && p.Y == y {
			return true
		}
	}

	return false
}

func (g *Game) Update(screen *ebiten.Image) error {
	err := g.As.GetInputs(g)
	if err != nil {
		return err
	}
	err = g.As.Update(g)
	if err != nil {
		return err
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.As.Draw(g, screen)
}

func (g *Game) Load(str string, entrypoint int) {
	err := g.Ows.tileMap.OpenFile(str)
	if err != nil {
		panic(err)
	}
	g.Player.Location = str
	index := g.Ows.tileMap.GetEntryWithId(entrypoint)
	g.Player.X = g.Ows.tileMap.Entries[index].X
	g.Player.Y = g.Ows.tileMap.Entries[index].Y
	g.Player.Gx = float64(g.Player.X * TileSize)
	g.Player.Gy = float64(g.Player.Y * TileSize)
}

func (g *Game) Save() {
	bytes, err := json.Marshal(g.Ows.tileMap)
	if err != nil {
		fmt.Println(err)
	}
	ioutil.WriteFile(g.Player.Location, bytes, 0644)
}

func (g *Game) DrawPlayer(player *Player) {
	playerOpt := &ebiten.DrawImageOptions{}
	playerOpt.GeoM.Scale(2,2)

	if player.Dir == Left || player.Dir == Down {
		playerOpt.GeoM.Scale(-1, 1)
		playerOpt.GeoM.Translate(64, 0)
	}

	x := player.Gx + playerOffsetX
	y := player.Gy + playerOffsetY

	playerRect := image.Rect(
		player.Tx,
		player.Ty,
		player.Tx + TileSize,
		player.Ty + TileSize,
	)

	g.Rend.Draw(&RenderTarget{
		playerOpt,
		PlayerImgs[player.TexId],
		&playerRect,
		x,
		y,
		3,
	})

	tagOpt := &ebiten.DrawImageOptions{}

	tag := g.Ows.PlayerNameTags[g.Player.TexId]
	g.Rend.Draw(&RenderTarget{
		tagOpt,
		tag,
		nil,
		player.Gx + TileSize / 2 - float64(tag.Bounds().Dx() / 2),
		y - 16,
		69,
	})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WindowWidth, WindowHeight
}

func (g *Game) ChangeState(state GameState) {
	if g.As != nil {
		g.As.ChangeFrom(g)
	}
	g.As = state
	g.As.ChangeTo(g)
}

func (g *Game) PlayAudio() {
	g.Audio.audioPlayer.Play()
}
