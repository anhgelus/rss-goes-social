package config

import (
	"errors"
	"github.com/pelletier/go-toml/v2"
	"os"
)

type Config struct {
	Version string `toml:"version"`
	Feeds   []Feed `toml:"feed"`
}

type Feed struct {
	Url     string `toml:"url"`
	Token   string `toml:"token"`
	Enabled bool   `toml:"enabled"`
}

const (
	Location = "config.toml"
)

func (cfg *Config) Load() {
	data, err := os.ReadFile(Location)
	if errors.Is(err, os.ErrNotExist) {
		data, err = toml.Marshal(Config{
			Version: "1",
			Feeds: []Feed{
				{
					Url:     "rss_feed_url",
					Token:   "account_token",
					Enabled: false,
				},
			},
		})
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(Location, data, 0666)
		if err != nil {
			panic(err)
		}
	}
	err = toml.Unmarshal(data, cfg)
	if err != nil {
		panic(err)
	}
}
