package ui

import (
	"fmt"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type BannerText struct {
	Endpoint string
	Port     string
}

type Banner struct {
	text *BannerText
	done chan struct{}
	c    chan interface{}
	*widgets.Paragraph
}

func (b *Banner) GetUpdateChannel() chan interface{} {
	return b.c
}

func (b *Banner) GetID() string {
	return ComponentIDInfo
}

func (b *Banner) Close() {
	b.done <- struct{}{}
}

func (b *Banner) OnResize(width int) {
	rect := b.GetRect()
	b.SetRect(rect.Min.X, rect.Min.Y, width, rect.Max.Y)

	b.Render()
}

func (b *Banner) Render() {
	b.Text = fmt.Sprintf("[Endpoint](fg:cyan):      %s\n[Local Address](fg:cyan): %s", b.text.Endpoint, b.text.Port)
	termui.Render(b)
}

func (b *Banner) UpdateOn(c chan interface{}) {
	b.c = c
}

func (b *Banner) run() {
	defer close(b.c)
	b.Render()

	for {
		select {
		case <-b.done:
			close(b.done)
			return
		case u, ok := <-b.c:
			if !ok {
				return
			}

			if text, isType := u.(*BannerText); isType {
				b.text = text
				b.Render()
			}
		}
	}
}

func NewBanner(listenAddress string) *Banner {
	width, _ := termui.TerminalDimensions()

	p := widgets.NewParagraph()
	p.Title = "tnl"
	p.TextStyle = termui.NewStyle(termui.ColorWhite)
	p.SetRect(0, 0, width, 5)

	return &Banner{
		text:      &BannerText{Port: listenAddress},
		done:      make(chan struct{}, 1),
		Paragraph: p,
	}
}
