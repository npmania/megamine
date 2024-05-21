package game

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"image"
	_ "image/png"
	"io"
	"io/ioutil"

	"github.com/hajimehoshi/ebiten/v2"
)

type SpriteInfo []struct {
	Name string
	Rect image.Rectangle
}

const (
	ImgCell1 = iota
	ImgCell2
	ImgCell3
	ImgCell4
	ImgCell5
	ImgCell6
	ImgCell7
	ImgCell8
	ImgUnopened
	ImgOpened
	ImgFlagged
	ImgGuess
	ImgOpenedGuess
	ImgOpenedMined
	ImgOpenedExploded
	ImgWrongFlag
)

const (
	ImgSmile = iota
	ImgSmilePress
	ImgOops
	ImgSunglass
	ImgDead
)

const (
	ImgNum1 = iota
	ImgNum2
	ImgNum3
	ImgNum4
	ImgNum5
	ImgNum6
	ImgNum7
	ImgNum8
	ImgNum9
	ImgNum0
	ImgNumHyphen
	ImgNumEmpty
)

type GameImages struct {
	Cell [16]*ebiten.Image
	Face [5]*ebiten.Image
	Num  [12]*ebiten.Image
	Full *ebiten.Image
}

func (si SpriteInfo) getRect(name string) image.Rectangle {
	for _, s := range si {
		if name == s.Name {
			return s.Rect
		}
	}
	return image.Rectangle{}
}

func (si *SpriteInfo) loadSubImage(str string) *ebiten.Image {
	return Ass.Images.Full.SubImage(si.getRect(str)).(*ebiten.Image)
}

//go:embed ass/ms.png
var msPng []byte

func msPNG() io.Reader { return bytes.NewReader(msPng) }

//go:embed ass/ms.json
var msJson []byte

func msJSON() io.Reader { return bytes.NewReader(msJson) }

func ImportGameImages() error {
	img, _, err := image.Decode(msPNG())
	if err != nil {
		return err
	}

	gi := &Ass.Images
	gi.Full = ebiten.NewImageFromImage(img)

	data, _ := ioutil.ReadAll(msJSON())

	si := SpriteInfo{}
	err = json.Unmarshal(data, &si)
	if err != nil {
		return err
	}

	gi.Cell[ImgCell1] = si.loadSubImage("cell1")
	gi.Cell[ImgCell2] = si.loadSubImage("cell2")
	gi.Cell[ImgCell3] = si.loadSubImage("cell3")
	gi.Cell[ImgCell4] = si.loadSubImage("cell4")
	gi.Cell[ImgCell5] = si.loadSubImage("cell5")
	gi.Cell[ImgCell6] = si.loadSubImage("cell6")
	gi.Cell[ImgCell7] = si.loadSubImage("cell7")
	gi.Cell[ImgCell8] = si.loadSubImage("cell8")

	gi.Cell[ImgUnopened] = si.loadSubImage("unopened")
	gi.Cell[ImgOpened] = si.loadSubImage("opened")
	gi.Cell[ImgFlagged] = si.loadSubImage("flagged")
	gi.Cell[ImgGuess] = si.loadSubImage("guess")
	gi.Cell[ImgOpenedGuess] = si.loadSubImage("openedguess")
	gi.Cell[ImgOpenedMined] = si.loadSubImage("openedmined")
	gi.Cell[ImgOpenedExploded] = si.loadSubImage("openedexploded")
	gi.Cell[ImgWrongFlag] = si.loadSubImage("wrongflag")

	gi.Face[ImgSmile] = si.loadSubImage("smile")
	gi.Face[ImgSmilePress] = si.loadSubImage("smilepress")
	gi.Face[ImgOops] = si.loadSubImage("oops")
	gi.Face[ImgSunglass] = si.loadSubImage("sunglass")
	gi.Face[ImgDead] = si.loadSubImage("dead")

	gi.Num[ImgNum1] = si.loadSubImage("num1")
	gi.Num[ImgNum2] = si.loadSubImage("num2")
	gi.Num[ImgNum3] = si.loadSubImage("num3")
	gi.Num[ImgNum4] = si.loadSubImage("num4")
	gi.Num[ImgNum5] = si.loadSubImage("num5")
	gi.Num[ImgNum6] = si.loadSubImage("num6")
	gi.Num[ImgNum7] = si.loadSubImage("num7")
	gi.Num[ImgNum8] = si.loadSubImage("num8")
	gi.Num[ImgNum9] = si.loadSubImage("num9")
	gi.Num[ImgNum0] = si.loadSubImage("num0")
	gi.Num[ImgNumHyphen] = si.loadSubImage("numhyphen")
	gi.Num[ImgNumEmpty] = si.loadSubImage("numempty")

	return nil
}
