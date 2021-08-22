package ui

import (
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type Table struct {
	done chan struct{}
	c    chan interface{}
	*widgets.Table
}

func (t *Table) Close() {
	t.done <- struct{}{}
}

func (t *Table) GetUpdateChannel() chan interface{} {
	return t.c
}

func (t *Table) GetID() string {
	return ComponentIDRequests
}

func (t *Table) OnResize(width int) {
	rect := t.GetRect()
	t.SetRect(rect.Min.X, rect.Min.Y, width, rect.Max.Y)

	t.Render()
}

func (t *Table) Render() {
	// ensure only previous N rows are shown
	rows := len(t.Rows)
	if rows > maxDataRows+1 {
		t.Rows = append(
			t.Rows[:1], // headers
			t.Rows[rows-maxDataRows:]...,
		)
	}
	termui.Render(t)
}

func (t *Table) UpdateOn(updates chan interface{}) {
	t.c = updates
}

func (t *Table) run() {
	defer close(t.c)

	t.Render()

	for {
		select {
		case <-t.done:
			close(t.done)
			return
		case v, ok := <-t.c:
			if !ok {
				return
			}

			t.Rows = append(t.Rows, v.([]string))
			t.Render()
		}
	}
}

func NewTable(t *widgets.Table) *Table {
	return &Table{
		done:  make(chan struct{}, 1),
		Table: t,
	}
}
