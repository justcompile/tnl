package ui

import (
	"context"
	"fmt"
	"time"

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
	// go pollRequests(requestUpdates)

	infoUpdates := make(chan interface{})
	banner.UpdateOn(infoUpdates)
	// go pollInfo(infoUpdates)

	return window, nil
}

func pollInfo(updateChan chan interface{}) {
	time.AfterFunc(time.Second, func() {
		updateChan <- &BannerText{
			Endpoint: "https://meh.tnl.justcompile.io",
			Port:     ":3000",
		}
	})
}

func pollRequests(updateChan chan interface{}) {
	t := time.NewTicker(time.Second * 2)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*10))
	defer cancel()
	i := 3
	for {
		select {
		case v := <-t.C:
			updateChan <- []string{fmt.Sprintf("%d) col", i), "col2", v.String()}
			i++
		case <-ctx.Done():
			t.Stop()
			return
		}
	}
}
