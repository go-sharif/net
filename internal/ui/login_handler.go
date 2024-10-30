package ui

import (
	"fmt"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"

	"github.com/go-sharif/net/internal/model"
)

type LoginHandler struct {
	baseHandler
	sessionStatusTable *widgets.Table
	logList            *widgets.List
	pingList           *widgets.List
	byteDownSparkLine  *widgets.SparklineGroup
	byteUpSparkLine    *widgets.SparklineGroup
}

func (lh *LoginHandler) Init() error {
	if err := ui.Init(); err != nil {
		return fmt.Errorf("failed to initialize termui: %v", err)
	}

	lh.sessionStatusTable = initStatusTable()
	lh.logList = initLogList()
	lh.pingList = pingList()
	lh.byteDownSparkLine = initByteDownSparkLine()
	lh.byteUpSparkLine = initByteUpSparkLine()
	lh.gridSetup()

	return nil
}

func (lh *LoginHandler) gridSetup() {
	lh.grid = ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	lh.grid.SetRect(0, 0, termWidth, termHeight)

	lh.grid.Set(
		ui.NewRow(1.0/4,
			ui.NewCol(1.0, lh.sessionStatusTable),
		),
		ui.NewRow(1.0/4,
			ui.NewCol(1.0/2, lh.byteDownSparkLine),
			ui.NewCol(1.0/2, lh.byteUpSparkLine),
		),
		ui.NewRow(2.0/4,
			ui.NewCol(1.0/2, lh.logList),
			ui.NewCol(1.0/2, lh.pingList),
		),
	)
	ui.Render(lh.grid)
}

func (lh *LoginHandler) UpdateStatusTable(status *model.SessionStatus) {
	// change the second row
	lh.sessionStatusTable.Rows[1] = []string{
		status.Username,
		status.IPAddress,
		status.SessionTime,
		status.TimeLeft,
		status.BytesUp.ToString(),
		status.BytesDown.ToString(),
	}
}

var (
	LogInfo = ui.Theme.List.Text
	LogErr  = ui.NewStyle(ui.ColorRed)
	LogSucc = ui.NewStyle(ui.ColorGreen)
)

func (lh *LoginHandler) AddLog(l string, style ui.Style) {
	lh.logList.Rows = append(lh.logList.Rows, l)

	lh.logList.TextStyle = style

	lh.logList.ScrollDown()

	lh.Refresh()
}

func (lh *LoginHandler) AddPing(p string, style ui.Style) {
	lh.pingList.Rows = append(lh.pingList.Rows, p)

	lh.pingList.TextStyle = ui.Theme.List.Text

	lh.pingList.ScrollDown()

	lh.Refresh()
}

// Since both BytesDown and BytesUp data updates in sync, no need
// for separate updator. This way, will double ui refresh.
func (lh *LoginHandler) AddBytesData(u, d float64) {

	lh.byteUpSparkLine.Sparklines[0].Data = append(lh.byteUpSparkLine.Sparklines[0].Data, u)
	lh.byteDownSparkLine.Sparklines[0].Data = append(lh.byteDownSparkLine.Sparklines[0].Data, d)

	// check if it's full, pop the first entries each time after adding the new entry
	dataLen := len(lh.byteDownSparkLine.Sparklines[0].Data)
	if lh.byteDownSparkLine.Inner.Dx() <= dataLen {
		d := dataLen - lh.byteDownSparkLine.Inner.Dx() + 1
		lh.byteDownSparkLine.Sparklines[0].Data = lh.byteDownSparkLine.Sparklines[0].Data[d:]
		lh.byteUpSparkLine.Sparklines[0].Data = lh.byteUpSparkLine.Sparklines[0].Data[d:]
	}

}

func initStatusTable() *widgets.Table {

	t := widgets.NewTable()
	t.Title = "Session Status"
	// set the header
	t.Rows = [][]string{
		{"Username", "IP Address", "Session Time", "Time Left", "Bytes Up", "Bytes Down"},
		{"", "", "", "", "", ""},
	}

	t.TextStyle = ui.NewStyle(ui.ColorWhite)
	t.TextAlignment = ui.AlignCenter
	t.RowSeparator = false
	t.FillRow = true

	return t
}

func initLogList() *widgets.List {
	l := widgets.NewList()
	l.Title = "Logs"
	l.Rows = make([]string, 0)
	l.Border = true

	return l
}

func pingList() *widgets.List {
	l := widgets.NewList()
	l.Title = "Internet Status"
	l.Rows = make([]string, 0)
	l.Border = true

	return l
}

func initByteDownSparkLine() *widgets.SparklineGroup {
	sl := widgets.NewSparkline()
	sl.Data = make([]float64, 0)
	sl.LineColor = ui.ColorBlue

	slg := widgets.NewSparklineGroup(sl)
	slg.Title = "Bytes Down"
	return slg
}

func initByteUpSparkLine() *widgets.SparklineGroup {
	sl := widgets.NewSparkline()
	sl.Data = make([]float64, 0)
	sl.LineColor = ui.ColorMagenta

	slg := widgets.NewSparklineGroup(sl)

	slg.Title = "Bytes Up"

	return slg
}
