package config

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

type TopCastConfig struct {
	Podcasts []PodcastConfig `yaml:"podcasts"`
}

type PodcastConfig struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

func NewTopCastConfig(file string) (TopCastConfig, error) {
	config := TopCastConfig{}

	configData, err := ioutil.ReadFile(file)

	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(configData, &config)

	if err != nil {
		return config, err
	}

	return config, err
}
