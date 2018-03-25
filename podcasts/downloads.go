package podcasts

import (
	"fmt"
	"path/filepath"

	"github.com/cavaliercoder/grab"
	humanize "github.com/dustin/go-humanize"
)

type DownloadManager struct {
	downloadsDir string
	client       *grab.Client
}

func (dm *DownloadManager) StartDownload(episode *Episode) (*Download, error) {
	podcastDir := filepath.Join(dm.downloadsDir, episode.Podcast.Slug) + "/"
	req, err := grab.NewRequest(podcastDir, episode.URL)

	if err != nil {
		return nil, err
	}

	download := &Download{
		episode:  episode,
		response: dm.client.Do(req),
	}

	return download, err
}

func NewDownloadManager(downloadsDir string) *DownloadManager {
	return &DownloadManager{
		downloadsDir: downloadsDir,
		client:       grab.NewClient(),
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
	return fmt.Sprintf("%s/s", humanize.Bytes(uint64(d.response.BytesPerSecond())))
}

func (d *Download) TimeLeft() string {
	return humanize.Time(d.response.ETA())
}

func (d *Download) IsComplete() bool {
	return d.response.IsComplete()
}
