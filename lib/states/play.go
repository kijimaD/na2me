package states

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kijimaD/na2me/lib/touch"
	"github.com/kijimaD/na2me/lib/utils"
	"github.com/kijimaD/nova/event"
	"github.com/kijimaD/nova/lexer"
	"github.com/kijimaD/nova/parser"
)

const (
	screenWidth  = 720
	screenHeight = 720
	fontSize     = 26
	padding      = 60
	paddingSmall = 30
)

type PlayState struct {
	ui             *ebitenui.UI
	statsContainer *widget.Container

	// 選択中のシナリオファイルのバイト列
	scenario []byte
	// 指定された章で再生開始する。外部ステートから指定するときに使う
	startLabel *string
	// ステート遷移
	trans *Transition
	// イベントキュー
	eventQ event.Queue
	// アニメーション状態が切り替わったかを判断する用
	prevOnAnim bool

	bgImage     *ebiten.Image
	promptImage *ebiten.Image
	startTime   time.Time
}

func (st *PlayState) OnPause() {}

func (st *PlayState) OnResume() {}

func (st *PlayState) OnStart() {
	if len(st.scenario) == 0 {
		log.Fatal("シナリオが選択されていない")
	}

	st.startTime = time.Now()

	l := lexer.NewLexer(string(st.scenario))
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
		st.eventQ.Evaluator.Play(*st.startLabel)
	}

	{
		eimg := utils.LoadImage("black.png")
		if err != nil {
			log.Fatal(err)
		}
		st.bgImage = eimg
	}
	{
		eimg := utils.LoadImage("ui/prompt.png")
		if err != nil {
			log.Fatal(err)
		}
		st.promptImage = eimg
	}

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
		return Transition{Type: TransPush, NewStates: []State{&PauseState{scenario: st.scenario}}}
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
			eimg := utils.LoadImage(event.Source)
			st.bgImage = eimg
		}
	default:
	}

	// 状態が切り替わったときだけ実行する
	if st.prevOnAnim != st.eventQ.OnAnim {
		st.updateStatsContainer()

		st.prevOnAnim = !st.prevOnAnim
	}

	return Transition{Type: TransNone}
}

func (st *PlayState) Draw(screen *ebiten.Image) {
	{
		// 背景画像
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, screenHeight/4)
		screen.DrawImage(st.bgImage, op)
	}

	{
		// 背景色
		black := color.RGBA{0x10, 0x10, 0x10, 0x80}
		vector.DrawFilledRect(screen, paddingSmall, padding, screenWidth-paddingSmall*2, screenHeight-padding*2, black, false)
	}

	// 待ち状態表示
	if st.eventQ.OnAnim {
		elapsed := time.Since(st.startTime).Seconds()
		offsetY := 4 * math.Sin(elapsed*4) // sin関数で上下に動かす
		bounds := st.promptImage.Bounds()
		bounds.Min.Y = int(20 + offsetY) // 初期位置 + オフセット
		bounds.Max.Y = bounds.Min.Y + bounds.Dy()

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(720-float64(bounds.Max.X*2), 720-float64(bounds.Min.Y*2))
		screen.DrawImage(st.promptImage, op)
	}

	{
		japaneseText := st.eventQ.Display()
		const lineSpacing = fontSize + 8
		x, y := padding-20, padding+paddingSmall
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		op.LineSpacing = lineSpacing
		text.Draw(screen, japaneseText, utils.BodyFont, op)
	}

	st.ui.Draw(screen)
}

func (st *PlayState) initUI() *ebitenui.UI {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Top:    10,
				Bottom: 10,
				Left:   10,
				Right:  10,
			}),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)

	topContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Top:    10,
				Bottom: 10,
				Left:   10,
				Right:  10,
			}),
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
		widget.ButtonOpts.Image(utils.LoadButtonImage()),
		widget.ButtonOpts.Text("一覧", utils.BodyFont, &widget.ButtonTextColor{
			Idle: color.RGBA{0xaa, 0xaa, 0xaa, 0xff},
		}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   10,
			Right:  10,
			Top:    5,
			Bottom: 5,
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			st.trans = &Transition{Type: TransPush, NewStates: []State{&PauseState{scenario: st.scenario}}}
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
		widget.TextOpts.Text(st.eventQ.Evaluator.CurrentLabel, utils.BodyFont, color.NRGBA{100, 100, 100, 255}),
	)

	progressbar := widget.NewProgressBar(
		widget.ProgressBarOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(300, 16),
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter},
			),
		),
		widget.ProgressBarOpts.Images(
			&widget.ProgressBarImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{40, 40, 40, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{40, 40, 40, 255}),
			},
			&widget.ProgressBarImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			},
		),
		widget.ProgressBarOpts.Values(0, len(st.eventQ.Evaluator.Events), st.eventQ.Evaluator.CurrentEventIdx+1),
	)
	rate := float64(st.eventQ.Evaluator.CurrentEventIdx+1) / float64(len(st.eventQ.Evaluator.Events)) * 100
	progressBarLabel := widget.NewText(
		widget.TextOpts.Text(
			fmt.Sprintf("%.1f%%", rate),
			utils.BodyFont,
			color.NRGBA{100, 100, 100, 255},
		),
	)

	st.statsContainer.AddChild(text, progressbar, progressBarLabel)
}
