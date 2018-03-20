package screen

import (
	"fmt"

	ui "github.com/gizak/termui"
	"github.com/jtharris/topcast/podcasts"
)

// TODO:  It is weird to use an object that manipulates the global ui.Body - rethink this!

type Screen struct {
	downloads []*DownloadElement
}

func (s *Screen) AddDownload(download *podcasts.Download) {
	elem := newDownloadElement(download)
	s.downloads = append(s.downloads, elem)

	ui.Body.AddRows(elem.row)
}

func (s *Screen) Update() {
	for _, download := range s.downloads {
		download.Update()
	}

	ui.Clear()
	ui.Body.Align()
	ui.Render(ui.Body)
}

func NewScreen() *Screen {
	headerCol := ui.NewCol(8, 2, newHeader())
	ui.Body.AddRows(ui.NewRow(headerCol))

	return &Screen{}
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
	guage.Height = 5

	return &DownloadElement{
		download: d,
		row:      ui.NewRow(ui.NewCol(12, 0, guage)),
		guage:    guage,
	}
}

func newHeader() *ui.Par {
	p := ui.NewPar("PRESS q TO QUIT")
	p.Height = 3
	p.BorderLabel = "TopCast"
	p.BorderFg = ui.ColorCyan

	return p
}
