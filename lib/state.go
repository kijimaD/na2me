// 参考: https://github.com/x-hgg-x/goecsengine/blob/master/states/lib.go
package lib

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type TransType int

const (
	TransNone TransType = iota
	TransSwitch
)

type Transition struct {
	Type      TransType
	NewStates []State
}

type State interface {
	// Executed when the state begins
	OnStart()
	// Executed when the state exits
	OnStop()
	// Executed when a new state is pushed over this one
	OnPause()
	// Executed when the state become active again (states pushed over this one have been popped)
	OnResume()
	// Executed on every frame when the state is active
	Update() Transition
	// 描画
	Draw(screen *ebiten.Image)
}

type StateMachine struct {
	states         []State
	lastTransition Transition
}

// Init creates a new state machine with an initial state
func Init(s State) StateMachine {
	s.OnStart()
	return StateMachine{[]State{s}, Transition{TransNone, []State{}}}
}

func (sm *StateMachine) Update() {
	switch sm.lastTransition.Type {
	case TransSwitch:
		sm._Switch(sm.lastTransition.NewStates)
	}

	if len(sm.states) < 1 {
		os.Exit(0)
	}

	// Run state update function with game systems
	sm.lastTransition = sm.states[len(sm.states)-1].Update()
}

// Draw draws the screen after a state update
func (sm *StateMachine) Draw(screen *ebiten.Image) {
	sm.states[len(sm.states)-1].Draw(screen)
}

// Remove the active state and replace it by a new one
func (sm *StateMachine) _Switch(newStates []State) {
	if len(newStates) != 1 {
		log.Fatal("switch transition accept only one new state")
	}

	sm.states[len(sm.states)-1].OnStop()
	newStates[0].OnStart()
	sm.states[len(sm.states)-1] = newStates[0]
}
