package states

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"path"
	"time"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	embeds "github.com/kijimaD/na2me/embeds"
	"github.com/kijimaD/na2me/lib/resources"
	"github.com/kijimaD/na2me/lib/touch"
	"github.com/kijimaD/na2me/lib/utils"
	"github.com/kijimaD/nova/event"
	"github.com/kijimaD/nova/lexer"
	"github.com/kijimaD/nova/parser"
)

const (
	padding      = 60
	paddingSmall = 30
)

type PlayState struct {
	ui             *ebitenui.UI
	statsContainer *widget.Container

	// 選択中のシナリオ
	scenario embeds.Scenario
	// 指定された章で再生開始する。外部ステートから指定するときに使う
	startLabel *string

	// ステート遷移
	trans *Transition
	// イベントキュー
	eventQ event.Queue
	// アニメーション状態が切り替わったかを判断する用
	prevOnAnim bool

	bgImage   *ebiten.Image
	startTime time.Time
}

func (st *PlayState) OnPause() {}

func (st *PlayState) OnResume() {}

func (st *PlayState) OnStart() {
	if len(st.scenario.Body) == 0 {
		log.Fatal("シナリオが選択されていない")
	}

	st.startTime = time.Now()

	l := lexer.NewLexer(string(st.scenario.Body))
	p := parser.NewParser(l)
	program, err := p.ParseProgram()
	if err != nil {
		log.Fatal(err)
	}
	e := event.NewEvaluator()
	e.Eval(program)
	st.eventQ = event.NewQueue(e)
	st.eventQ.Start()
	if st.startLabel != nil {
		st.eventQ.Play(*st.startLabel)
	}
	st.bgImage = utils.LoadImage("bg/black.png")

	st.ui = st.initUI()
}

func (st *PlayState) OnStop() {}

func (st *PlayState) Update() Transition {
	st.ui.Update()

	// transの書き換えで遷移する
	if st.trans != nil {
		next := *st.trans
		st.trans = nil
		return next
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		st.transPause()
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		st.eventQ.Run()
	}

	if touch.IsTouchJustReleased() {
		st.eventQ.Run()
	}

	select {
	case v := <-st.eventQ.NotifyChan:
		switch event := v.(type) {
		case *event.ChangeBg:
			st.bgImage = utils.LoadImage(path.Join("bg", event.Source))
		}
	default:
	}

	// 状態が切り替わったときだけ実行する
	if st.prevOnAnim != st.eventQ.OnAnim {
		st.updateStatsContainer()

		st.startTime = time.Now()
		st.prevOnAnim = !st.prevOnAnim
	}

	return Transition{Type: TransNone}
}

func (st *PlayState) Draw(screen *ebiten.Image) {
	{
		// 背景画像
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, resources.ScreenHeight/4)
		screen.DrawImage(st.bgImage, op)
	}

	{
		// 背景色
		black := color.RGBA{0x10, 0x10, 0x10, 0x80}
		vector.DrawFilledRect(screen, paddingSmall, padding, resources.ScreenWidth-paddingSmall*2, resources.ScreenHeight-padding*2, black, false)
	}

	// 待ち状態表示
	if st.eventQ.OnAnim {
		promptImage := resources.Master.Backgrounds.PromptIcon
		elapsed := time.Since(st.startTime).Seconds()
		offsetY := 2 * math.Cos(elapsed*4) // cos関数で上下に動かす
		bounds := promptImage.Bounds()
		bounds.Min.Y = int(20 + offsetY) // 初期位置 + オフセット
		bounds.Max.Y = bounds.Min.Y + bounds.Dy()

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(720-float64(bounds.Max.X*2), 720-float64(bounds.Min.Y*2))
		screen.DrawImage(promptImage, op)
	}

	{
		japaneseText := st.eventQ.Display()
		const fontSize = 26
		const lineSpacing = fontSize + 8
		x, y := padding-20, padding+paddingSmall
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		op.LineSpacing = lineSpacing
		text.Draw(screen, japaneseText, resources.Master.Fonts.BodyFace, op)
	}

	st.ui.Draw(screen)
}

func (st *PlayState) transPause() {
	newState := NewPauseState(
		st.scenario,
		st.eventQ.CurrentLabel,
	)
	st.trans = &Transition{Type: TransPush, NewStates: []State{&newState}}
}

func (st *PlayState) initUI() *ebitenui.UI {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(10)),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)

	topContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(10)),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)

	button := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		widget.ButtonOpts.Image(resources.Master.Button.Image),
		widget.ButtonOpts.Text(
			"一覧",
			resources.Master.Button.Face,
			resources.Master.Button.TextColor,
		),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   10,
			Right:  10,
			Top:    6,
			Bottom: 6,
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			st.transPause()
		}),
	)

	st.statsContainer = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)
	st.updateStatsContainer()

	topContainer.AddChild(button, st.statsContainer)
	rootContainer.AddChild(topContainer)

	return &ebitenui.UI{Container: rootContainer}
}

func (st *PlayState) updateStatsContainer() {
	st.statsContainer.RemoveChildren()

	text := widget.NewText(
		widget.TextOpts.Text(st.eventQ.CurrentLabel, resources.Master.Fonts.BodyFace, color.NRGBA{100, 100, 100, 255}),
	)

	idx := len(st.eventQ.Evaluator.Events) - len(st.eventQ.WaitingQueue)
	progressbar := widget.NewProgressBar(
		widget.ProgressBarOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(300, 16),
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter},
			),
		),
		widget.ProgressBarOpts.Images(
			resources.Master.ProgressBar.TrackImage,
			resources.Master.ProgressBar.FillImage,
		),
		widget.ProgressBarOpts.Values(0, len(st.eventQ.Evaluator.Events), idx+1),
	)
	rate := float64(idx+1) / float64(len(st.eventQ.Evaluator.Events)) * 100
	progressBarLabel := widget.NewText(
		widget.TextOpts.Text(
			fmt.Sprintf("%.1f%%", rate),
			resources.Master.Fonts.BodyFace,
			resources.TextSecondaryColor,
		),
	)

	st.statsContainer.AddChild(text, progressbar, progressBarLabel)
}
