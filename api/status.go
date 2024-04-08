package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mmcdole/gofeed"
	"go-to-social-rss/config"
	"net/http"
	"strings"
)

var (
	ErrRequestFailed = errors.New("request failed")
)

const (
	lengthMax = 500
)

type postStatus struct {
	Status      string `json:"status"`
	Visibility  string `json:"visibility"`
	Language    string `json:"language"`
	ContentType string `json:"content_type"`
	Sensitive   bool   `json:"sensitive"`
	SpoilerText string `json:"spoiler_text"`
	Federated   bool   `json:"federated"`
	Boostable   bool   `json:"boostable"`
	Replyable   bool   `json:"replyable"`
	Likeable    bool   `json:"likeable"`
}

func PostNewContent(item *gofeed.Item, f *config.Feed, cfg *config.Config) error {
	b, err := json.Marshal(genStatus(item, f))
	if err != nil {
		return err
	}
	req, err := newRequest(http.MethodPost, f.GetUrl("/api/v1/statuses"), f.Token, strings.NewReader(string(b)))
	if err != nil {
		return err
	}
	resp, err := doRequest(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Join(ErrRequestFailed, errors.New(fmt.Sprintf("with the status code %d", resp.StatusCode)))
	}
	return nil
}

func genStatus(item *gofeed.Item, f *config.Feed) *postStatus {
	// length max - title - link - "..." - "\n\n" - "\n\n"
	l := lengthMax - len(item.Title) - len(item.Link) - 3 - 2 - 2
	content := ""
	split := strings.Split(item.Content, " ")
	i := 0
	for len(content+" "+split[i]) < l && i < len(split) {
		content += " " + split[i]
		i++
	}
	return &postStatus{
		Status:      fmt.Sprintf("%s\n\n%s...\n%s", item.Title, content, item.Link),
		Visibility:  "public",
		Language:    f.Language,
		ContentType: "text/plain",
		Sensitive:   false,
		SpoilerText: "",
		Federated:   true,
		Boostable:   true,
		Replyable:   true,
		Likeable:    true,
	}
}
