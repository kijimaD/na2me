package states

import (
	"fmt"
	"strings"

	"github.com/ebitenui/ebitenui/widget"
	embeds "github.com/kijimaD/na2me/embeds"
	"github.com/kijimaD/na2me/lib/eui"
)

type page struct {
	title   string
	content widget.PreferredSizeLocateableWidget
}

type pageContainer struct {
	widget    widget.PreferredSizeLocateableWidget
	titleText *widget.Text
	flipBook  *widget.FlipBook
}

func newPageContentContainer() *widget.Container {
	return widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
		)))
}

func (st *MainMenuState) bookListPage() *page {
	c := newPageContentContainer()

	entries := []any{}
	for _, s := range embeds.ScenarioMaster.Scenarios {
		entries = append(entries, s.ID)
	}

	listContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(1),
				widget.GridLayoutOpts.Stretch([]bool{false}, []bool{false}),
				widget.GridLayoutOpts.Spacing(10, 0))),
	)

	list := eui.NewList(
		widget.ListOpts.Entries(entries),
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			key := e.(string)
			scenario := embeds.ScenarioMaster.GetScenario(key)

			whitespace := strings.Repeat("　", 18-(len([]rune(scenario.Title))+len([]rune(scenario.AuthorName))))

			return fmt.Sprintf("%s%s%s", scenario.Title, whitespace, scenario.AuthorName)
		}),
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			key := args.Entry.(string)
			scenario := embeds.ScenarioMaster.GetScenario(key)

			st.trans = &Transition{Type: TransSwitch, NewStates: []State{&PlayState{scenario: scenario}}}
		}),
		widget.ListOpts.ContainerOpts(
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.GridLayoutData{
					MaxHeight: 560,
				}),
			)),
	)
	listContainer.AddChild(list)
	c.AddChild(listContainer)

	return &page{
		title:   "作品一覧",
		content: c,
	}
}
