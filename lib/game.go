package lib

import (
	_ "embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/na2me/lib/resources"
	"github.com/kijimaD/na2me/lib/states"
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
	return resources.ScreenWidth, resources.ScreenHeight
}

func (g *Game) Start() {
	g.StateMachine = states.Init(&states.MainMenuState{})
	ebiten.SetWindowSize(resources.ScreenWidth, resources.ScreenHeight)
	ebiten.SetWindowTitle("demo")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
