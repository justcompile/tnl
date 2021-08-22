package ui

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gizak/termui/v3"
)

type Window struct {
	components []Component
	onStart    func()
}

func (w *Window) GetComponent(id string) Component {
	for _, c := range w.components {
		if c.GetID() == id {
			return c
		}
	}
	return nil
}

func (w *Window) close() {
	for _, c := range w.components {
		c.Close()
	}

	termui.Close()
}

func (w *Window) resize() {
	width, _ := termui.TerminalDimensions()
	for _, c := range w.components {
		c.OnResize(width)
	}
}

func (w *Window) Run() error {
	defer w.close()

	for _, c := range w.components {
		go c.run()
	}

	if w.onStart != nil {
		w.onStart()
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	uiEvents := termui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.Type {
			case termui.ResizeEvent:
				w.resize()
			case termui.KeyboardEvent:
				if e.ID == "q" || e.ID == "<C-c>" {
					return nil
				}
			}
		case <-interrupt:
			return nil
		}
	}
}
