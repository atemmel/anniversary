package common

import (
	"github.com/hajimehoshi/ebiten"
)

type GroupedSpinnerStateConfig struct {
	strings [][]string
	winners []int
}

type GroupedSpinnerState struct {
	states []*SpinnerState
	currentIndex int
}

func NewGroupedSpinnerState() *GroupedSpinnerState {
	gss := &GroupedSpinnerState{}

	gss.states = make([]*SpinnerState, 1)
	strs := make([]string, 3)
	strs[0] = "DUELL"
	strs[1] = "RED VS BLU"
	strs[2] = "FREE FOR ALL"
	winner := 0

	for i := range gss.states {
		gss.states[i] = NewSpinnerState(strs, winner)
	}

	return gss
}

func (gss *GroupedSpinnerState) Draw(g *Game, screen *ebiten.Image) {
	gss.states[gss.currentIndex].Draw(g, screen)
}

func (gss *GroupedSpinnerState) Update(g *Game) error {
	return gss.states[gss.currentIndex].Update(g)
}

func (gss *GroupedSpinnerState) GetInputs(g *Game) error {
	return gss.states[gss.currentIndex].GetInputs(g)
}

func (gss *GroupedSpinnerState) ChangeFrom(g *Game) {
}

func (gss *GroupedSpinnerState) ChangeTo(g *Game) {
}
