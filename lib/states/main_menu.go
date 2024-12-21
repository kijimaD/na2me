package states

import (
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.TrackHover(false)),
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(1),
				widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
				widget.GridLayoutOpts.Padding(widget.Insets{
					Top:    10,
					Bottom: 10,
				}),
				widget.GridLayoutOpts.Spacing(0, 0),
			),
		),
	)

	footerContainer := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewRowLayout(
		widget.RowLayoutOpts.Padding(widget.Insets{
			Left:  25,
			Right: 25,
		}),
	)))
	footerContainer.AddChild(widget.NewText(
		widget.TextOpts.Text("github.com/kijimaD/na2me -- 電子紙芝居方式流通推進連盟", utils.UIFont, color.NRGBA{100, 100, 100, 255})))

	var ui *ebitenui.UI
	rootContainer.AddChild(
		st.headerContainer(),
		st.actionContainer(func() *ebitenui.UI {
			return ui
		}),
		footerContainer,
	)

	return &ebitenui.UI{Container: rootContainer}
}

func (st *MainMenuState) headerContainer() widget.PreferredSizeLocateableWidget {
	c := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10))),
	)

	c.AddChild(st.header("話灯機",
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
	))

	c2 := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Top:    10,
				Bottom: 10,
				Left:   25,
				Right:  25,
			}),
		)),
	)
	c.AddChild(c2)

	c2.AddChild(widget.NewText(
		widget.TextOpts.Text("注意力散漫たる現代において、歴史的読書方法は競争力を失っている。\n電子紙芝居方式の優れた威力を万人へ宣伝し、方式普及を推進する。", utils.UIFont, color.NRGBA{100, 100, 100, 255})),
	)

	return c
}

func (st *MainMenuState) header(label string, opts ...widget.ContainerOpt) widget.PreferredSizeLocateableWidget {
	c := widget.NewContainer(append(opts, []widget.ContainerOpt{
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.TrackHover(false)),
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{255, 255, 255, 140})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(widget.AnchorLayoutOpts.Padding(widget.Insets{
			Top:    4,
			Bottom: 4,
			Left:   30,
			Right:  30,
		}))),
	}...)...)

	c.AddChild(widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionStart,
			VerticalPosition:   widget.AnchorLayoutPositionStart,
		})),
		widget.TextOpts.Text(label, utils.UIFont, color.NRGBA{0, 0, 0, 200}),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter),
	))

	return c
}

func (st *MainMenuState) actionContainer(ui func() *ebitenui.UI) widget.PreferredSizeLocateableWidget {
	actionContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Padding(widget.Insets{
				Left:  25,
				Right: 25,
			}),
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Stretch([]bool{false, true}, []bool{true}),
			widget.GridLayoutOpts.Spacing(20, 0),
		)))

	pages := []interface{}{
		st.recentPage(),
		st.bookListPage(),
	}

	pageContainer := st.newPageContainer()

	pageList := eui.NewList(
		widget.ListOpts.Entries(pages),
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			return e.(*page).title
		}),
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			pageContainer.setPage(args.Entry.(*page))
		}),
	)
	pageList.SetSelectedEntry(pages[0])

	actionContainer.AddChild(
		pageList,
		pageContainer.widget,
	)

	return actionContainer
}

func (st *MainMenuState) newPageContainer() *pageContainer {
	c := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.TrackHover(false)),
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{255, 255, 255, 140})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(2)),
			widget.RowLayoutOpts.Spacing(15))),
	)

	titleText := widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.TextOpts.Text("", utils.UIFont, color.NRGBA{255, 255, 255, 255}))
	c.AddChild(titleText)

	flipBook := widget.NewFlipBook(
		widget.FlipBookOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		}))),
	)
	c.AddChild(flipBook)

	return &pageContainer{
		widget:    c,
		titleText: titleText,
		flipBook:  flipBook,
	}
}

func (p *pageContainer) setPage(page *page) {
	p.titleText.Label = page.title
	p.flipBook.SetPage(page.content)
	p.flipBook.RequestRelayout()
}
