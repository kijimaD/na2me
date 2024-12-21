package states

import (
	"fmt"
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	embeds "github.com/kijimaD/na2me/embeds"
	"github.com/kijimaD/na2me/lib/bookmark"
	"github.com/kijimaD/na2me/lib/eui"
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
	bgImage       *ebiten.Image
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
	st.bgImage = utils.LoadImage("ui/door.jpg")
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
	screen.DrawImage(st.bgImage, nil)
	st.ui.Draw(screen)
}

func (st *PauseState) initUI() *ebitenui.UI {
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
	st.rootContainer = rootContainer
	st.reloadUI()

	return &ebitenui.UI{Container: st.rootContainer}
}

func (st *PauseState) reloadUI() {
	st.rootContainer.RemoveChildren()

	buttonContainer := widget.NewContainer(
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
	emptyContainer := widget.NewContainer(
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
	title := widget.NewText(
		widget.TextOpts.Text(st.scenario.Title, utils.BodyFont, color.White),
	)
	buttonContainer.AddChild(
		title,
		st.backButton(utils.BodyFont),
		st.mainMenuButton(utils.BodyFont),
		emptyContainer,
		st.saveButton(utils.BodyFont),
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

func (st *PauseState) mainMenuButton(face text.Face) *widget.Button {
	button := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		widget.ButtonOpts.Image(utils.LoadButtonImage()),
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
	button := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		widget.ButtonOpts.Image(utils.LoadButtonImage()),
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

func (st *PauseState) saveButton(face text.Face) *widget.Button {
	button := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		widget.ButtonOpts.Image(utils.LoadButtonImage()),
		widget.ButtonOpts.Text("保存", face, &widget.ButtonTextColor{
			Idle: color.RGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			bm := bookmark.NewBookmark(
				st.scenario.ID,
				st.scenario.Title,
				st.currentLabel,
			)
			bookmark.Bookmarks.Add(bm)
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
		widget.TextOpts.Text(str, utils.BodyFont, color.NRGBA{100, 100, 100, 255}),
	)

	return text
}
