package feed

import (
	"context"
	"errors"
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/redis/go-redis/v9"
	"go-to-social-rss/api"
	"go-to-social-rss/config"
	"log/slog"
	"sort"
)

const (
	KeyLastFeedUrl   = "last_feed:url"
	KeyLastFeedTitle = "last_feed:title"
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
		err := api.VerifyToken(f.GetUrl("/api/v1/accounts/verify_credentials"), f.Token)
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

func UpdateLast(item *gofeed.Item, f *config.Feed, cfg *config.Config) error {
	r, err := cfg.GetRedis()
	if err != nil {
		return err
	}
	ctx := context.Background()

	if r.Set(ctx, f.RssFeedUrl+":"+KeyLastFeedUrl, item.Link, 0).Err() != nil {
		return err
	}
	if r.Set(ctx, f.RssFeedUrl+":"+KeyLastFeedTitle, item.Title, 0).Err() != nil {
		return err
	}
	return nil
}

func checkFeed(f *config.Feed, cfg *config.Config) (*gofeed.Item, error) {
	fp := gofeed.NewParser()
	parsed, err := fp.ParseURL(f.RssFeedUrl)
	if err != nil {
		return nil, err
	}
	if len(parsed.Items) == 0 {
		return nil, nil
	}
	sort.Sort(parsed)
	last := parsed.Items[len(parsed.Items)-1]
	r, err := cfg.GetRedis()
	if err != nil {
		return nil, err
	}

	v, err := checkNewValue(r, last.Link, f.RssFeedUrl+":"+KeyLastFeedUrl)
	if err != nil {
		return nil, err
	}
	err = r.Close()
	if err != nil {
		slog.Error(err.Error())
	}
	if v {
		return last, nil
	}
	v, err = checkNewValue(r, last.Title, f.RssFeedUrl+":"+KeyLastFeedTitle)
	if err != nil {
		return nil, err
	}
	if v {
		return last, nil
	}
	return nil, nil
}

func checkNewValue(r *redis.Client, last string, key string) (bool, error) {
	ctx := context.Background()

	res := r.Get(ctx, key)
	not := false
	err := res.Err()
	if err != nil {
		not = errors.Is(err, redis.Nil)
		if !not {
			return false, err
		} else {
			return true, nil
		}
	}
	return res.String() != last, nil
}
