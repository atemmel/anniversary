package common

import(
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"io/ioutil"
	"image"
	"image/color"
)

var NameTagFont font.Face = nil

func NewNameTag(str string) *ebiten.Image {
	const nameTagSize = 12
	if NameTagFont == nil {
		NameTagFont = loadFont(ResourceDir + "MonsterFriendFore.otf", nameTagSize)
	}

	rect := text.BoundString(NameTagFont, str)
	rect = image.Rect(0, 0, rect.Bounds().Dx() + 12, rect.Bounds().Dy() + 8)
	img, _ := ebiten.NewImage(rect.Bounds().Dx(), rect.Bounds().Dy(), ebiten.FilterDefault)
	img.Fill(color.RGBA{0,0,0,105})
	text.Draw(img, str, NameTagFont, 0 + 6, nameTagSize + 3, color.White)
	return img
}

func loadFont(path string, size float64) font.Face {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic( err)
	}

	tt, err := opentype.Parse(bytes)
	if err != nil {
		panic( err)
	}

	const dpi = 72
	font, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	return font
}
