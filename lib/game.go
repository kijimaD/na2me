package lib

import (
	"bytes"
	_ "embed"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	embeds "github.com/kijimaD/na2me/file"
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

type Game struct {
	bgImage     *ebiten.Image
	promptImage *ebiten.Image
	startTime   time.Time
}

func (g *Game) Update() error {
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
			g.bgImage = eimg
		}
	default:
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	{
		// 背景画像
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, screenHeight/4)
		screen.DrawImage(g.bgImage, op)
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
		elapsed := time.Since(g.startTime).Seconds()
		offsetY := 4 * math.Sin(elapsed*4) // sin関数で上下に動かす
		bounds := g.promptImage.Bounds()
		bounds.Min.Y = int(20 + offsetY) // 初期位置 + オフセット
		bounds.Max.Y = bounds.Min.Y + bounds.Dy()

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(720-float64(bounds.Max.X*2), 720-float64(bounds.Min.Y*2))
		screen.DrawImage(g.promptImage, op)
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

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Start() {
	g.startTime = time.Now()
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
		g.bgImage = eimg
	}
	{
		eimg, err := loadImage("ui/prompt.png")
		if err != nil {
			log.Fatal(err)
		}
		g.promptImage = eimg
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("demo")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
