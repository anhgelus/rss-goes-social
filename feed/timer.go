package feed

import (
	"github.com/mmcdole/gofeed"
	"go-to-social-rss/config"
	"log/slog"
	"time"
)

type Up struct {
	Item *gofeed.Item
	F    *config.Feed
}

func InitFeedsChecker(cfg *config.Config) (chan<- bool, <-chan *Up) {
	t := time.NewTicker(5 * time.Minute)
	done := make(chan bool)
	up := make(chan *Up)

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

func checkFeeds(cfg *config.Config, up chan<- *Up) {
	for _, f := range validFeed {
		go func() {
			item, err := checkFeed(f, cfg)
			if err != nil {
				slog.Error(err.Error())
			}
			up <- &Up{Item: item, F: f}
		}()
	}
}
