package main

import (
	"go-to-social-rss/api"
	"go-to-social-rss/config"
)

func main() {
	cfg := config.Config{}
	cfg.Load()
	api.LoadFeed(&cfg)
}
