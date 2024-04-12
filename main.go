package main

import (
	"github.com/anhgelus/rss-goes-social/api"
	"github.com/anhgelus/rss-goes-social/config"
	"github.com/anhgelus/rss-goes-social/feed"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {

}

func run() {
	cfg := config.Config{}
	slog.Info("Loading config...")
	cfg.Load()
	slog.Info("Config loaded")
	slog.Info("Loading feed(s)...")
	feeds := feed.LoadFeed(&cfg)
	if len(feeds) == 0 {
		slog.Warn("No feed loaded. Exiting.")
		os.Exit(0)
	}
	done, up := feed.InitFeedsChecker(&cfg)
	slog.Info("Feed(s) loaded and initialized")

	go func() {
		for u := range up {
			slog.Info("New post", "feed", u.F.RssFeedUrl)
			err := api.PostNewContent(u.Item, u.F, &cfg)
			if err != nil {
				slog.Error(err.Error())
				continue
			}
			err = feed.UpdateLast(u.Item, u.F, &cfg)
			if err != nil {
				slog.Error(err.Error())
			}
		}
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	slog.Info("Exiting.")
	done <- true
}
