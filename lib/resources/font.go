package resources

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/kijimaD/na2me/embeds"
)

type fonts struct {
	BodyFace     text.Face
	UIFace       text.Face
	SmallFace    text.Face
	TitleFace    text.Face
	BigTitleFace text.Face
}

func newFonts() *fonts {
	return &fonts{
		BodyFace: loadFont("ui/JF-Dot-Kappa20B.ttf", 26),
		UIFace:   loadFont("ui/JF-Dot-Kappa20B.ttf", 22),
	}
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
