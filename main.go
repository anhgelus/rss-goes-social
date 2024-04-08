package main

import (
	"go-to-social-rss/api"
	"go-to-social-rss/config"
	"go-to-social-rss/feed"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Config{}
	slog.Info("Loading config...")
	cfg.Load()
	slog.Info("Config loaded")
	slog.Info("Loading feed...")
	feed.LoadFeed(&cfg)
	done, up := feed.InitFeedsChecker(&cfg)
	slog.Info("Feed loaded and initialized")

	go func() {
		for u := range up {
			err := api.PostNewContent(u.Item, u.F, &cfg)
			if err != nil {
				slog.Error(err.Error())
			}
		}
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	slog.Info("Good bye!")
	done <- true
}
