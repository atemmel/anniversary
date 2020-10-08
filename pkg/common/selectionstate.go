package common

import(
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image"
	"image/color"
)

const (
	gridOffsetY = 40
	handvel = 4
	gridW = 8
)

type character struct {
	x, y float64
}

type SelectionState struct {
	hand *ebiten.Image
	disc *ebiten.Image
	grid *ebiten.Image
	readytofight *ebiten.Image
	handanim handanim
	rects []image.Rectangle
}

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
	loadPlayerImgs()
	hand, _, err := ebitenutil.NewImageFromFile(ResourceDir + "textures/hand.png", ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}
	disc, _, err := ebitenutil.NewImageFromFile(ResourceDir + "textures/disc.png", ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}
	border, _, err := ebitenutil.NewImageFromFile(ResourceDir + "textures/border.png", ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}
	readytofight, _, err := ebitenutil.NewImageFromFile(ResourceDir + "textures/readytofight.png", ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}
	const factor = 0.5
	x, y := hand.Size()
	scaledhand, _ := ebiten.NewImage(int(float64(x) * factor), int(float64(y) * factor), ebiten.FilterDefault)
	x, y = disc.Size()
	scaleddisc, _ := ebiten.NewImage(int(float64(x) * factor), int(float64(y) * factor), ebiten.FilterDefault)
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Scale(factor, factor)
	scaledhand.DrawImage(hand, opt)
	scaleddisc.DrawImage(disc, opt)

	grid := buildGrid(len(PlayerImgs), border)

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
		scaledhand,
		scaleddisc,
		grid,
		readytofight,
		handanim,
		buildRects(len(PlayerImgs), border),
	}
}

func buildRects(m int, border *ebiten.Image) []image.Rectangle {
	rects := make([]image.Rectangle, 0)

	bw, by := border.Size()
	w := gridW * bw
	x := WindowWidth / 2 - w / 2
	y := gridOffsetY

	rows := 0
	rem := m
	for ; rem > gridW; rem -= gridW {
		rows++
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < gridW; j++ {
			r := image.Rect(j * bw, i * by, j * bw + bw, i * by + by)
			r = r.Add(image.Point{x, y})
			rects = append(rects, r)
		}
	}

	w = rem * bw
	x = WindowWidth / 2 - w / 2
	y = gridOffsetY + by * rows

	for j := 0; j < gridW; j++ {
		r := image.Rect(x + (j * bw), y, x + ((j + 1) * bw), y + by)
		rects = append(rects, r)
	}

	return rects
}

func buildGrid(m int, border *ebiten.Image) *ebiten.Image {
	n := 0
	rem := m
	for ; rem > gridW; rem -= gridW {
		n++
	}

	borderW, borderH := border.Size()
	img, _ := ebiten.NewImage(borderW * gridW, borderH * (n + 1), ebiten.FilterDefault)
	ix, iy := PlayerImgs[0].Size()
	r := image.Rect(0, 0, ix, iy / 3)

	index := 0
	for y := 0; y < n; y++ {
		for x := 0; x < gridW; x++ {
			popt := &ebiten.DrawImageOptions{}
			popt.GeoM.Scale(2, 2)
			popt.GeoM.Translate(float64(x * borderW + 1), float64(y * borderH + 1))
			img.DrawImage(PlayerImgs[index].SubImage(r).(*ebiten.Image), popt)
			bopt := &ebiten.DrawImageOptions{}
			bopt.GeoM.Translate(float64(x * borderW), float64(y * borderH))
			img.DrawImage(border, bopt)
			index++
		}
	}

	offsetx := (gridW * borderW) / 2 - (rem * borderW) / 2

	for x := 0; x < rem; x++ {
		popt := &ebiten.DrawImageOptions{}
		popt.GeoM.Scale(2, 2)
		popt.GeoM.Translate(float64(x * borderW + 1 + offsetx), float64(n * borderH + 1))
		img.DrawImage(PlayerImgs[index].SubImage(r).(*ebiten.Image), popt)
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(float64(x * borderW + offsetx), float64(n * borderH))
		img.DrawImage(border, opt)
		index++
	}

	return img
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
		p1 := image.Point{int(s.handanim.x), int(s.handanim.y)}
		for i, r := range s.rects {
			if p1.In(r) {
				fmt.Println(i)
				g.Player.Id = i
				g.Ows.SetPlayerTag(i)
				img, _ := ebiten.NewImage(WindowWidth, WindowHeight, ebiten.FilterDefault)
				s.Draw(g, img)
				g.ChangeState(NewTransitionState(img, s, g.Ows, 40))
			}
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
	drawGrid(s.grid, screen)
	hopt := ebiten.DrawImageOptions{}
	hopt.GeoM.Translate(s.handanim.x, s.handanim.y)
	dopt := ebiten.DrawImageOptions{}
	dopt.GeoM.Translate(s.handanim.x - 5, s.handanim.y - 9)
	screen.DrawImage(s.disc, &dopt)
	if s.handanim.selected {
		screen.DrawImage(s.readytofight, &ebiten.DrawImageOptions{})
	}
	rect := image.Rect(int(s.handanim.tx), int(s.handanim.ty), int(s.handanim.tx + s.handanim.w), int(s.handanim.ty + s.handanim.h))
	screen.DrawImage(s.hand.SubImage(rect).(*ebiten.Image), &hopt)
}

func drawGrid(g *ebiten.Image, screen *ebiten.Image) {
	opt := &ebiten.DrawImageOptions{}
	w, _ := g.Size()
	x := WindowWidth / 2 - w / 2
	y := gridOffsetY
	opt.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(g, opt)
}

func (s *SelectionState) ChangeFrom(g *Game) {
	g.Audio.characterselectPlayer.Pause()
	g.Audio.characterselectPlayer.Rewind()
}

func (s *SelectionState) ChangeTo(g *Game) {
	g.Audio.characterselectPlayer.Play()
}
