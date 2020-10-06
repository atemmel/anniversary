package common

import(
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type viewport struct {
	x16 int
	y16 int
}

func (p *viewport) Move(is *IntroState) {
	w, h := is.tile.Size()
	maxX16 := w * 16
	maxY16 := h * 16

	p.x16 += w / 32
	p.y16 += h / 32
	p.x16 %= maxX16
	p.y16 %= maxY16
}

func (p *viewport) Position() (int, int) {
	return p.x16, p.y16
}

type IntroState struct {
	tile *ebiten.Image
	logo *ebiten.Image
	viewport viewport
}

func NewIntroState() *IntroState {
	tile, _, err := ebitenutil.NewImageFromFile("resources/textures/tile.png", ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}
	logo, _, err := ebitenutil.NewImageFromFile("resources/textures/logo.png", ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}
	return &IntroState{
		tile,
		logo,
		viewport{},
	}
}

func (is *IntroState) GetInputs(g *Game) error {
	return nil
}

func (is *IntroState) Update(g *Game) error {
	is.viewport.Move(is)
	if accept() {
		img, _ := ebiten.NewImage(WindowWidth, WindowHeight, ebiten.FilterDefault)
		is.Draw(g, img)
		g.ChangeState(NewTransitionState(img, is, &g.Ows, 40))
		//g.ChangeState(&g.Ows)
	}
	return nil
}

func (is *IntroState) Draw(g *Game, screen *ebiten.Image) {
	x16, y16 := is.viewport.Position()
	offsetX, offsetY := float64(-x16)/16, float64(-y16)/16

	// Draw bgImage on the screen repeatedly.
	const repeat = 6
	w, h := is.tile.Size()
	for j := 0; j < repeat; j++ {
		for i := 0; i < repeat; i++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(w*i), float64(h*j))
			op.GeoM.Translate(offsetX, offsetY)
			screen.DrawImage(is.tile, op)
		}
	}

	op := &ebiten.DrawImageOptions{}
	w, h = is.logo.Size()
	x := w / 4 + WindowWidth / 4
	y := h / 4 + WindowHeight / 4
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(is.logo, op)
}

func (is *IntroState) ChangeFrom(g *Game) {
	g.Audio.titlePlayer.Pause()
	g.Audio.titlePlayer.Rewind()
}

func (is *IntroState) ChangeTo(g *Game) {
	g.Audio.titlePlayer.Play()
}
