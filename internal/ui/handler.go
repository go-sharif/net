package ui

import (
	"fmt"

	"github.com/gizak/termui/v3"
	ui "github.com/gizak/termui/v3"
)

type Handler interface {
	Init() error
	GridSetup()
	Close()
	PollEvents()
	Refresh()
}

type baseHandler struct {
	grid *termui.Grid
}

type refreshOpt struct {
	Width  int
	Height int
	Clear  bool
}

func (uih *baseHandler) Refresh(opts ...refreshOpt) {
	if len(opts) > 0 {
		opt := opts[0]
		if opt.Height > 0 || opt.Width > 0 {
			uih.grid.SetRect(0, 0, opt.Width, opt.Height)
		}
		if opt.Clear {
			ui.Clear()
		}
	}
	ui.Render(uih.grid)
}

func (uih *baseHandler) Init() error {
	if err := ui.Init(); err != nil {
		return fmt.Errorf("failed to initialize termui: %v", err)
	}

	uih.grid = ui.NewGrid()
	return nil
}

func (uih *baseHandler) Close() {
	defer ui.Close()
}

func (uih *baseHandler) PollEvents(uiQuitChan chan<- struct{}) {
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.Type {
		case ui.KeyboardEvent:
			if e.ID == "q" || e.ID == "<C-c>" {
				close(uiQuitChan)
				return
			}
		case ui.ResizeEvent:
			payload := e.Payload.(ui.Resize)
			uih.Refresh(refreshOpt{Width: payload.Width, Height: payload.Height, Clear: true})
		}
	}
}
