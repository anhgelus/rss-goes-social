package main

import (
	"github.com/anhgelus/rss-goes-social/api"
	"github.com/anhgelus/rss-goes-social/cli"
	"github.com/anhgelus/rss-goes-social/config"
	"github.com/anhgelus/rss-goes-social/feed"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	c := cli.CLI{
		Help: "RSS Goes Social is an application that adds RSS feeds to the Fediverse through a Mastodon-compatible API.",
		Commands: []*cli.Command{
			{
				Name:    "run",
				Help:    "Start the application",
				Flags:   nil,
				Handler: run,
			},
			{
				Name: "setup",
				Help: "Setup the Mastodon application and get the token. You have to give the url to the instance " +
					"that you want to use",
				Flags: []*cli.Flag{
					{
						Name: "id",
						Type: cli.FlagString,
						Help: "client_id of the application (not required)",
					},
					{
						Name: "secret",
						Type: cli.FlagString,
						Help: "client_secret of the application (not required)",
					},
				},
				Handler: api.Setup,
			},
		},
	}
	c.Handle()
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
