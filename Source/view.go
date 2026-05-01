package main

import (
	"gopher/gopher"
	"strings"

	"github.com/playdate-go/pdgo"
)

const (
	fontSize   = 20
	marginLeft = 15
)

type GraphicsInterface interface {
	DrawText(text string, x, y int) int
	DrawRect(x, y, width, height int, color pdgo.LCDColor)
	Clear(color pdgo.LCDColor)
}

type CranckInterface interface {
	GetCrankChange() float32
}

type Data struct {
	File  string
	Items []gopher.Item
}

type View struct {
	data     Data
	graphics GraphicsInterface
	system   CranckInterface

	crnk float32

	Cursor gopher.Item
}

func NewView(graphics GraphicsInterface, system CranckInterface) *View {
	return &View{
		graphics: graphics,
		system:   system,

		data: Data{
			Items: []gopher.Item{},
		},
	}
}

func (v *View) Crnk(crnk float32) {
	v.crnk += crnk
}

func (v *View) Data(d Data) {
	v.data = d
}

func (v *View) Render() {
	pos := -v.crnk

	if pos > 0 {
		pos = 0
		v.crnk = 0
	}

	// render text
	v.graphics.Clear(pdgo.SolidWhite)

	// Draw status text
	if len(v.data.Items) == 0 && v.data.File == "" {
		v.graphics.DrawText("Connecting:", 10, 10)

		return
	}

	cursor := gopher.Item{}
	for _, item := range v.data.Items {
		pos += fontSize

		switch item.Type {
		case 'i':
			v.graphics.DrawText(item.Name, marginLeft, int(pos))
			pos += float32(strings.Count(item.Name, "\n")) * fontSize
		case '0':
			if int(pos) >= 40 && int(pos) <= 40+5 {
				cursor = item
			}

			v.graphics.DrawText("TXT |  "+item.Name, marginLeft, int(pos))
		case '1':
			if int(pos) >= 40 && int(pos) <= 40+5 {
				cursor = item
			}

			v.graphics.DrawText("DIR |  "+item.Name, marginLeft, int(pos))
		}
	}

	v.Cursor = cursor

	if v.data.File != "" {
		v.graphics.DrawText(v.data.File, marginLeft, int(pos))
	}

	// pd.System.LogToConsole(v.Cursor.Selector)

	// v.graphics.DrawRect(-1, 40, 10, 15, pdgo.SolidBlack)
	v.graphics.DrawText(">", 1, 40)
}
