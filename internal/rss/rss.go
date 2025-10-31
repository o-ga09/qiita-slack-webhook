package rss

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/o-ga09/qiita-slack-webhook/internal/config"
	"github.com/o-ga09/qiita-slack-webhook/internal/notifier"
)

type RSSFeed struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Items       []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Author      string `xml:"author"`
}

// fetchRSSFeed はRSSフィードから最新記事を取得
func FetchRSSFeed(feedURL string) (*notifier.SlackMessage, error) {
	resp, err := http.Get(feedURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch RSS feed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("RSS feed request failed with status: %d", resp.StatusCode)
	}

	var feed RSSFeed
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, fmt.Errorf("failed to parse RSS feed: %w", err)
	}

	return &notifier.SlackMessage{}, nil
}

// getLatestRSSArticles は指定した件数の最新記事を取得
func GetLatestRSSArticles(cfg config.Config) (*notifier.SlackMessage, error) {
	items, err := FetchRSSFeed(cfg.RSSFeedURL)
	if err != nil {
		return nil, err
	}

	// 日付でソート（新しい順）
	sort.Slice(items, func(i, j int) bool {
		ti, _ := time.Parse(time.RFC1123Z, items[i].PubDate)
		tj, _ := time.Parse(time.RFC1123Z, items[j].PubDate)
		return ti.After(tj)
	})

	// 指定件数まで取得
	count := min(limit, len(items))
	return &notifier.SlackMessage{}, nil
}
