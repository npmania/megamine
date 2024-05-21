package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	KeyUp = 1 << iota
	KeyDown
	KeyJust
)

type KeyState int

type CursorEvent struct {
	X, Y                int
	Left, Middle, Right KeyState
}

func GetCursorEvent() *CursorEvent {
	ce := &CursorEvent{}
	ce.X, ce.Y = ebiten.CursorPosition()
	ce.Middle = KeyUp
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		ce.Left = (KeyJust | KeyDown)
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		ce.Left = KeyDown
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		ce.Left = (KeyJust | KeyUp)
	} else {
		ce.Left = KeyUp
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		ce.Right = (KeyJust | KeyDown)
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		ce.Right = KeyDown
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
		ce.Right = (KeyJust | KeyUp)
	} else {
		ce.Right = KeyUp
	}
	return ce
}