package main

import (
	"go-to-social-rss/config"
	"go-to-social-rss/feed"
)

func main() {
	cfg := config.Config{}
	cfg.Load()
	feed.LoadFeed(&cfg)
}
