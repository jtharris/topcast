package screen

import (
	"fmt"

	ui "github.com/gizak/termui"
	"github.com/jtharris/topcast/podcasts"
)

// TODO:  It is weird to use an object that manipulates the global ui.Body - rethink this!
type Screen struct {
	Header *ui.Par

	Downloads *DownloadsList
	Info      *InfoList
}

func (s *Screen) Update() {
	s.Downloads.Update()

	ui.Clear()
	ui.Body.Align()
	ui.Render(ui.Body)
}

func NewScreen() *Screen {
	// TODO:  Hard coding!
	infoQueue := make(chan string, 25)
	info := NewInfoList(infoQueue, 11)
	ui.Body.AddRows(info.row)

	return &Screen{
		Downloads: newDownloadsList(infoQueue, ui.Body),
		Info:      info,
	}
}

type InfoList struct {
	InfoQueue chan string
	maxItems  int
	row       *ui.Row
	list      *ui.List
}

func (il *InfoList) Start() {
	for m := range il.InfoQueue {
		il.addInfo(m)
	}
}

func (il *InfoList) addInfo(info string) {
	il.list.Items = append(il.list.Items, info)

	itemLength := len(il.list.Items)
	if itemLength > il.maxItems {
		il.list.Items = il.list.Items[itemLength-il.maxItems : itemLength]
	}

	ui.Render(il.row)
}

func NewInfoList(infoQueue chan string, maxItems int) *InfoList {
	info := ui.NewList()
	info.Height = 12
	info.PaddingBottom = 2
	info.BorderLabel = "Topcast"
	info.BorderLabelFg = ui.ColorYellow
	info.ItemFgColor = ui.ColorCyan
	info.Items = []string{"Welcome to Topcast!  Press 'q' at any time to quit."}

	return &InfoList{
		InfoQueue: infoQueue,
		maxItems:  maxItems,
		row:       ui.NewRow(ui.NewCol(12, 0, info)),
		list:      info,
	}
}

type DownloadsList struct {
	infoQueue chan string
	widget    *ui.Grid
	elements  []*DownloadElement
}

func (dl *DownloadsList) AddDownload(download *podcasts.Download) {
	elem := newDownloadElement(download)
	dl.elements = append(dl.elements, elem)

	dl.widget.AddRows(elem.row)
}

func (dl *DownloadsList) Update() {
	for _, element := range dl.elements {
		if element.download != nil {
			element.Update()

			if element.download.IsComplete() {
				// This will effectively hide the element
				element.guage.Height = 0
				dl.infoQueue <- "Download Complete:  " + element.download.Title()
				element.download = nil
			}
		}
	}
}

func newDownloadsList(infoQueue chan string, parentGrid *ui.Grid) *DownloadsList {
	return &DownloadsList{
		infoQueue: infoQueue,
		widget:    parentGrid,
	}
}

type DownloadElement struct {
	download *podcasts.Download
	row      *ui.Row
	guage    *ui.Gauge
}

func (d *DownloadElement) Update() {
	d.guage.Percent = d.download.PercentComplete()
	d.guage.Label = fmt.Sprintf("%d%% (%s)  ETA: %s", d.download.PercentComplete(), d.download.Rate(), d.download.TimeLeft())
}

func newDownloadElement(d *podcasts.Download) *DownloadElement {
	guage := ui.NewGauge()
	guage.BorderLabel = d.Title()
	guage.Height = 3

	return &DownloadElement{
		download: d,
		row:      ui.NewRow(ui.NewCol(12, 0, guage)),
		guage:    guage,
	}
}
