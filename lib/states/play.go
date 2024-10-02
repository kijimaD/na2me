package states

import (
	"bytes"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	embeds "github.com/kijimaD/na2me/embeds"
	"github.com/kijimaD/nova/event"
	"github.com/kijimaD/nova/lexer"
	"github.com/kijimaD/nova/parser"
	"golang.org/x/text/language"
)

const (
	screenWidth  = 720
	screenHeight = 720
	fontSize     = 26
	padding      = 40
)

var japaneseFaceSource *text.GoTextFaceSource
var eventQ event.Queue

type PlayState struct {
	trans *Transition

	bgImage     *ebiten.Image
	promptImage *ebiten.Image
	startTime   time.Time
}

func (st *PlayState) OnPause() {}

func (st *PlayState) OnResume() {}

func (st *PlayState) OnStart() {
	st.startTime = time.Now()
	{
		font, err := embeds.FS.ReadFile("ui/JF-Dot-Kappa20B.ttf")
		if err != nil {
			log.Fatal(err)
		}
		s, err := text.NewGoTextFaceSource(bytes.NewReader(font))
		if err != nil {
			log.Fatal(err)
		}
		japaneseFaceSource = s
	}

	l := lexer.NewLexer(string(embeds.Input))
	p := parser.NewParser(l)
	program, err := p.ParseProgram()
	if err != nil {
		log.Fatal(err)
	}
	e := event.NewEvaluator()
	e.Eval(program)
	eventQ = event.NewQueue(e)
	eventQ.Start()

	{
		eimg, err := loadImage("forest.jpg")
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
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return Transition{Type: TransQuit}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		eventQ.Run()
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		eventQ.Reset()
	}

	select {
	case v := <-eventQ.NotifyChan:
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
		black := color.RGBA{0x00, 0x00, 0x00, 0x80}
		vector.DrawFilledRect(screen, 0, 0, screenWidth, screenHeight, black, false)
	}

	f := &text.GoTextFace{
		Source:   japaneseFaceSource,
		Size:     fontSize,
		Language: language.Japanese,
	}

	// 待ち状態表示
	if eventQ.OnAnim {
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
		japaneseText := eventQ.Display()
		const lineSpacing = fontSize + 8
		x, y := padding, padding
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		op.LineSpacing = lineSpacing
		text.Draw(screen, japaneseText, f, op)
	}
}

func (st *PlayState) updateMenuContainer() {}
