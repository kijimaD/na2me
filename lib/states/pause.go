package states

import (
	"fmt"
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	embeds "github.com/kijimaD/na2me/embeds"
	"github.com/kijimaD/na2me/lib/bookmark"
	"github.com/kijimaD/na2me/lib/eui"
	"github.com/kijimaD/na2me/lib/resources"
	"github.com/kijimaD/na2me/lib/utils"
	"github.com/kijimaD/nova/event"
	"github.com/kijimaD/nova/lexer"
	"github.com/kijimaD/nova/parser"
)

// TODO: Newを作ったほうがよさそう
type PauseState struct {
	trans *Transition

	// 再生シナリオ
	scenario embeds.Scenario
	// 再生中ラベル
	currentLabel string

	// ラベル一覧
	labels []string

	ui            *ebitenui.UI
	rootContainer *widget.Container
}

func NewPauseState(scenario embeds.Scenario, currentLabel string) PauseState {
	return PauseState{
		scenario:     scenario,
		currentLabel: currentLabel,
	}
}

func (st *PauseState) OnPause() {}

func (st *PauseState) OnResume() {}

func (st *PauseState) OnStart() {
	if len(st.scenario.Body) == 0 {
		log.Fatal("シナリオが選択されていない")
	}

	l := lexer.NewLexer(string(st.scenario.Body))
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
	screen.DrawImage(resources.Master.Backgrounds.PauseBG, nil)
	st.ui.Draw(screen)
}

func (st *PauseState) initUI() *ebitenui.UI {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(1),
				widget.GridLayoutOpts.Spacing(4, 4),
				widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true}),
				widget.GridLayoutOpts.Padding(widget.NewInsetsSimple(10)),
			),
		),
	)
	st.rootContainer = rootContainer
	st.reloadUI()

	return &ebitenui.UI{Container: st.rootContainer}
}

func (st *PauseState) reloadUI() {
	st.rootContainer.RemoveChildren()

	buttonContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(10)),
			widget.RowLayoutOpts.Spacing(10), // ボタンの間隔
		)),
	)
	emptyContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(10)),
			widget.RowLayoutOpts.Spacing(10), // ボタンの間隔
		)),
	)
	title := widget.NewText(
		widget.TextOpts.Text(st.scenario.Title, resources.Master.Fonts.BodyFace, resources.TextPrimaryColor),
	)
	buttonContainer.AddChild(
		title,
		st.backButton(),
		st.mainMenuButton(),
		emptyContainer,
		st.saveButton(),
		st.deleteButton(),
		st.saveText(),
	)

	entries := []any{}
	for _, label := range st.labels {
		entries = append(entries, label)
	}
	list := eui.NewList(
		widget.ListOpts.Entries(entries),
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			key := e.(string)
			if key == st.currentLabel {
				key += " ←"
			}

			return key
		}),
		widget.ListOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(300, 0),
		)),
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			key := args.Entry.(string)
			st.trans = &Transition{Type: TransSwitch, NewStates: []State{&PlayState{scenario: st.scenario, startLabel: utils.GetPtr(key)}}}
		}),
	)
	// Listは高さのあるレイアウトのコンテナに入れないと、スクロールされない
	listContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewAnchorLayout(),
		),
	)
	listContainer.AddChild(
		list,
		buttonContainer,
	)

	st.rootContainer.AddChild(listContainer)
}

func (st *PauseState) mainMenuButton() *widget.Button {
	button := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		widget.ButtonOpts.Image(resources.Master.Button.Image),
		widget.ButtonOpts.Text(
			"終了",
			resources.Master.Button.Face,
			resources.Master.Button.TextColor,
		),
		widget.ButtonOpts.TextPadding(resources.Master.Button.Padding),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			st.trans = &Transition{Type: TransSwitch, NewStates: []State{&MainMenuState{}}}
		}),
	)
	return button
}

func (st *PauseState) backButton() *widget.Button {
	button := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		widget.ButtonOpts.Image(resources.Master.Button.Image),
		widget.ButtonOpts.Text(
			"戻る",
			resources.Master.Button.Face,
			resources.Master.Button.TextColor,
		),
		widget.ButtonOpts.TextPadding(resources.Master.Button.Padding),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			st.trans = &Transition{Type: TransPop}
		}),
	)
	return button
}

func (st *PauseState) saveButton() *widget.Button {
	button := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		widget.ButtonOpts.Image(resources.Master.Button.Image),
		widget.ButtonOpts.Text(
			"保存",
			resources.Master.Button.Face,
			resources.Master.Button.TextColor,
		),
		widget.ButtonOpts.TextPadding(resources.Master.Button.Padding),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			bm := bookmark.NewBookmark(
				st.scenario.ID,
				st.scenario.Title,
				st.currentLabel,
			)
			bookmark.Bookmarks.Add(bm)
			if err := bookmark.GlobalSave(); err != nil {
				log.Fatal(err)
			}
			st.reloadUI()
		}),
	)
	return button
}

func (st *PauseState) deleteButton() *widget.Button {
	button := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		widget.ButtonOpts.Image(resources.Master.Button.Image),
		widget.ButtonOpts.Text(
			"破棄",
			resources.Master.Button.Face,
			resources.Master.Button.TextColor,
		),
		widget.ButtonOpts.TextPadding(resources.Master.Button.Padding),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			bookmark.Bookmarks.Delete(st.scenario.ID)
			if err := bookmark.GlobalSave(); err != nil {
				log.Fatal(err)
			}
			st.reloadUI()
		}),
	)
	return button
}

func (st *PauseState) saveText() *widget.Text {
	str := ""
	bookmark, ok := bookmark.Bookmarks.Get(st.scenario.ID)
	if !ok {
		str = "未保存"
	} else {
		str = fmt.Sprintf("%s", bookmark.Label)
	}

	text := widget.NewText(
		widget.TextOpts.Text(str, resources.Master.Fonts.BodyFace, color.NRGBA{100, 100, 100, 255}),
	)

	return text
}
