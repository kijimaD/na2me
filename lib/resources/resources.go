package resources

import (
	"image/color"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var Master uiResources

func init() {
	fonts := newFonts()
	res := uiResources{
		Fonts:       fonts,
		ProgressBar: newProgressBarResources(),
		Button:      newButtonResources(fonts.UIFace),
	}
	Master = res
}

const (
	ScreenWidth  = 720
	ScreenHeight = 720
)

var (
	WhiteColor      = color.NRGBA{220, 220, 220, 255}
	WhiteTransColor = color.NRGBA{220, 220, 220, 140}
	BlackColor      = color.NRGBA{40, 40, 40, 255}
	GrayColor       = color.NRGBA{100, 100, 100, 255}

	TextPrimaryColor      = WhiteColor
	TextSecondaryColor    = GrayColor
	progressBarTrackColor = color.NRGBA{40, 40, 40, 255}
	progressBarFillColor  = TextSecondaryColor
	buttonIdleColor       = TextPrimaryColor
	buttonDisabledColor   = TextSecondaryColor
)

type uiResources struct {
	Fonts       *fonts
	ProgressBar *progressBarResources
	Button      *buttonResources
}

type progressBarResources struct {
	TrackImage *widget.ProgressBarImage
	FillImage  *widget.ProgressBarImage
}

type buttonResources struct {
	Image     *widget.ButtonImage
	TextColor *widget.ButtonTextColor
	Face      text.Face
	Padding   widget.Insets
}

func newProgressBarResources() *progressBarResources {
	track := image.NewNineSliceColor(progressBarTrackColor)
	fill := image.NewNineSliceColor(progressBarFillColor)

	return &progressBarResources{
		TrackImage: &widget.ProgressBarImage{
			Idle:     track,
			Hover:    track,
			Disabled: track,
		},
		FillImage: &widget.ProgressBarImage{
			Idle:     fill,
			Hover:    fill,
			Disabled: fill,
		},
	}
}

func newButtonResources(face text.Face) *buttonResources {
	idle := image.NewNineSliceColor(color.NRGBA{R: 70, G: 30, B: 10, A: 255})
	hover := image.NewNineSliceColor(color.NRGBA{R: 110, G: 180, B: 130, A: 255})
	pressed := image.NewNineSliceColor(color.NRGBA{R: 80, G: 110, B: 100, A: 255})

	return &buttonResources{
		Image: &widget.ButtonImage{
			Idle:    idle,
			Hover:   hover,
			Pressed: pressed,
		},
		TextColor: &widget.ButtonTextColor{
			Idle:     buttonIdleColor,
			Disabled: buttonDisabledColor,
		},
		Face: face,
		Padding: widget.Insets{
			Left:   20,
			Right:  20,
			Top:    8,
			Bottom: 8,
		},
	}
}
