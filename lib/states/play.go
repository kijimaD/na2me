package states

import (
	"image/color"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kijimaD/nova/event"
	"github.com/kijimaD/nova/lexer"
	"github.com/kijimaD/nova/parser"
)

const (
	screenWidth  = 720
	screenHeight = 720
	fontSize     = 26
	padding      = 40
)

type PlayState struct {
	scenario   []byte
	startLabel *string
	trans      *Transition
	eventQ     event.Queue

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
}

func (st *PlayState) OnStop() {}

func (st *PlayState) Update() Transition {
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
		vector.DrawFilledRect(screen, 0, 0, screenWidth, screenHeight, black, false)
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
		x, y := padding, padding
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		op.LineSpacing = lineSpacing
		text.Draw(screen, japaneseText, st.faceFont, op)
	}
}
