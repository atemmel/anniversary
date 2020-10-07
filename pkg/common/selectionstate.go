package common

import(
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image"
	"image/color"
)

type character struct {
	x, y float64
}

type SelectionState struct {
	images []*ebiten.Image
	hand *ebiten.Image
	disc *ebiten.Image
	border *ebiten.Image
	readytofight *ebiten.Image
	handanim handanim

}

const handvel = 4

type handanim struct {
	x, y, tx, ty, w, h float64
	selected bool
}

func (h *handanim) Move(x, y float64) {
	h.x += x
	h.y += y

	if h.x + h.w > WindowWidth {
		h.x = WindowWidth - h.w
	} else if h.x < 0 {
		h.x = 0
	}

	if h.y < 0 {
		h.y = 0
	} else if h.y + h.h > WindowHeight {
		h.y = WindowHeight - h.h
	}
}

func (h *handanim) Place() {
	//h.tx = 
}

func NewSelectionState() *SelectionState {
	hand, _, err := ebitenutil.NewImageFromFile("resources/textures/hand.png", ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}
	disc, _, err := ebitenutil.NewImageFromFile("resources/textures/disc.png", ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}
	border, _, err := ebitenutil.NewImageFromFile("resources/textures/border.png", ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}
	readytofight, _, err := ebitenutil.NewImageFromFile("resources/textures/readytofight.png", ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}
	x, y := hand.Size()
	const factor = 0.5
	scaledhand, _ := ebiten.NewImage(int(float64(x) * factor), int(float64(y) * factor), ebiten.FilterDefault)
	x, y = disc.Size()
	scaleddisc, _ := ebiten.NewImage(int(float64(x) * factor), int(float64(y) * factor), ebiten.FilterDefault)
	x, y = border.Size()
	scaledborder, _ := ebiten.NewImage(int(float64(x) * factor), int(float64(y) * factor), ebiten.FilterDefault)
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Scale(factor, factor)
	scaledhand.DrawImage(hand, opt)
	scaleddisc.DrawImage(disc, opt)
	scaledborder.DrawImage(border, opt)
	x, y = scaledhand.Size()
	handanim := handanim{
		50,
		50,
		float64(x) / 3 * 2,
		0,
		float64(x) / 3,
		float64(y),
		false,
	}
	return &SelectionState{
		make([]*ebiten.Image, 0),
		scaledhand,
		scaleddisc,
		border,
		readytofight,
		handanim,
	}
}

func (s *SelectionState) GetInputs(g *Game) error {
	if movingDown() {
		s.handanim.Move(0, handvel)
	}

	if movingUp() {
		s.handanim.Move(0, -handvel)
	}

	if movingLeft() {
		s.handanim.Move(-handvel, 0)
	}

	if movingRight() {
		s.handanim.Move(handvel, 0)
	}

	if accept() {
		b1 := image.Rect(0, 0, s.border.Bounds().Dx(), s.border.Bounds().Dy())
		p1 := image.Point{int(s.handanim.x), int(s.handanim.y)}
		if p1.In(b1) {
			img, _ := ebiten.NewImage(WindowWidth, WindowHeight, ebiten.FilterDefault)
			s.Draw(g, img)
			g.ChangeState(NewTransitionState(img, s, g.Ows, 40))
		}
		/*
		s.handanim.tx += s.handanim.w
		if s.handanim.tx == s.handanim.w * 3 {
			s.handanim.tx = 0
		}
		*/
	}

	return nil
}

func (s *SelectionState) Update(g *Game) error {
	return nil
}

func (s *SelectionState) Draw(g *Game, screen *ebiten.Image) {
	screen.Fill(color.Black)
	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Translate(s.handanim.x, s.handanim.y)
	rect := image.Rect(int(s.handanim.tx), int(s.handanim.ty), int(s.handanim.tx + s.handanim.w), int(s.handanim.ty + s.handanim.h))
	screen.DrawImage(s.border, &ebiten.DrawImageOptions{})
	screen.DrawImage(s.disc, &opt)
	if s.handanim.selected {
		screen.DrawImage(s.readytofight, &ebiten.DrawImageOptions{})
	}
	screen.DrawImage(s.hand.SubImage(rect).(*ebiten.Image), &opt)
}

func (s *SelectionState) ChangeFrom(g *Game) {
	g.Audio.characterselectPlayer.Pause()
	g.Audio.characterselectPlayer.Rewind()
}

func (s *SelectionState) ChangeTo(g *Game) {
	g.Audio.characterselectPlayer.Play()
}
