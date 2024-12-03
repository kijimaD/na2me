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
	"github.com/kijimaD/na2me/lib/eui"
	"github.com/kijimaD/na2me/lib/utils"
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
	for _, s := range embeds.ScenarioMaster.Scenarios {
		entries = append(entries, s.Name)
	}
	list := eui.NewList(
		widget.ListOpts.Entries(entries),
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			key := e.(string)
			scenario := embeds.ScenarioMaster.GetScenario(key)

			return scenario.LabelName
		}),
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			key := args.Entry.(string)
			scenario := embeds.ScenarioMaster.GetScenario(key)

			st.trans = &Transition{Type: TransSwitch, NewStates: []State{&PlayState{scenario: scenario.Body}}}
		}),
	)
	rootContainer.AddChild(list)

	return &ebitenui.UI{Container: rootContainer}
}

func (st *MainMenuState) scenarioSelectButton(label string, face text.Face, scenario []byte) *widget.Button {
	buttonImage := utils.LoadButtonImage()
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
