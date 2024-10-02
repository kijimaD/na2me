package states

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	embeds "github.com/kijimaD/na2me/embeds"
)

func loadImage(filename string) (*ebiten.Image, error) {
	bs, err := embeds.FS.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	dec, _, err := image.Decode(bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	img := ebiten.NewImageFromImage(dec)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func loadFont(path string, size float64) text.Face {
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
