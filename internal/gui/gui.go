package gui

import (
	"image/color"

	"gobrowser/internal/page"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

const (
	width, height = 800, 600
	hStep, vStep  = 13, 18
)

type displayItem struct {
	X, Y float32
	Char rune
}

type sizedLayout struct{ size fyne.Size }

func (l *sizedLayout) Layout([]fyne.CanvasObject, fyne.Size) {}
func (l *sizedLayout) MinSize([]fyne.CanvasObject) fyne.Size { return l.size }

type Browser struct {
	app    fyne.App
	window fyne.Window
}

func NewBrowser() *Browser {
	a := app.New()
	w := a.NewWindow("Browser")
	w.Resize(fyne.NewSize(width, height))
	return &Browser{app: a, window: w}
}

func (b *Browser) Render(tokens []page.Token) {
	items := b.layout(tokens)
	scroll := b.draw(items)
	b.window.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		switch key.Name {
		case fyne.KeyUp:
			scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(0, vStep)})
		case fyne.KeyDown:
			scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(0, -vStep)})
		}
	})
	b.window.ShowAndRun()
}

func (b *Browser) layout(tokens []page.Token) []displayItem {
	var items []displayItem
	cursorX, cursorY := float32(hStep), float32(vStep)
	for _, token := range tokens {
		text, ok := token.(page.Text)
		if !ok {
			continue
		}
		for _, c := range text.Data {
			if c == '\n' {
				cursorX = hStep
				cursorY += vStep
				continue
			}
			items = append(items, displayItem{cursorX, cursorY, c})
			cursorX += hStep
			if cursorX >= width-hStep {
				cursorY += vStep
				cursorX = hStep
			}
		}
	}
	return items
}

func (b *Browser) draw(items []displayItem) *container.Scroll {
	objects := make([]fyne.CanvasObject, len(items))
	for i, item := range items {
		t := canvas.NewText(string(item.Char), color.Black)
		t.Move(fyne.NewPos(item.X, item.Y))
		objects[i] = t
	}
	contentWidth := float32(width - hStep)
	contentHeight := items[len(items)-1].Y + vStep
	content := container.New(
		&sizedLayout{fyne.NewSize(contentWidth, contentHeight)},
		objects...,
	)
	scroll := container.NewScroll(content)
	b.window.SetContent(scroll)
	return scroll
}
