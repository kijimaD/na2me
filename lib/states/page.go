package states

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/ebitenui/ebitenui/widget"
	embeds "github.com/kijimaD/na2me/embeds"
	"github.com/kijimaD/na2me/lib/bookmark"
	"github.com/kijimaD/na2me/lib/eui"
	"github.com/kijimaD/na2me/lib/utils"
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
			key := e.(embeds.ScenarioIDType)
			scenario := embeds.ScenarioMaster.GetScenario(key)

			whitespace := strings.Repeat("　", 18-(len([]rune(scenario.Title))+len([]rune(scenario.AuthorName))))

			return fmt.Sprintf("%s%s%s", scenario.Title, whitespace, scenario.AuthorName)
		}),
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			key := args.Entry.(embeds.ScenarioIDType)
			scenario := embeds.ScenarioMaster.GetScenario(key)

			st.trans = &Transition{Type: TransSwitch, NewStates: []State{&PlayState{scenario: scenario}}}
		}),
		widget.ListOpts.ContainerOpts(
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.GridLayoutData{
					MaxHeight: 520,
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

func (st *MainMenuState) recentPage() *page {
	c := newPageContentContainer()

	entries := []any{}
	for _, s := range bookmark.Bookmarks.Bookmarks {
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
			key := e.(embeds.ScenarioIDType)
			bookmark, ok := bookmark.Bookmarks.Get(key)
			if !ok {
				return ""
			}
			whitespace := strings.Repeat("　", 18-(len([]rune(bookmark.ScenarioName))+len([]rune(bookmark.Label))))
			return fmt.Sprintf("%s%s%s", bookmark.ScenarioName, whitespace, bookmark.Label)
		}),
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			key := args.Entry.(embeds.ScenarioIDType)
			scenario := embeds.ScenarioMaster.GetScenario(key)
			bm, ok := bookmark.Bookmarks.Get(key)
			if !ok {
				return
			}

			st.trans = &Transition{Type: TransSwitch, NewStates: []State{&PlayState{scenario: scenario, startLabel: utils.GetPtr(bm.Label)}}}
		}),
		widget.ListOpts.ContainerOpts(
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.GridLayoutData{
					MaxHeight: 520,
				}),
			)),
	)
	if len(entries) > 0 {
		listContainer.AddChild(list)
	} else {
		noContentText := widget.NewText(
			widget.TextOpts.Text("保存なし", utils.UIFont, color.NRGBA{255, 255, 255, 255}),
		)
		listContainer.AddChild(noContentText)
	}
	c.AddChild(listContainer)

	return &page{
		title:   "再開",
		content: c,
	}
}