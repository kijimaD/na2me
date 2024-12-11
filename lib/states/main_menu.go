package states

import (
	"fmt"
	"strings"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	embeds "github.com/kijimaD/na2me/embeds"
	"github.com/kijimaD/na2me/lib/eui"
	"github.com/kijimaD/na2me/lib/utils"
)

type MainMenuState struct {
	ui            *ebitenui.UI
	trans         *Transition
	rootContainer *widget.Container

	bgImage *ebiten.Image
}

func (st *MainMenuState) OnPause() {}

func (st *MainMenuState) OnResume() {}

func (st *MainMenuState) OnStart() {
	st.ui = st.initUI()
	st.bgImage = utils.LoadImage("ui/desk.jpg")
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
	screen.DrawImage(st.bgImage, nil)
	st.ui.Draw(screen)
}

func (st *MainMenuState) updateMenuContainer() {}

func (st *MainMenuState) initUI() *ebitenui.UI {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(1),
				widget.GridLayoutOpts.Spacing(4, 4),
				widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true}),
				widget.GridLayoutOpts.Padding(widget.Insets{
					Top:    4,
					Bottom: 4,
					Left:   4,
					Right:  4,
				}),
			),
		),
	)
	listContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewAnchorLayout(),
		),
	)
	entries := []any{}
	for _, s := range embeds.ScenarioMaster.Scenarios {
		entries = append(entries, s.ID)
	}
	list := eui.NewList(
		widget.ListOpts.Entries(entries),
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			key := e.(string)
			scenario := embeds.ScenarioMaster.GetScenario(key)

			whitespace := strings.Repeat("　", 20-(len([]rune(scenario.Title))+len([]rune(scenario.AuthorName))))

			return fmt.Sprintf("%s%s%s", scenario.Title, whitespace, scenario.AuthorName)
		}),
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			key := args.Entry.(string)
			scenario := embeds.ScenarioMaster.GetScenario(key)

			st.trans = &Transition{Type: TransSwitch, NewStates: []State{&PlayState{scenario: scenario.Body}}}
		}),
	)
	listContainer.AddChild(list)
	rootContainer.AddChild(listContainer)

	return &ebitenui.UI{Container: rootContainer}
}
