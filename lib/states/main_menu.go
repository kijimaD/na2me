package states

import (
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type MainMenuState struct {
	ui            *ebitenui.UI
	trans         *Transition
	rootContainer *widget.Container
}

func (st *MainMenuState) OnPause() {}

func (st *MainMenuState) OnResume() {}

func (st *MainMenuState) OnStart() {
	st.ui = st.initUI()
}

func (st *MainMenuState) OnStop() {}

func (st *MainMenuState) Update() Transition {
	st.ui.Update()

	// transの書き換えで遷移する
	if st.trans != nil {
		next := *st.trans
		st.trans = nil
		return next
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return Transition{Type: TransQuit}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return Transition{Type: TransSwitch, NewStates: []State{&PlayState{}}}
	}

	return Transition{Type: TransNone}
}

func (st *MainMenuState) Draw(screen *ebiten.Image) {
	st.ui.Draw(screen)
}

func (st *MainMenuState) updateMenuContainer() {}

func (st *MainMenuState) initUI() *ebitenui.UI {
	faceFont := loadFont("ui/JF-Dot-Kappa20B.ttf", 20)

	buttonImage, _ := loadButtonImage()
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.RGBA{0x13, 0x1a, 0x22, 0xff})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	button := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("坊っちゃん", faceFont, &widget.ButtonTextColor{
			Idle: color.RGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			st.trans = &Transition{Type: TransSwitch, NewStates: []State{&PlayState{}}}
		}),
	)
	rootContainer.AddChild(button)
	return &ebitenui.UI{Container: rootContainer}
}

func loadButtonImage() (*widget.ButtonImage, error) {
	idle := image.NewNineSliceColor(color.RGBA{R: 170, G: 170, B: 180, A: 255})

	hover := image.NewNineSliceColor(color.RGBA{R: 130, G: 130, B: 150, A: 255})

	pressed := image.NewNineSliceColor(color.RGBA{R: 100, G: 100, B: 120, A: 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil
}
