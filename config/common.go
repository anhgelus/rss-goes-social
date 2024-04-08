package config

import (
	"context"
	"errors"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"github.com/redis/go-redis/v9"
	"os"
)

type Config struct {
	Version string `toml:"version"`
	Redis   Redis  `toml:"redis"`
	Feeds   []Feed `toml:"feed"`
}

type Redis struct {
	Host     string `toml:"host"`
	Port     uint   `toml:"port"`
	Password string `toml:"password"`
}

type Feed struct {
	RssFeedUrl string `toml:"rss_feed_url"`
	ServerUrl  string `toml:"server_url"`
	Token      string `toml:"token"`
	Enabled    bool   `toml:"enabled"`
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
					RssFeedUrl: "https://blog.example.org/rss",
					ServerUrl:  "https://gts.example.org",
					Token:      "account_token",
					Enabled:    false,
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
