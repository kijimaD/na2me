package states

import (
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return Transition{Type: TransPop, NewStates: []State{}}
	}

	return Transition{Type: TransNone}
}

func (st *PauseState) Draw(screen *ebiten.Image) {
	st.ui.Draw(screen)
}

func (st *PauseState) initUI() *ebitenui.UI {
	faceFont := loadFont("ui/JF-Dot-Kappa20B.ttf", 20)
	buttonImage, _ := loadButtonImage()

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
	entries := []any{}
	for _, label := range st.labels {
		entries = append(entries, label)
	}
	list := widget.NewList(
		widget.ListOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(150, 0),
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				StretchVertical:    true,
				Padding:            widget.NewInsetsSimple(50),
			}),
		)),
		widget.ListOpts.Entries(entries),
		widget.ListOpts.ScrollContainerOpts(
			widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
				Idle:     image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Disabled: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Mask:     image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			}),
		),
		widget.ListOpts.SliderOpts(
			widget.SliderOpts.Images(&widget.SliderTrackImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			}, buttonImage),
			widget.SliderOpts.MinHandleSize(5),
			widget.SliderOpts.TrackPadding(widget.NewInsetsSimple(2))),
		widget.ListOpts.HideHorizontalSlider(),
		widget.ListOpts.EntryFontFace(faceFont),
		widget.ListOpts.EntryColor(&widget.ListEntryColor{
			Selected:                   color.NRGBA{R: 0, G: 255, B: 0, A: 255},
			Unselected:                 color.NRGBA{R: 254, G: 255, B: 255, A: 255},
			SelectedBackground:         color.NRGBA{R: 130, G: 130, B: 200, A: 255},
			SelectingBackground:        color.NRGBA{R: 130, G: 130, B: 130, A: 255},
			SelectingFocusedBackground: color.NRGBA{R: 130, G: 140, B: 170, A: 255},
			SelectedFocusedBackground:  color.NRGBA{R: 130, G: 130, B: 170, A: 255},
			FocusedBackground:          color.NRGBA{R: 170, G: 170, B: 180, A: 255},
			DisabledUnselected:         color.NRGBA{R: 100, G: 100, B: 100, A: 255},
			DisabledSelected:           color.NRGBA{R: 100, G: 100, B: 100, A: 255},
			DisabledSelectedBackground: color.NRGBA{R: 100, G: 100, B: 100, A: 255},
		}),
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			key := e.(string)

			return key
		}),
		widget.ListOpts.EntryTextPadding(widget.NewInsetsSimple(5)),
		widget.ListOpts.EntryTextPosition(widget.TextPositionStart, widget.TextPositionCenter),
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			key := args.Entry.(string)
			st.trans = &Transition{Type: TransSwitch, NewStates: []State{&PlayState{scenario: st.scenario, startLabel: utils.GetPtr(key)}}}
		}),
	)
	rootContainer.AddChild(st.backButton(faceFont))
	rootContainer.AddChild(list)
	rootContainer.AddChild(st.mainMenuButton(faceFont))

	return &ebitenui.UI{Container: rootContainer}
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
		widget.ButtonOpts.Text("終了", face, &widget.ButtonTextColor{
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

func (st *PauseState) backButton(face text.Face) *widget.Button {
	buttonImage, _ := loadButtonImage()
	button := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("戻る", face, &widget.ButtonTextColor{
			Idle: color.RGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			st.trans = &Transition{Type: TransPop}
		}),
	)
	return button
}
