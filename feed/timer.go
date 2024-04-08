package feed

import (
	"github.com/mmcdole/gofeed"
	"go-to-social-rss/config"
	"log/slog"
	"time"
)

func InitFeedsChecker(cfg *config.Config) (chan<- bool, <-chan *gofeed.Item) {
	t := time.NewTicker(5 * time.Minute)
	done := make(chan bool)
	up := make(chan *gofeed.Item)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-t.C:
				checkFeeds(cfg, up)
			}
		}
	}()
	return done, up
}

func checkFeeds(cfg *config.Config, up chan<- *gofeed.Item) {
	for _, f := range validFeed {
		f, err := checkFeed(f, cfg)
		if err != nil {
			slog.Error(err.Error())
		}
		up <- f
	}
}
