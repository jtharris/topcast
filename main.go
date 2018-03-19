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
	manager := podcasts.NewDownloadManager()
	screen := setupScreen()
	config := readConfig()

	// TODO:  Provide more visibility here?
	for _, podcastConfig := range config.Podcasts {
		podcast := podcasts.NewPodcast(podcastConfig)
		podcast.Update()

		for _, episode := range podcast.GetLatestEpisodes(2) {
			dl, err := manager.StartDownload(episode)

			if err == nil {
				screen.AddDownload(dl)
			} else {
				// TODO:  Wayyyy too aggressive??
				panic(err)
			}
		}
	}

	go func() {
		for _ = range time.Tick(1 * time.Second) {
			screen.Update()
		}
	}()

	ui.Loop()
}