package common

import(
	"github.com/hajimehoshi/ebiten"
)

type SelectionState struct {
	images []*ebiten.Image
}

func (s *SelectionState) GetInputs(g *Game) error {
	return nil
}

func (s *SelectionState) Update(g *Game) error {
	return nil
}

func (s *SelectionState) Draw(g *Game, screen *ebiten.Image) {

}

func (s *SelectionState) ChangeTo(g *Game) {

}

func (s *SelectionState) ChangeFrom(g *Game) {

}
