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
	info := NewInfoList(5)
	ui.Body.AddRows(info.row)

	return &Screen{
		Downloads: newDownloadsList(ui.Body),
		Info:      info,
	}
}

type InfoList struct {
	maxItems int
	row      *ui.Row
	list     *ui.List
}

func (il *InfoList) AddInfo(info string) {
	// TODO:  Cap this at maxItems!
	il.list.Items = append(il.list.Items, info)
	ui.Render(il.row)
}

func NewInfoList(maxItems int) *InfoList {
	info := ui.NewList()
	info.Height = 12
	info.PaddingBottom = 2
	info.BorderLabel = "Topcast"
	info.BorderLabelFg = ui.ColorYellow
	info.ItemFgColor = ui.ColorCyan
	info.Items = []string{"Welcome to Topcast!  Press 'q' at any time to quit."}

	return &InfoList{
		maxItems: maxItems,
		row:      ui.NewRow(ui.NewCol(12, 0, info)),
		list:     info,
	}
}

type DownloadsList struct {
	widget   *ui.Grid
	elements []*DownloadElement
}

func (dl *DownloadsList) AddDownload(download *podcasts.Download) {
	elem := newDownloadElement(download)
	dl.elements = append(dl.elements, elem)

	dl.widget.AddRows(elem.row)
}

func (dl *DownloadsList) Update() {
	for _, element := range dl.elements {
		element.Update()

		if element.download.IsComplete() {
			// TODO:  Also send a message to the info queue
			// This will effectively hide the element
			element.guage.Height = 0
		}
	}
}

func newDownloadsList(parentGrid *ui.Grid) *DownloadsList {
	return &DownloadsList{
		widget: parentGrid,
	}
}

type DownloadElement struct {
	download *podcasts.Download
	row      *ui.Row
	guage    *ui.Gauge
}

func (d *DownloadElement) Update() {
	d.guage.Percent = d.download.PercentComplete()
	d.guage.Label = fmt.Sprintf("%d%% (%s)", d.download.PercentComplete(), d.download.Rate())
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
