package game

import (
	"errors"
	"fmt"
	"image"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type CellState int

const (
	CellMine = 1 << iota
	CellOpen
	CellFlag
	CellGuess
)

type XrayKind int

const (
	XrayOff = iota
	XrayNarrow
	XrayWide
	XrayDisabled
)

type Cell struct {
	State  CellState
	Nearby int
}

type Board struct {
	X, Y         int
	Mines        int
	Board        [][]Cell
	Pos          image.Point
	CellsLeft    int
	Flags        int
	XrayX, XrayY int
	XrayMode     XrayKind
	img          *ebiten.Image
}

var GameBoard *Board

func NewBoard(x, y, mines int) (*Board, error) {
	if x < 9 || y < 9 {
		return nil, errors.New(fmt.Sprintf("invalid size: %ux%u", x, y))
	} else if mines < 10 {
		return nil, errors.New(fmt.Sprintf("Too few mines"))
	}
	size := x * y
	if y != 0 && size/y != x {
		return nil, errors.New("Board size too large")
	}
	b := &Board{
		X:         x,
		Y:         y,
		Mines:     mines,
		Flags:     0,
		CellsLeft: x*y - mines,
	}

	b.Generate()
	b.img = ebiten.NewImage(b.X*16, b.Y*16)
	b.renderAll()

	return b, nil
}

func (b *Board) Generate() {
	rand.Seed(time.Now().UnixNano())

	b.Board = make([][]Cell, b.Y)
	for by := range b.Board {
		b.Board[by] = make([]Cell, b.X)
	}

	for i := 0; i < b.Mines; i++ {
		for {
			my := rand.Intn(b.Y)
			mx := rand.Intn(b.X)
			if b.Board[my][mx].State&CellMine != 0 {
				continue
			}
			b.Board[my][mx].State |= CellMine
			for yy := my - 1; yy <= my+1; yy++ {
				for xx := mx - 1; xx <= mx+1; xx++ {
					if xx < 0 || xx >= b.X || yy < 0 || yy >= b.Y {
						continue
					}
					b.Board[yy][xx].Nearby++
				}
			}
			break
		}
	}
}

func InitBoard() (err error) {
	GameBoard, err = NewBoard(30*1, 16*1, 99*1)
	if err != nil {
		return
	}
	return nil
}

func (b *Board) UpdatePos() {
	gw, gh := Game.Layout(ebiten.WindowSize())
	iw, ih := 16*b.X, 16*b.Y
	b.Pos.X = (gw - iw) / 2
	b.Pos.Y = (gh - ih) - b.Pos.X
}

func (b *Board) recursiveOpenCell(x, y int) {
	for yy := y - 1; yy <= y+1; yy++ {
		for xx := x - 1; xx <= x+1; xx++ {
			if xx < 0 || yy < 0 || xx >= b.X || yy >= b.Y {
				continue
			} else if xx == x && yy == y {
				continue
			}
			c := &b.Board[yy][xx]
			if c.State&(CellOpen|CellFlag|CellGuess) != 0 {
				continue
			}
			c.State |= CellOpen
			b.CellsLeft--
			if c.Nearby == 0 {
				b.recursiveOpenCell(xx, yy)
			}
		}
	}
}

func (b *Board) addMineOnLeftmost() {
	for y := 0; y < b.Y; y++ {
		for x := 0; x < b.X; x++ {
			c := &b.Board[y][x]
			if c.State&CellMine != 0 {
				continue
			}
			c.State |= CellMine
			for yy := y - 1; yy <= y+1; yy++ {
				for xx := x - 1; xx <= x+1; xx++ {
					if xx < 0 || yy < 0 || xx >= b.X || yy >= b.Y {
						continue
					}
					b.Board[yy][xx].Nearby++
				}
			}
			return
		}
	}
}

func (b *Board) moveMineToLeftmost(x, y int) {
	b.Board[y][x].State ^= CellMine
	for yy := y - 1; yy <= y+1; yy++ {
		for xx := x - 1; xx <= x+1; xx++ {
			if xx < 0 || yy < 0 || xx >= b.X || yy >= b.Y {
				continue
			}
			b.Board[yy][xx].Nearby--
		}
	}
	b.addMineOnLeftmost()
}

func (b *Board) openCell(x, y int) {
	c := &b.Board[y][x]
	if c.State&(CellFlag|CellGuess|CellOpen) != 0 {
		return
	}
	c.State |= CellOpen
	b.CellsLeft--
	if c.State&CellMine != 0 {
		if Game.State == GameReady {
			b.moveMineToLeftmost(x, y)
		} else {
			Game.State = GameDead
			b.renderAll()
			return
		}
	}
	if c.Nearby == 0 {
		b.recursiveOpenCell(x, y)
	}
	if b.CellsLeft == 0 {
		Counter.Set(0)
		Game.State = GameWin
	}
	return
}

func (b *Board) tryOpenCell(x, y int) bool {
	if b.XrayMode == XrayDisabled {
		return false
	}
	b.stopXray()
	b.openCell(x, y)
	b.renderMask(CellOpen)
	return true
}

func (b *Board) xrayNarrowCell(x, y int) {
	if b.XrayMode == XrayDisabled {
		return
	}
	b.unrenderXray()
	b.XrayX, b.XrayY = x, y
	b.XrayMode = XrayNarrow
	b.renderXray()
}

func (b *Board) tryChording(x, y int) bool {
	st := b.Board[y][x].State
	if st&CellOpen == 0 {
		return false
	}
	cnt := 0
	for yy := y - 1; yy <= y+1; yy++ {
		for xx := x - 1; xx <= x+1; xx++ {
			if xx < 0 || xx >= b.X || yy < 0 || yy >= b.Y {
				continue
			}
			st := b.Board[yy][xx].State
			if st&CellFlag != 0 {
				cnt++
			}
		}
	}
	if cnt != b.Board[y][x].Nearby {
		return false
	}
	for yy := y - 1; yy <= y+1; yy++ {
		for xx := x - 1; xx <= x+1; xx++ {
			if xx < 0 || xx >= b.X || yy < 0 || yy >= b.Y {
				continue
			} else if xx == x && yy == y {
				continue
			}
			st := b.Board[yy][xx].State
			if st&(CellFlag|CellOpen) != 0 {
				continue
			}
			b.openCell(xx, yy)
		}
	}
	b.renderMask(CellOpen | CellFlag | CellGuess)
	return true
}

func (b *Board) xrayCell(x, y int) {
	b.unrenderXray()
	b.XrayMode = XrayWide
	b.XrayX, b.XrayY = x, y
	b.renderXray()
}

func (b *Board) flagCell(x, y int) {
	cell := &b.Board[y][x]
	if cell.State&CellOpen != 0 {
		return
	}
	switch {
	case cell.State&CellGuess != 0:
		cell.State ^= CellGuess
	case cell.State&CellFlag != 0:
		cell.State ^= CellFlag
		cell.State |= CellGuess
		b.Flags--
	default:
		cell.State ^= CellFlag
		b.Flags++
	}
	b.renderCell(x, y)
}

func (b *Board) disableXray() {
	b.unrenderXray()
	b.XrayMode = XrayDisabled
}

func (b *Board) stopXray() {
	b.unrenderXray()
	b.XrayMode = XrayOff
}

func (b *Board) cursorCell(ce *CursorEvent) (int, int, bool) {
	dx, dy := ce.X-b.Pos.X, ce.Y-b.Pos.Y
	if dx >= b.X*16 || dy >= b.Y*16 || dx < 0 || dy < 0 {
		return 0, 0, false
	}
	return dx / 16, dy / 16, true
}

func (b *Board) HandleCursorEvent(ce *CursorEvent) (cellChanged, flagChanged bool) {
	flagChanged = false
	x, y, ok := b.cursorCell(ce)
	switch {
	case !ok:
		b.stopXray()
	case ce.Left&KeyDown != 0 && ce.Right&KeyDown != 0:
		cellChanged = b.tryChording(x, y)
		if !cellChanged {
			b.xrayCell(x, y)
		}
	case ce.Left&KeyDown != 0 && ce.Right == KeyJust|KeyUp:
		b.disableXray()
	case ce.Right&KeyDown != 0 && ce.Left == KeyJust|KeyUp:
		b.disableXray()
	case ce.Left&KeyDown != 0:
		b.xrayNarrowCell(x, y)
	case ce.Right == KeyJust|KeyDown:
		b.flagCell(x, y)
		flagChanged = true
	case ce.Left == KeyJust|KeyUp && ce.Right != KeyJust|KeyUp:
		cellChanged = b.tryOpenCell(x, y)
	default:
		b.stopXray()
	}
	return
}

func openedCellImage(c Cell) *ebiten.Image {
	switch c.Nearby {
	case 0:
		return Ass.Images.Cell[ImgOpened]
	case 1:
		return Ass.Images.Cell[ImgCell1]
	case 2:
		return Ass.Images.Cell[ImgCell2]
	case 3:
		return Ass.Images.Cell[ImgCell3]
	case 4:
		return Ass.Images.Cell[ImgCell4]
	case 5:
		return Ass.Images.Cell[ImgCell5]
	case 6:
		return Ass.Images.Cell[ImgCell6]
	case 7:
		return Ass.Images.Cell[ImgCell7]
	case 8:
		return Ass.Images.Cell[ImgCell8]
	}
	return nil
}

func matchCellImage(c Cell) *ebiten.Image {
	switch {
	case c.State&CellOpen != 0 && c.State&CellMine != 0:
		return Ass.Images.Cell[ImgOpenedExploded]
	case c.State&CellOpen != 0:
		return openedCellImage(c)
	case c.State&CellGuess != 0:
		return Ass.Images.Cell[ImgGuess]
	case Game.State == GameDead && c.State&CellFlag != 0 && c.State&CellMine == 0:
		return Ass.Images.Cell[ImgWrongFlag]
	case c.State&CellFlag != 0:
		return Ass.Images.Cell[ImgFlagged]
	case Game.State == GameDead && c.State&CellOpen == 0 && c.State&CellMine != 0:
		return Ass.Images.Cell[ImgOpenedMined]
	default:
		return Ass.Images.Cell[ImgUnopened]
	}
}

func (b *Board) renderCell(x, y int) {
	if y < 0 || y >= b.Y || x < 0 || x > b.X {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x*16), float64(y*16))
	b.img.DrawImage(matchCellImage(b.Board[y][x]), op)
}

// Utility function to use while rendering xray
func (b *Board) renderOpenCell(x, y int) {
	if y < 0 || y >= b.Y || x < 0 || x > b.X {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x*16), float64(y*16))
	b.img.DrawImage(Ass.Images.Cell[ImgOpened], op)
}

func (b *Board) unrenderXray() {
	if b.XrayMode == XrayOff || b.XrayMode == XrayDisabled {
		return
	}
	if b.XrayMode == XrayNarrow {
		b.renderCell(b.XrayX, b.XrayY)
	} else if b.XrayMode == XrayWide {
		for y := b.XrayY - 1; y <= b.XrayY+1; y++ {
			for x := b.XrayX - 1; x <= b.XrayX+1; x++ {
				if x < 0 || x >= b.X || y < 0 || y > b.Y {
					continue
				}
				b.renderCell(x, y)
			}
		}
	}
}

func (b *Board) renderXray() {
	if b.XrayMode == XrayOff || b.XrayMode == XrayDisabled {
		return
	}
	if b.XrayMode == XrayNarrow {
		if b.Board[b.XrayY][b.XrayX].State&(CellOpen|CellFlag|CellGuess) != 0 {
			return
		}
		b.renderOpenCell(b.XrayX, b.XrayY)
	} else if b.XrayMode == XrayWide {
		for yy := b.XrayY - 1; yy <= b.XrayY+1; yy++ {
			for xx := b.XrayX - 1; xx <= b.XrayX+1; xx++ {
				if xx < 0 || xx >= b.X || yy < 0 || yy >= b.Y {
					continue
				}
				if b.Board[yy][xx].State&(CellOpen|CellFlag|CellGuess) != 0 {
					continue
				}
				b.renderOpenCell(xx, yy)
			}
		}
	}
}
func (b *Board) renderMask(st CellState) {
	for y := 0; y < b.Y; y++ {
		for x := 0; x < b.X; x++ {
			if b.Board[y][x].State&st == 0 {
				continue
			}
			b.renderCell(x, y)
		}
	}
}

func (b *Board) renderAll() {
	for y := 0; y < b.Y; y++ {
		for x := 0; x < b.X; x++ {
			b.renderCell(x, y)
		}
	}
}

func (b *Board) Image() *ebiten.Image {
	return b.img
}
