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

	for i := 0; i < max; i++ {
		episodes[i] = newEpisodeFromFeedItem(p.feed.Items[i])
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
	Title       string
	Description string
	URL         string
}

func newEpisodeFromFeedItem(item *gofeed.Item) *Episode {
	return &Episode{
		Title:       item.Title,
		Description: item.Description,
		URL:         item.Link,
	}
}
