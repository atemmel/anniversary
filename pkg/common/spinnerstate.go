package common

import(
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	"image"
	"image/color"
)

var(
	SpinnerFont font.Face = nil
)
const spinnerFontSize = 14

type SpinnerState struct {
	spinnerState *ebiten.Image
	roll spinnerRoll
}

type spinnerRoll struct {
	y, dy, ddy, dddy float64
	y2 float64
	img *ebiten.Image
	clipped *ebiten.Image
}

func NewSpinnerState() *SpinnerState {
	s := &SpinnerState{}

	s.roll.dy = 10
	s.roll.ddy = -0.0001
	s.roll.dddy = -0.0001
	SpinnerFont = loadFont(ResourceDir + "Mario-Kart-DS.ttf", spinnerFontSize)

	strs := make([]string, 3)
	strs[0] = "DUELL"
	strs[1] = "RED VS BLU"
	strs[2] = "FREE FOR ALL"
	img := buildSpinnerTexture(strs)
	s.roll.img = img
	w, h := img.Size()
	s.roll.y2 = -float64(h)
	s.roll.clipped, _ = ebiten.NewImage(w, h / len(strs) * 3, ebiten.FilterDefault)

	return s
}

func buildSpinnerTexture(strs []string) *ebiten.Image {
	bounds := make([]image.Rectangle, len(strs))

	for i := range bounds {
		bounds[i] = text.BoundString(SpinnerFont, strs[i])
	}

	const padimgw = 300
	const padimgh = 2
	const vspace = 10
	padimg, _ := ebiten.NewImage(padimgw, padimgh, ebiten.FilterDefault)
	padimg.Fill(color.Black)

	imgh := -padimgh * 2
	for i := range bounds {
		imgh += spinnerFontSize
		imgh += bounds[i].Dy()
		imgh += padimgh
		imgh += vspace * 2
	}

	imgh -= padimgh * 2 + vspace * 2 + spinnerFontSize

	img, _ := ebiten.NewImage(padimgw, imgh, ebiten.FilterDefault)
	img.Fill(color.White)

	y := float64(-padimgh * 2)
	for i := range bounds {
		x := padimgw / 2 - bounds[i].Dx() / 2
		text.Draw(img, strs[i], SpinnerFont, x, int(y) + spinnerFontSize + vspace, color.RGBA{255,0,0,255})
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(0, y + spinnerFontSize + padimgh + vspace * 2)
		img.DrawImage(padimg, opt)
		y += padimgh + float64(bounds[i].Dy() + vspace * 2)
	}
	df := &ebiten.DrawImageOptions{}
	df.GeoM.Translate(0, y + padimgh)
	img.DrawImage(padimg, df)
	return img
}

func (s *SpinnerState) Draw(g *Game, screen *ebiten.Image) {
	screen.Fill(color.RGBA{255,255,0,255})
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(0, s.roll.y)
	opt2 := &ebiten.DrawImageOptions{}
	opt2.GeoM.Translate(0, s.roll.y2)
	s.roll.clipped.DrawImage(s.roll.img, opt)
	s.roll.clipped.DrawImage(s.roll.img, opt2)
	screen.DrawImage(s.roll.clipped, &ebiten.DrawImageOptions{})
}

func (s *SpinnerState) Update(g *Game) error {
	s.roll.ddy += s.roll.dddy
	s.roll.dy += s.roll.ddy
	if s.roll.dy < 0.1 {
		s.roll.ddy = 0
		s.roll.dy = 0
	}

	s.roll.y += s.roll.dy
	s.roll.y2 += s.roll.dy

	_, h := s.roll.img.Size()
	if s.roll.y > float64(h) {
		s.roll.y -= float64(h) * 2
	}
	if s.roll.y2 > float64(h) {
		s.roll.y2 -= float64(h) * 2
	}
	return nil
}

func (s *SpinnerState) GetInputs(g *Game) error {
	return nil
}

func (s *SpinnerState) ChangeFrom(g *Game) {
}

func (s *SpinnerState) ChangeTo(g *Game) {
}
