package main

import (
	"io/ioutil"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	"image/color"
	"log"
	"strings"
)

const (
	WindowWidth = 1280
	WindowHeight = 720

	ChatWidth = WindowWidth
	ChatHeight = 300
)

var showChat = false
var chatBackgroundColor = color.RGBA{0, 0, 0, 205}
var chatImage *ebiten.Image
var chatFont font.Face
var messages = make([]string, 0)

type Game struct {}

func (g *Game) Update(screen *ebiten.Image) error {
	if showChat == false && inpututil.IsKeyJustPressed(ebiten.KeyT) {
		showChat = true
	}

	if showChat == true && inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		showChat = false
	}

	return nil
}

func (g* Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White);

	if showChat {
		DrawChat(chatImage, screen)
	}
}

func DrawChat(chat, screen *ebiten.Image) {
	opt := &ebiten.DrawImageOptions{}
	joined := strings.Join(messages, "\n")
	rect := text.BoundString(chatFont, joined)
	rect.Add(image.Point{0, WindowHeight - rect.Bounds().Dy() * 2})
	opt.GeoM.Translate(0, float64(rect.Max.Y))
	subimg, _ := ebiten.NewImageFromImage(screen.SubImage(rect), ebiten.FilterDefault)
	subimg.Fill(chatBackgroundColor)
	text.Draw(subimg, joined, chatFont, 0, 0, color.White)
	screen.DrawImage(subimg, opt)
	/*
	opt.GeoM.Translate(0, WindowHeight - float64(chat.Bounds().Dy()))
	chat.Fill(chatBackgroundColor)
	text.Draw(chat, strings.Join(messages, "\n"), chatFont, 10, 50, color.White )
	screen.DrawImage(chat, opt)
	*/
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return WindowWidth, WindowHeight
}

func NewChatImage() *ebiten.Image {
	img, _ := ebiten.NewImage(ChatWidth, ChatHeight, ebiten.FilterDefault)
	img.Fill(color.RGBA{0,0,0,205})
	return img
}

func LoadFont(path string) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	tt, err := opentype.Parse(bytes)
	if err != nil {
		return err
	}

	const dpi = 72
	chatFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	return err
}

func main() {
	chatImage = NewChatImage()
	err := LoadFont("./MonsterFriendFore.otf")
	if err != nil {
		log.Fatal(err)
	}
	messages = append(messages, "SOMEONE: Tjo", "ASDF: Asdsd dfdsfds", "JSDFSDFSDF: QWE QWE QW EQWE ")
	ebiten.SetWindowSize(WindowWidth, WindowHeight)
	ebiten.SetWindowTitle("Anniversary")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
