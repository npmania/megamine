package game

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	GameReady = iota
	GameActive
	GameWin
	GameDead
)

type GameObject struct {
	X, Y    int
	State   int
	BeginAt time.Time
}

var Game GameObject

func (g *GameObject) Update() error {
	ce := GetCursorEvent()
	switch {
	case Game.State == GameWin || Game.State == GameDead:
		Face.HandleCursorEvent(ce)
	case Game.State == GameActive:
		Clock.TrySet(int(time.Now().Sub(Game.BeginAt).Seconds()))
		Face.HandleCursorEvent(ce)
		_, flagChanged := GameBoard.HandleCursorEvent(ce)
		if flagChanged {
			Counter.Set(GameBoard.Mines - GameBoard.Flags)
		}
	case Game.State == GameReady:
		Face.HandleCursorEvent(ce)
		boardChanged, flagChanged := GameBoard.HandleCursorEvent(ce)
		if (boardChanged || flagChanged) && Game.State != GameDead {
			g.State = GameActive
			Game.BeginAt = time.Now()
			Counter.Set(GameBoard.Mines - GameBoard.Flags)
		}
	}
	return nil
}

func (g *GameObject) Draw(screen *ebiten.Image) {
	DrawScreen(screen)
}

func (g *GameObject) Layout(outw, outh int) (w, h int) {
	return g.X, g.Y
}

func InitGame() error {

	err := ImportGameImages()
	if err != nil {
		return err
	}
	ebiten.SetWindowSize(756, 482)
	ebiten.SetWindowTitle("MegaMine!")
	ebiten.SetWindowResizable(true)

	Game = GameObject{
		X:     500 * 1,
		Y:     318 * 1,
		State: GameReady,
	}
	err = InitBoard()
	if err != nil {
		return err
	}
	InitSegDisp()
	Clock.Set(0)
	Counter.Set(GameBoard.Mines)
	UpdatePos()
	return nil
}

func ResetGame() error {
	board, err := NewBoard(30*1, 16*1, 99*1)
	if err != nil {
		return err
	}
	board.Generate()
	GameBoard = board
	Game.State = GameReady
	UpdatePos()
	Clock.Set(0)
	Counter.Set(GameBoard.Mines)
	return nil
}

func UpdatePos() {
	GameBoard.UpdatePos()
	Face.UpdatePos()
	UpdatePosSegDisp()
}
