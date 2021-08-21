package ui

const maxDataRows = 3

type Component interface {
	GetUpdateChannel() chan interface{}
	GetID() string
	Close()
	OnResize(int)
	Render()
	UpdateOn(chan interface{})
	run()
}
