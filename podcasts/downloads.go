package podcasts

import (
	"fmt"

	"github.com/cavaliercoder/grab"
	humanize "github.com/dustin/go-humanize"
)

type DownloadManager struct {
	Destination string
	client      *grab.Client
}

func (dm *DownloadManager) StartDownload(episode *Episode) (*Download, error) {
	req, err := grab.NewRequest(dm.Destination, episode.URL)

	if err != nil {
		return nil, err
	}

	download := &Download{
		episode:  episode,
		response: dm.client.Do(req),
	}

	return download, err
}

func NewDownloadManager() *DownloadManager {
	return &DownloadManager{
		Destination: "./downloads/", // Temp!  Get this from config
		client:      grab.NewClient(),
	}
}

type Download struct {
	episode  *Episode
	response *grab.Response
}

func (d *Download) Title() string {
	return d.episode.Title
}

func (d *Download) PercentComplete() int {
	return int(d.response.Progress() * 100)
}

func (d *Download) Rate() string {
	return fmt.Sprintf("%s / s", humanize.Bytes(uint64(d.response.BytesPerSecond())))
}
