package main

import (
	"bytes"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

func loadImage(filename string) (*ebiten.Image, error) {
	bs, err := FS.ReadFile(filename)
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
