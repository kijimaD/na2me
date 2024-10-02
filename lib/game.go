package lib

import (
	_ "embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/na2me/lib/states"
)

const (
	screenWidth  = 720
	screenHeight = 720
	fontSize     = 26
	padding      = 40
)

type Game struct {
	StateMachine states.StateMachine
}

func (g *Game) Update() error {
	g.StateMachine.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.StateMachine.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Start() {
	g.StateMachine = states.Init(&states.MainMenuState{})
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("demo")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
