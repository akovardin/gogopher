package main

import (
	"testing"

	"gopher/gopher"

	"github.com/playdate-go/pdgo"
)

type mockGraphics struct {
	drawTextFunc func(text string, x, y int) int
	drawRectFunc func(x, y, width, height int, color pdgo.LCDColor)
	clearFunc    func(color pdgo.LCDColor)
}

func (m *mockGraphics) DrawText(text string, x, y int) int {
	return m.drawTextFunc(text, x, y)
}

func (m *mockGraphics) DrawRect(x, y, width, height int, color pdgo.LCDColor) {
	if m.drawRectFunc != nil {
		m.drawRectFunc(x, y, width, height, color)
	}
}

func (m *mockGraphics) Clear(color pdgo.LCDColor) {
	m.clearFunc(color)
}

type mockCrank struct {
	getCrankChangeFunc func() float32
}

func (m *mockCrank) GetCrankChange() float32 {
	return m.getCrankChangeFunc()
}

func TestNewView(t *testing.T) {
	g := &mockGraphics{
		drawTextFunc: func(text string, x, y int) int { return 0 },
		clearFunc:    func(color pdgo.LCDColor) {},
	}
	s := &mockCrank{
		getCrankChangeFunc: func() float32 { return 0 },
	}

	v := NewView(g, s)
	if v == nil {
		t.Fatal("NewView returned nil")
	}
	if len(v.data.Items) != 0 {
		t.Errorf("expected 0 default items, got %d", len(v.data.Items))
	}
	if v.crnk != 0 {
		t.Errorf("expected crnk to be 0, got %f", v.crnk)
	}
}

func TestViewData(t *testing.T) {
	v := &View{}
	items := []gopher.Item{
		{Type: '0', Name: "File1"},
	}
	v.Data(Data{Items: items})
	if len(v.data.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(v.data.Items))
	}
	if v.data.Items[0].Name != "File1" {
		t.Errorf("expected File1, got %s", v.data.Items[0].Name)
	}
}

func TestRender_EmptyItems(t *testing.T) {
	var cleared bool
	var drawnText string
	var drawnX, drawnY int

	g := &mockGraphics{
		drawTextFunc: func(text string, x, y int) int {
			drawnText = text
			drawnX = x
			drawnY = y
			return 0
		},
		clearFunc: func(color pdgo.LCDColor) {
			cleared = true
			if color != pdgo.SolidWhite {
				t.Errorf("expected SolidWhite, got %v", color)
			}
		},
	}
	s := &mockCrank{
		getCrankChangeFunc: func() float32 { return 0 },
	}

	v := &View{graphics: g, system: s}
	v.Render()

	if !cleared {
		t.Error("expected Clear to be called")
	}
	if drawnText != "Connecting:" {
		t.Errorf("expected 'Connecting:', got %q", drawnText)
	}
	if drawnX != 10 || drawnY != 10 {
		t.Errorf("expected (10,10), got (%d,%d)", drawnX, drawnY)
	}
}

func TestRender_InfoItem(t *testing.T) {
	var draws []struct {
		text string
		x, y int
	}

	g := &mockGraphics{
		drawTextFunc: func(text string, x, y int) int {
			draws = append(draws, struct {
				text string
				x, y int
			}{text, x, y})
			return 0
		},
		clearFunc: func(color pdgo.LCDColor) {},
	}
	s := &mockCrank{
		getCrankChangeFunc: func() float32 { return 0 },
	}

	v := &View{
		graphics: g,
		system:   s,
		data: Data{Items: []gopher.Item{
			{Type: 'i', Name: "info line"},
		}},
	}
	v.Render()

	if len(draws) != 2 {
		t.Fatalf("expected 2 DrawText calls, got %d", len(draws))
	}
	if draws[0].text != "info line" {
		t.Errorf("expected 'info line', got %q", draws[0].text)
	}
	if draws[0].x != 15 || draws[0].y != 20 {
		t.Errorf("expected (15,20), got (%d,%d)", draws[0].x, draws[0].y)
	}
	if draws[1].text != ">" {
		t.Errorf("expected '>', got %q", draws[1].text)
	}
	if draws[1].x != 1 || draws[1].y != 40 {
		t.Errorf("expected (1,40), got (%d,%d)", draws[1].x, draws[1].y)
	}
}

func TestRender_TextItem(t *testing.T) {
	var draws []struct {
		text string
		x, y int
	}

	g := &mockGraphics{
		drawTextFunc: func(text string, x, y int) int {
			draws = append(draws, struct {
				text string
				x, y int
			}{text, x, y})
			return 0
		},
		clearFunc: func(color pdgo.LCDColor) {},
	}
	s := &mockCrank{
		getCrankChangeFunc: func() float32 { return 0 },
	}

	v := &View{
		graphics: g,
		system:   s,
		data: Data{Items: []gopher.Item{
			{Type: '0', Name: "readme.txt"},
		}},
	}
	v.Render()

	if len(draws) != 2 {
		t.Fatalf("expected 2 DrawText calls, got %d", len(draws))
	}
	if draws[0].text != "TXT |  readme.txt" {
		t.Errorf("expected 'TXT |  readme.txt', got %q", draws[0].text)
	}
	if draws[1].text != ">" {
		t.Errorf("expected '>', got %q", draws[1].text)
	}
}

func TestRender_DirectoryItem(t *testing.T) {
	var draws []struct {
		text string
		x, y int
	}

	g := &mockGraphics{
		drawTextFunc: func(text string, x, y int) int {
			draws = append(draws, struct {
				text string
				x, y int
			}{text, x, y})
			return 0
		},
		clearFunc: func(color pdgo.LCDColor) {},
	}
	s := &mockCrank{
		getCrankChangeFunc: func() float32 { return 0 },
	}

	v := &View{
		graphics: g,
		system:   s,
		data: Data{Items: []gopher.Item{
			{Type: '1', Name: "documents"},
		}},
	}
	v.Render()

	if len(draws) != 2 {
		t.Fatalf("expected 2 DrawText calls, got %d", len(draws))
	}
	if draws[0].text != "DIR |  documents" {
		t.Errorf("expected 'DIR |  documents', got %q", draws[0].text)
	}
}

func TestRender_MixedItems(t *testing.T) {
	var draws []struct {
		text string
		x, y int
	}

	g := &mockGraphics{
		drawTextFunc: func(text string, x, y int) int {
			draws = append(draws, struct {
				text string
				x, y int
			}{text, x, y})
			return 0
		},
		clearFunc: func(color pdgo.LCDColor) {},
	}
	s := &mockCrank{
		getCrankChangeFunc: func() float32 { return 0 },
	}

	v := &View{
		graphics: g,
		system:   s,
		data: Data{Items: []gopher.Item{
			{Type: 'i', Name: "header"},
			{Type: '0', Name: "file.txt"},
			{Type: '1', Name: "subdir"},
		}},
	}
	v.Render()

	if len(draws) != 4 {
		t.Fatalf("expected 4 DrawText calls, got %d", len(draws))
	}

	if draws[0].text != "header" {
		t.Errorf("draw[0]: expected 'header', got %q", draws[0].text)
	}
	if draws[0].x != 15 || draws[0].y != 20 {
		t.Errorf("draw[0]: expected (15,20), got (%d,%d)", draws[0].x, draws[0].y)
	}
	if draws[1].text != "TXT |  file.txt" {
		t.Errorf("draw[1]: expected 'TXT |  file.txt', got %q", draws[1].text)
	}
	if draws[1].x != 15 || draws[1].y != 40 {
		t.Errorf("draw[1]: expected (15,40), got (%d,%d)", draws[1].x, draws[1].y)
	}
	if draws[2].text != "DIR |  subdir" {
		t.Errorf("draw[2]: expected 'DIR |  subdir', got %q", draws[2].text)
	}
	if draws[2].x != 15 || draws[2].y != 60 {
		t.Errorf("draw[2]: expected (15,60), got (%d,%d)", draws[2].x, draws[2].y)
	}
	if draws[3].text != ">" {
		t.Errorf("draw[3]: expected '>', got %q", draws[3].text)
	}
}

func TestRender_CrankScrolling(t *testing.T) {
	var draws []struct {
		text string
		x, y int
	}

	g := &mockGraphics{
		drawTextFunc: func(text string, x, y int) int {
			draws = append(draws, struct {
				text string
				x, y int
			}{text, x, y})
			return 0
		},
		clearFunc: func(color pdgo.LCDColor) {},
	}
	s := &mockCrank{
		getCrankChangeFunc: func() float32 { return 0 },
	}

	v := &View{
		graphics: g,
		system:   s,
		data: Data{Items: []gopher.Item{
			{Type: '0', Name: "file.txt"},
		}},
	}
	v.Render()

	if len(draws) != 2 {
		t.Fatalf("expected 2 DrawText calls, got %d", len(draws))
	}
	if draws[0].y != 20 {
		t.Errorf("expected y=20, got %d", draws[0].y)
	}
}

func TestRender_CrankAccumulates(t *testing.T) {
	crankCalls := 0
	g := &mockGraphics{
		drawTextFunc: func(text string, x, y int) int { return 0 },
		clearFunc:    func(color pdgo.LCDColor) {},
	}
	s := &mockCrank{
		getCrankChangeFunc: func() float32 {
			crankCalls++
			return 10
		},
	}

	v := &View{
		graphics: g,
		system:   s,
		data: Data{Items: []gopher.Item{
			{Type: '0', Name: "a"},
			{Type: '0', Name: "b"},
		}},
	}

	v.Crnk(s.GetCrankChange())
	if v.crnk != 10 {
		t.Errorf("expected crnk=10 after first Crnk, got %f", v.crnk)
	}

	v.Crnk(s.GetCrankChange())
	if v.crnk != 20 {
		t.Errorf("expected crnk=20 after second Crnk, got %f", v.crnk)
	}

	if crankCalls != 2 {
		t.Errorf("expected 2 GetCrankChange calls, got %d", crankCalls)
	}
}

func TestRender_UnhandledType(t *testing.T) {
	var draws []struct {
		text string
		x, y int
	}

	g := &mockGraphics{
		drawTextFunc: func(text string, x, y int) int {
			draws = append(draws, struct {
				text string
				x, y int
			}{text, x, y})
			return 0
		},
		clearFunc: func(color pdgo.LCDColor) {},
	}
	s := &mockCrank{
		getCrankChangeFunc: func() float32 { return 0 },
	}

	v := &View{
		graphics: g,
		system:   s,
		data: Data{Items: []gopher.Item{
			{Type: '7', Name: "search"},
			{Type: '9', Name: "binary"},
		}},
	}
	v.Render()

	if len(draws) != 1 {
		t.Fatalf("expected 1 DrawText call (for cursor), got %d", len(draws))
	}
	if draws[0].text != ">" {
		t.Errorf("expected '>', got %q", draws[0].text)
	}
}
