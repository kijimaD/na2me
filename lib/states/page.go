package states

import (
	"fmt"
	"strings"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/kijimaD/na2me/lib/bookmark"
	"github.com/kijimaD/na2me/lib/eui"
	"github.com/kijimaD/na2me/lib/resources"
	"github.com/kijimaD/na2me/lib/scenario"
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
			key := e.(scenario.ScenarioIDType)
			bookmark, ok := bookmark.Bookmarks.Get(key)
			if !ok {
				return ""
			}
			whitespace := strings.Repeat("　", 18-(len([]rune(bookmark.ScenarioName))+len([]rune(bookmark.Label))))
			return fmt.Sprintf("%s%s%s", bookmark.ScenarioName, whitespace, bookmark.Label)
		}),
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			key := args.Entry.(scenario.ScenarioIDType)
			scenario := scenario.ScenarioMaster.GetScenario(key)
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
			widget.TextOpts.Text("保存なし", resources.Master.Fonts.UIFace, resources.TextPrimaryColor),
		)
		listContainer.AddChild(noContentText)
	}
	c.AddChild(listContainer)

	return &page{
		title:   "再開",
		content: c,
	}
}

func (st *MainMenuState) bookListPage() *page {
	c := newPageContentContainer()

	entries := []any{}
	for _, s := range scenario.ScenarioMaster.Scenarios {
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
			key := e.(scenario.ScenarioIDType)
			scenario := scenario.ScenarioMaster.GetScenario(key)

			whitespace := strings.Repeat("　", 18-(len([]rune(scenario.Title))+len([]rune(scenario.AuthorName))))

			return fmt.Sprintf("%s%s%s", scenario.Title, whitespace, scenario.AuthorName)
		}),
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			key := args.Entry.(scenario.ScenarioIDType)
			scenario := scenario.ScenarioMaster.GetScenario(key)

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
		title:   "全部",
		content: c,
	}
}

func (st *MainMenuState) unreadPage() *page {
	c := newPageContentContainer()

	entries := []any{}
	for _, status := range scenario.ScenarioMaster.Statuses {
		if !status.IsRead {
			entries = append(entries, status.ID)
		}
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
			key := e.(scenario.ScenarioIDType)
			scenario := scenario.ScenarioMaster.GetScenario(key)

			whitespace := strings.Repeat("　", 18-(len([]rune(scenario.Title))+len([]rune(scenario.AuthorName))))

			return fmt.Sprintf("%s%s%s", scenario.Title, whitespace, scenario.AuthorName)
		}),
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			key := args.Entry.(scenario.ScenarioIDType)
			scenario := scenario.ScenarioMaster.GetScenario(key)

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
		title:   "未読",
		content: c,
	}
}

func (st *MainMenuState) donePage() *page {
	c := newPageContentContainer()

	entries := []any{}
	for _, status := range scenario.ScenarioMaster.Statuses {
		if status.IsRead {
			entries = append(entries, status.ID)
		}
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
			key := e.(scenario.ScenarioIDType)
			scenario := scenario.ScenarioMaster.GetScenario(key)

			whitespace := strings.Repeat("　", 18-(len([]rune(scenario.Title))+len([]rune(scenario.AuthorName))))

			return fmt.Sprintf("%s%s%s", scenario.Title, whitespace, scenario.AuthorName)
		}),
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			key := args.Entry.(scenario.ScenarioIDType)
			scenario := scenario.ScenarioMaster.GetScenario(key)

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
		title:   "既読",
		content: c,
	}
}

func (st *MainMenuState) infoPage() *page {
	c := newPageContentContainer()

	bookCount := widget.NewText(widget.TextOpts.Text(fmt.Sprintf("収録数 %d", len(scenario.ScenarioMaster.Scenarios)), resources.Master.Fonts.UIFace, resources.TextPrimaryColor))
	c.AddChild(bookCount)

	return &page{
		title:   "情報",
		content: c,
	}
}
