package podcasts

import (
	"github.com/jtharris/topcast/config"
	"github.com/mmcdole/gofeed"
)

type Podcast struct {
	Slug string
	URL  string

	Title string
	feed  *gofeed.Feed
}

func (p *Podcast) Update() error {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(p.URL)

	p.Title = feed.Title
	p.feed = feed

	return err
}

func (p *Podcast) GetLatestEpisodes(max int) []*Episode {
	if p.feed == nil {
		return []*Episode{}
	}

	if len(p.feed.Items) < max {
		max = len(p.feed.Items)
	}

	episodes := make([]*Episode, max)
	episodeIndex := 0

	for _, item := range p.feed.Items {
		episode := newEpisodeFromFeedItem(item)

		if episode.IsValid() {
			episode.Podcast = p
			episodes[episodeIndex] = episode
			episodeIndex++

			if episodeIndex >= max {
				break
			}
		}
	}

	return episodes
}

func NewPodcast(config config.PodcastConfig) *Podcast {
	return &Podcast{
		Slug:  config.Name,
		URL:   config.URL,
		Title: "(Unknown)",
	}
}

type Episode struct {
	Podcast     *Podcast
	Title       string
	Description string
	URL         string
}

func (e *Episode) IsValid() bool {
	return len(e.URL) > 3
}

func newEpisodeFromFeedItem(item *gofeed.Item) *Episode {
	var episodeURL string

	if len(item.Enclosures) > 0 {
		episodeURL = item.Enclosures[0].URL
	}

	return &Episode{
		Title:       item.Title,
		Description: item.Description,
		URL:         episodeURL,
	}
}
