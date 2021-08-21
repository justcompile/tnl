package ui

import (
	"fmt"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func ConstructUI(listenAddress string) (*Window, error) {
	if err := termui.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize termui: %v", err)
	}

	width, height := termui.TerminalDimensions()

	exampleTable := widgets.NewTable()
	exampleTable.Rows = [][]string{
		{"Path", "Status", "Response Time"},
	}
	exampleTable.TextStyle = termui.NewStyle(termui.ColorWhite)
	exampleTable.Border = false

	banner := NewBanner(listenAddress)
	exampleTable.SetRect(0, banner.GetRect().Max.Y, width, height-banner.GetRect().Max.Y)

	table1 := NewTable(exampleTable)

	window := &Window{
		components: []Component{
			banner,
			table1,
		},
	}

	requestUpdates := make(chan interface{})
	table1.UpdateOn(requestUpdates)

	infoUpdates := make(chan interface{})
	banner.UpdateOn(infoUpdates)

	return window, nil
}
