package config

import (
	"context"
	"errors"
	"fmt"
	"github.com/anhgelus/rss-goes-social/utils"
	"github.com/pelletier/go-toml/v2"
	"github.com/redis/go-redis/v9"
	"os"
)

type Config struct {
	Version            string `toml:"version"`
	FetchEveryXMinutes uint   `toml:"fetch_every_X_minutes"`
	Redis              Redis  `toml:"redis"`
	Feeds              []Feed `toml:"feed"`
}

type Redis struct {
	Host     string `toml:"host"`
	Port     uint   `toml:"port"`
	Password string `toml:"password"`
}

type Feed struct {
	RssFeedUrl string   `toml:"rss_feed_url"`
	ServerUrl  string   `toml:"server_url"`
	Token      string   `toml:"token"`
	Enabled    bool     `toml:"enabled"`
	Language   string   `toml:"language"`
	Tags       []string `toml:"tags"`
}

const (
	Location = "config/config.toml"
)

func (cfg *Config) Load() {
	err := os.Mkdir("config", 0666)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
	data, err := os.ReadFile(Location)
	if errors.Is(err, os.ErrNotExist) {
		data, err = toml.Marshal(Config{
			Version:            "2",
			FetchEveryXMinutes: 5,
			Feeds: []Feed{
				{
					RssFeedUrl: "https://blog.example.org/rss",
					ServerUrl:  "https://gts.example.org",
					Token:      "account_token",
					Enabled:    false,
					Language:   "language of the feed (e.g. en, fr, de)",
					Tags:       []string{"tag-one", "tag-two"},
				},
			},
			Redis: Redis{
				Host:     "localhost",
				Port:     6379,
				Password: "",
			},
		})
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(Location, data, 0666)
		if err != nil {
			panic(err)
		}
	} else {
		panic(err)
	}
	err = toml.Unmarshal(data, cfg)
	if err != nil {
		panic(err)
	}
}

func (cfg *Config) GetRedis() (*redis.Client, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       0,
	})
	return c, c.Ping(context.Background()).Err()
}

func (f *Feed) GetUrl(uri string) string {
	return utils.GetFullUrl(f.ServerUrl, uri)
}
