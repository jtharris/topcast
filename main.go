package main

import (
	"time"

	ui "github.com/gizak/termui"
	"github.com/jtharris/topcast/config"
	"github.com/jtharris/topcast/podcasts"
	"github.com/jtharris/topcast/screen"
)

func setupScreen() *screen.Screen {
	err := ui.Init()
	if err != nil {
		panic(err)
	}

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	return screen.NewScreen()
}

func readConfig() config.TopCastConfig {
	config, err := config.NewTopCastConfig("./topcast.yaml")

	if err != nil {
		panic(err)
	}

	return config
}

func loadEpisodes(podcast *podcasts.Podcast, maxEpisodes int, episodes chan *podcasts.Episode, info chan string) {
	info <- "Fetching episodes for:  " + podcast.Slug
	podcast.Update()
	info <- "Got episodes for:  " + podcast.Title

	for _, episode := range podcast.GetLatestEpisodes(maxEpisodes) {
		episodes <- episode
	}
}

func main() {
	defer ui.Close()
	config := readConfig()
	manager := podcasts.NewDownloadManager(config.Settings.DownloadsDir)
	screen := setupScreen()

	go screen.Info.Start()

	episodeChan := make(chan *podcasts.Episode, 100)

	go func() {
		for _ = range time.Tick(500 * time.Millisecond) {
			screen.Update()
		}
	}()

	go func() {
		for episode := range episodeChan {
			dl, err := manager.StartDownload(episode)

			if err == nil {
				screen.Downloads.AddDownload(dl)
			} else {
				// TODO:  Wayyyy too aggressive??
				panic(err)
			}
		}
	}()

	for _, podcastConfig := range config.Podcasts {
		podcast := podcasts.NewPodcast(podcastConfig)
		go loadEpisodes(podcast, config.Settings.MaxEpisodes, episodeChan, screen.Info.InfoQueue)
	}

	ui.Loop()
}
