package states

import (
	"image/color"
	"log"
	"math"
	"time"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kijimaD/na2me/lib/touch"
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
	ui *ebitenui.UI

	// 選択中のシナリオファイルのバイト列
	scenario []byte
	// 指定された章で再生開始する。外部ステートから指定するときに使う
	startLabel *string
	// ステート遷移
	trans *Transition
	// イベントキュー
	eventQ event.Queue

	bgImage     *ebiten.Image
	promptImage *ebiten.Image
	startTime   time.Time
	faceFont    text.Face
}

func (st *PlayState) OnPause() {}

func (st *PlayState) OnResume() {}

func (st *PlayState) OnStart() {
	if len(st.scenario) == 0 {
		log.Fatal("シナリオが選択されていない")
	}

	st.faceFont = loadFont("ui/JF-Dot-Kappa20B.ttf", fontSize)
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
		eimg, err := loadImage("black.png")
		if err != nil {
			log.Fatal(err)
		}
		st.bgImage = eimg
	}
	{
		eimg, err := loadImage("ui/prompt.png")
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
			eimg, err := loadImage(event.Source)
			if err != nil {
				log.Fatal(err)
			}
			st.bgImage = eimg
		}
	default:
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
		vector.DrawFilledRect(screen, paddingSmall, paddingSmall, screenWidth-paddingSmall*2, screenHeight-paddingSmall*2, black, false)
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
		x, y := padding-20, padding
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		op.LineSpacing = lineSpacing
		text.Draw(screen, japaneseText, st.faceFont, op)
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

	buttonImage, _ := loadButtonImage()
	button := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("一覧", st.faceFont, &widget.ButtonTextColor{
			Idle: color.RGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			st.trans = &Transition{Type: TransPush, NewStates: []State{&PauseState{scenario: st.scenario}}}
		}),
	)
	rootContainer.AddChild(button)

	return &ebitenui.UI{Container: rootContainer}
}
