package game

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	segNum1 = iota
	segNum2
	segNum3
	segNum4
	segNum5
	segNum6
	segNum7
	segNum8
	segNum9
	segNum0
	segMinus
	segEmpty
)

const (
	digitWidth  = 13
	digitHeight = 23
)

type segDigit struct {
	img   *ebiten.Image
	value int
}

type SegDisp struct {
	disp  [3]segDigit
	img   *ebiten.Image
	value int
	Pos   image.Point
}

var Counter SegDisp
var Clock SegDisp

func UpdatePosSegDisp() {
	gw, _ := Game.Layout(ebiten.WindowSize())
	Counter.Pos.Y, Clock.Pos.Y = Face.Pos.Y, Face.Pos.Y
	Counter.Pos.X = Counter.Pos.Y
	Clock.Pos.X = gw - Clock.Pos.Y - digitWidth*3
}

// digit over 9 defaults to 9, negative value represents hyphen
func (ss *segDigit) Set(v int) {
	switch v {
	case 1:
		ss.img = Ass.Images.Num[ImgNum1]
	case 2:
		ss.img = Ass.Images.Num[ImgNum2]
	case 3:
		ss.img = Ass.Images.Num[ImgNum3]
	case 4:
		ss.img = Ass.Images.Num[ImgNum4]
	case 5:
		ss.img = Ass.Images.Num[ImgNum5]
	case 6:
		ss.img = Ass.Images.Num[ImgNum6]
	case 7:
		ss.img = Ass.Images.Num[ImgNum7]
	case 8:
		ss.img = Ass.Images.Num[ImgNum8]
	case 9:
		ss.img = Ass.Images.Num[ImgNum9]
	case 0:
		ss.img = Ass.Images.Num[ImgNum0]
	default:
		if v > 9 {
			defer ss.Set(9)
		} else {
			ss.img = Ass.Images.Num[ImgNumHyphen]
		}
	}
}

func (ss *segDigit) SetHyphen() {
	ss.img = Ass.Images.Num[ImgNumHyphen]
}

func (ss *segDigit) SetEmpty() {
	ss.img = Ass.Images.Num[ImgNumEmpty]
}

func (sd *SegDisp) Get() int {
	return sd.value
}

func (sd *SegDisp) TrySet(v int) {
	if v == sd.value {
		return
	}
	sd.Set(v)
}
func (sd *SegDisp) Set(v int) {
	defer sd.RenderAll()
	sd.value = v
	n1, n2, n3 := v/100, v/10, v%10
	if v < 0 {
		sd.disp[0].SetHyphen()
	} else if n1 > 9 {
		sd.disp[0].Set(9)
		sd.disp[1].Set(9)
		sd.disp[2].Set(9)
		return
	} else {
		sd.disp[0].Set(n1 % 10)
	}
	if n2 < -9 {
		sd.disp[1].Set(9)
		sd.disp[2].Set(9)
		return
	} else if n2 < 0 {
		sd.disp[1].Set(-n2 % 10)
	} else {
		sd.disp[1].Set(n2 % 10)
	}
	if n3 < 0 {
		sd.disp[2].Set(-n3)
	} else {
		sd.disp[2].Set(n3)
	}
}

func (sd *SegDisp) RenderAll() {
	for i := 0; i < 3; i++ {
		sd.Render(i)
	}
}

func (sd *SegDisp) Render(i int) {
	op := ebiten.DrawImageOptions{}
	switch i {
	case 0:
	case 1:
		op.GeoM.Translate(digitWidth, 0)
	case 2:
		op.GeoM.Translate(digitWidth*2, 0)
	default:
		return
	}
	sd.img.DrawImage(sd.disp[i].img, &op)
}

func InitSegDisp() {
	Clock.img = ebiten.NewImage(digitWidth*3, digitHeight)
	Counter.img = ebiten.NewImage(digitWidth*3, digitHeight)
}
