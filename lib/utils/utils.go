package utils

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/na2me/embeds"
)

func GetPtr[T any](x T) *T {
	return &x
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
