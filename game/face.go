package game

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	FaceWidth  = 24
	FaceHeight = 24
)

type FaceObject struct {
	Pos     image.Point
	Clicked bool
}

var Face FaceObject

func (f *FaceObject) Current() *ebiten.Image {
	switch {
	case f.Clicked:
		return Ass.Images.Face[ImgSmilePress]
	case Game.State == GameDead:
		return Ass.Images.Face[ImgDead]
	case Game.State == GameWin:
		return Ass.Images.Face[ImgSunglass]
	case GameBoard.XrayMode != XrayOff && GameBoard.XrayMode != XrayDisabled:
		return Ass.Images.Face[ImgOops]
	default:
		return Ass.Images.Face[ImgSmile]
	}
}

func (f *FaceObject) UpdatePos() {
	gw, _ := Game.Layout(ebiten.WindowSize())
	fw, fh := FaceWidth, FaceHeight
	f.Pos.X = (gw - fw) / 2
	f.Pos.Y = (GameBoard.Pos.Y - fh) / 2
}

func InitFace() {
	Face = FaceObject{
		Clicked: false,
	}
}

func (f *FaceObject) underCursor(cx, cy int) bool {
	if cx < f.Pos.X || cx >= f.Pos.X+FaceWidth ||
		cy < f.Pos.Y || cy >= f.Pos.Y+FaceHeight {
		return false
	}
	return true
}

func (f *FaceObject) HandleCursorEvent(ce *CursorEvent) {
	switch {
	case ce.Left&(KeyJust|KeyDown) != 0:
		f.HandleLeftClick(ce)
	case ce.Left&(KeyJust|KeyUp) != 0:
		f.HandleLeftClickUp(ce)
	}
}

func (f *FaceObject) HandleLeftClick(ce *CursorEvent) {
	if f.underCursor(ce.X, ce.Y) {
		f.Clicked = true
	}
}

func (f *FaceObject) HandleLeftClickUp(ce *CursorEvent) {
	if !f.Clicked {
		return
	}
	f.Clicked = false
	if f.underCursor(ce.X, ce.Y) {
		Game.State = GameReady
		ResetGame()
	}
}
