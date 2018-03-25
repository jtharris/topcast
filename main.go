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

func main() {
	defer ui.Close()
	config := readConfig()
	manager := podcasts.NewDownloadManager(config.Settings.DownloadsDir)
	screen := setupScreen()
	screen.Update()

	go func() {
		for _, podcastConfig := range config.Podcasts {
			podcast := podcasts.NewPodcast(podcastConfig)
			screen.Info.AddInfo("Fetching episodes for:  " + podcast.Slug)
			podcast.Update()
			screen.Info.AddInfo("Got episodes for:  " + podcast.Title)

			for _, episode := range podcast.GetLatestEpisodes(config.Settings.MaxEpisodes) {
				dl, err := manager.StartDownload(episode)

				if err == nil {
					screen.Downloads.AddDownload(dl)
				} else {
					// TODO:  Wayyyy too aggressive??
					panic(err)
				}
			}
		}

		for _ = range time.Tick(500 * time.Millisecond) {
			screen.Update()
		}
	}()

	ui.Loop()
}
