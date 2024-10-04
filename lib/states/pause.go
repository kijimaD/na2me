package states

import (
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/kijimaD/nova/event"
	"github.com/kijimaD/nova/lexer"
	"github.com/kijimaD/nova/parser"
	"github.com/kijimaD/nova/utils"
)

type PauseState struct {
	scenario []byte
	trans    *Transition
	faceFont text.Face

	ui     *ebitenui.UI
	labels []string
}

func (st *PauseState) OnPause() {}

func (st *PauseState) OnResume() {}

func (st *PauseState) OnStart() {
	if len(st.scenario) == 0 {
		log.Fatal("シナリオが選択されていない")
	}

	st.faceFont = loadFont("ui/JF-Dot-Kappa20B.ttf", fontSize)
	l := lexer.NewLexer(string(st.scenario))
	p := parser.NewParser(l)
	program, err := p.ParseProgram()
	if err != nil {
		log.Fatal(err)
	}
	e := event.NewEvaluator()
	e.Eval(program)
	st.labels = e.Labels()

	st.ui = st.initUI()
}

func (st *PauseState) OnStop() {}

func (st *PauseState) Update() Transition {
	st.ui.Update()

	// transの書き換えで遷移する
	if st.trans != nil {
		next := *st.trans
		st.trans = nil
		return next
	}

	return Transition{Type: TransNone}
}

func (st *PauseState) Draw(screen *ebiten.Image) {
	st.ui.Draw(screen)
}

func (st *PauseState) initUI() *ebitenui.UI {
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
	rootContainer.AddChild(st.mainMenuButton(faceFont))
	for _, label := range st.labels {
		rootContainer.AddChild(st.labelSelectButton(label, faceFont, st.scenario))
	}

	return &ebitenui.UI{Container: rootContainer}
}

func (st *PauseState) labelSelectButton(startLabel string, face text.Face, scenario []byte) *widget.Button {
	buttonImage, _ := loadButtonImage()
	button := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text(startLabel, face, &widget.ButtonTextColor{
			Idle: color.RGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			st.trans = &Transition{Type: TransSwitch, NewStates: []State{&PlayState{scenario: scenario, startLabel: utils.GetPtr(startLabel)}}}
		}),
	)
	return button
}

func (st *PauseState) mainMenuButton(face text.Face) *widget.Button {
	buttonImage, _ := loadButtonImage()
	button := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("メインメニュー", face, &widget.ButtonTextColor{
			Idle: color.RGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			st.trans = &Transition{Type: TransSwitch, NewStates: []State{&MainMenuState{}}}
		}),
	)
	return button
}
