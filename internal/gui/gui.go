package gui

import (
	"image/color"
	"strings"

	"gobrowser/internal/page"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

const (
	width, height = 800, 600
	padding       = 13
)

type displayItem struct {
	X, Y float32
	Char string
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
			scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(0, 20)})
		case fyne.KeyDown:
			scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(0, -20)})
		}
	})
	b.window.ShowAndRun()
}

func (b *Browser) layout(tokens []page.Token) []displayItem {
	var items []displayItem
	cursorX, cursorY := float32(padding), float32(padding)
	spaceWidth := fyne.MeasureText(" ", 16, fyne.TextStyle{}).Width
	for _, token := range tokens {
		text, ok := token.(page.Text)
		if !ok {
			continue
		}
		for _, word := range strings.Fields(text.Data) {
			size := fyne.MeasureText(word, 16, fyne.TextStyle{})
			if cursorX+size.Width > width-padding {
				cursorX = padding
				cursorY += size.Height * 1.25
			}
			items = append(items, displayItem{cursorX, cursorY, word})
			cursorX += size.Width + spaceWidth
		}
	}
	return items
}

func (b *Browser) draw(items []displayItem) *container.Scroll {
	objects := make([]fyne.CanvasObject, len(items))
	for i, item := range items {
		t := canvas.NewText(item.Char, color.Black)
		t.Move(fyne.NewPos(item.X, item.Y))
		objects[i] = t
	}
	contentWidth := float32(width - padding)
	contentHeight := items[len(items)-1].Y + padding
	content := container.New(
		&sizedLayout{fyne.NewSize(contentWidth, contentHeight)},
		objects...,
	)
	scroll := container.NewScroll(content)
	b.window.SetContent(scroll)
	return scroll
}
