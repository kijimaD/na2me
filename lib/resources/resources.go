package resources

import (
	"image/color"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
)

var Master uiResources

func init() {
	Master.Fonts = newFonts()
	Master.ProgressBar = newProgressBarResources()
}

const (
	ScreenWidth  = 720
	ScreenHeight = 720
)

var (
	TextIdleColor         = "dff4ffFF"
	TextDisabledColor     = "646464FF"
	progressBarTrackColor = color.NRGBA{40, 40, 40, 255}
	progressBarFillColor  = color.NRGBA{100, 100, 100, 255}
)

type uiResources struct {
	Fonts       *fonts
	ProgressBar *progressBarResources
}

type progressBarResources struct {
	TrackImage *widget.ProgressBarImage
	FillImage  *widget.ProgressBarImage
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
