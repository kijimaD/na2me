package resources

import (
	"image/color"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/kijimaD/na2me/lib/utils"
)

var Master uiResources

func init() {
	fonts := newFonts()
	res := uiResources{
		Fonts:       fonts,
		ProgressBar: newProgressBarResources(),
		Button:      newButtonResources(fonts.UIFace),
		List:        newListResources(),
		Backgrounds: newBackgroundResources(),
	}
	Master = res
}

const (
	ScreenWidth  = 720
	ScreenHeight = 720
)

var (
	whiteColor      = color.NRGBA{240, 240, 240, 255}
	lightColor      = color.NRGBA{220, 220, 220, 255}
	lightTransColor = color.NRGBA{220, 220, 220, 140}
	blackColor      = color.NRGBA{40, 40, 40, 255}
	grayColor       = color.NRGBA{100, 100, 100, 255}

	TextPrimaryColor    = lightColor
	TextSecondaryColor  = grayColor
	TextBodyColor       = whiteColor
	BGPrimaryTransColor = lightTransColor

	progressBarTrackColor = color.NRGBA{40, 40, 40, 255}
	progressBarFillColor  = TextSecondaryColor
	buttonIdleColor       = TextPrimaryColor
	buttonDisabledColor   = TextSecondaryColor
	listIdleColor         = lightColor
)

type uiResources struct {
	Fonts       *fonts
	ProgressBar *progressBarResources
	Button      *buttonResources
	List        *listResources
	Backgrounds *backgroundResources
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

type listResources struct {
	UnselectedColor *color.NRGBA
}

type backgroundResources struct {
	MainMenuBG *ebiten.Image
	PauseBG    *ebiten.Image
	PromptIcon *ebiten.Image
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

func newListResources() *listResources {
	return &listResources{
		UnselectedColor: &listIdleColor,
	}
}

func newBackgroundResources() *backgroundResources {
	return &backgroundResources{
		MainMenuBG: utils.LoadImage("ui/desk.jpg"),
		PauseBG:    utils.LoadImage("ui/door.jpg"),
		PromptIcon: utils.LoadImage("ui/prompt.png"),
	}
}
