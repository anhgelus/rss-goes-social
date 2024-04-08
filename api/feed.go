package api

import (
	"fmt"
	"go-to-social-rss/config"
	"log/slog"
)

var (
	validFeed []*config.Feed
)

func LoadFeed(cfg *config.Config) []*config.Feed {
	enabled := 0
	loaded := 0
	for _, f := range cfg.Feeds {
		if !f.Enabled {
			continue
		}
		enabled++
		err := VerifyToken(f.Url, f.Token)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		validFeed = append(validFeed, &f)
		loaded++
	}
	slog.Info(fmt.Sprintf("Loaded %d feed(s)/%d", loaded, enabled))
	return validFeed
}
