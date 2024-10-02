package states

import (
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	embeds "github.com/kijimaD/na2me/embeds"
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

	return Transition{Type: TransNone}
}

func (st *MainMenuState) Draw(screen *ebiten.Image) {
	st.ui.Draw(screen)
}

func (st *MainMenuState) updateMenuContainer() {}

func (st *MainMenuState) initUI() *ebitenui.UI {
	faceFont := loadFont("ui/JF-Dot-Kappa20B.ttf", 20)

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.RGBA{0x13, 0x1a, 0x22, 0xff})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Top:    10,
				Bottom: 10,
				Left:   10,
				Right:  10,
			}),
			widget.RowLayoutOpts.Spacing(10), // ボタンの間隔
		)),
	)
	rootContainer.AddChild(st.scenarioSelectButton("坊っちゃん", faceFont, embeds.Bochan))
	rootContainer.AddChild(st.scenarioSelectButton("吾輩は猫である", faceFont, embeds.WagahaiHaNekoDearu))

	return &ebitenui.UI{Container: rootContainer}
}

func (st *MainMenuState) scenarioSelectButton(label string, face text.Face, scenario []byte) *widget.Button {
	buttonImage, _ := loadButtonImage()
	button := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text(label, face, &widget.ButtonTextColor{
			Idle: color.RGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			st.trans = &Transition{Type: TransSwitch, NewStates: []State{&PlayState{scenario: scenario}}}
		}),
	)
	return button
}

func loadButtonImage() (*widget.ButtonImage, error) {
	idle := image.NewNineSliceColor(color.RGBA{R: 110, G: 110, B: 180, A: 255})
	hover := image.NewNineSliceColor(color.RGBA{R: 110, G: 180, B: 130, A: 255})
	pressed := image.NewNineSliceColor(color.RGBA{R: 80, G: 110, B: 100, A: 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil
}
