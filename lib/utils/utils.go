package utils

import (
	"bytes"
	"image"
	"image/color"
	"log"

	uimage "github.com/ebitenui/ebitenui/image"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/kijimaD/na2me/embeds"
)

var (
	BodyFont text.Face
)

func init() {
	BodyFont = LoadFont("ui/JF-Dot-Kappa20B.ttf", 26)
}

func GetPtr[T any](x T) *T {
	return &x
}

func LoadButtonImage() *widget.ButtonImage {
	idle := uimage.NewNineSliceColor(color.RGBA{R: 110, G: 110, B: 180, A: 255})
	hover := uimage.NewNineSliceColor(color.RGBA{R: 110, G: 180, B: 130, A: 255})
	pressed := uimage.NewNineSliceColor(color.RGBA{R: 80, G: 110, B: 100, A: 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}
}

func LoadImage(filename string) *ebiten.Image {
	bs, err := embeds.FS.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	dec, _, err := image.Decode(bytes.NewReader(bs))
	if err != nil {
		log.Fatal(err)
	}
	img := ebiten.NewImageFromImage(dec)
	if err != nil {
		log.Fatal(err)
	}

	return img
}

func LoadFont(path string, size float64) text.Face {
	fontFile, err := embeds.FS.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	s, err := text.NewGoTextFaceSource(fontFile)
	if err != nil {
		log.Fatal(err)
	}

	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}
}
