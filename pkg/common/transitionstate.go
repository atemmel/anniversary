package common

import (
	"github.com/hajimehoshi/ebiten"
	"image/color"
)

type TransitionState struct {
	Ticks int
	nTransitionTicks int

	fromState GameState
	toState GameState
	magnitude int
	fadeFrom *ebiten.Image
	currentFade *ebiten.Image
}

func NewTransitionState(src *ebiten.Image, fromState, toState GameState, ticks int) *TransitionState {
	fade, _ := ebiten.NewImage(src.Bounds().Max.X, src.Bounds().Max.Y, ebiten.FilterDefault)
	fade.Fill(color.RGBA{0, 0, 0, 0})
	return &TransitionState{
		0,
		ticks,
		fromState,
		toState,
		1,
		src,
		fade,
	}
}

func (t *TransitionState) GetInputs(g *Game) error {
	return nil
}

func (t *TransitionState) Update(g *Game) error {
	t.Ticks += t.magnitude
	if t.Ticks > t.nTransitionTicks {
		t.toState.ChangeTo(g)
		t.toState.Update(g)
		t.toState.Draw(g, t.fadeFrom)
		t.magnitude = -1;
		return nil
	} else if t.Ticks == 0 {
		g.As = &g.Ows
		return nil
	}
	scale := float64(t.Ticks) / float64(t.nTransitionTicks)
	t.currentFade.Fill(color.RGBA{0, 0, 0, uint8(255.0 * scale)})
	return nil
}

func (t *TransitionState) Draw(g *Game, screen *ebiten.Image) {
	screen.DrawImage(t.fadeFrom, &ebiten.DrawImageOptions{})
	screen.DrawImage(t.currentFade, &ebiten.DrawImageOptions{})
}

func (t *TransitionState) ChangeFrom(g *Game) {
}

func (t *TransitionState) ChangeTo(g *Game) {
}
