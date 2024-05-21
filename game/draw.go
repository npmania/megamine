package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var bgColor = color.RGBA{
	R: 0xc0,
	G: 0xc0,
	B: 0xc0,
	A: 0xff,
}

func DrawBoard(s *ebiten.Image) {
	img := GameBoard.Image()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(GameBoard.Pos.X), float64(GameBoard.Pos.Y))
	s.DrawImage(img, op)
}

func DrawFace(s *ebiten.Image) {
	img := Face.Current()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(Face.Pos.X), float64(Face.Pos.Y))
	s.DrawImage(img, op)
}

func DrawSegDisp(s *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(Counter.Pos.X), float64(Counter.Pos.Y))
	s.DrawImage(Counter.img, op)
	op.GeoM.Reset()
	op.GeoM.Translate(float64(Clock.Pos.X), float64(Clock.Pos.Y))
	s.DrawImage(Clock.img, op)
}

func DrawScreen(s *ebiten.Image) {
	s.Fill(bgColor)
	DrawBoard(s)
	DrawFace(s)
	DrawSegDisp(s)
}
