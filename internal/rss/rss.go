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
	XMLName     xml.Name  `xml:"feed"`
	Title       string    `xml:"title"`
	Description string    `xml:"description"`
	Updated     string    `xml:"updated"`
	Entries     []RSSItem `xml:"entry"`
}

type RSSItem struct {
	Title     string     `xml:"title"`
	Link      AtomLink   `xml:"link"`
	URL       string     `xml:"url"`
	Content   string     `xml:"content"`
	Published string     `xml:"published"`
	Updated   string     `xml:"updated"`
	Author    AtomAuthor `xml:"author"`
}

type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type AtomAuthor struct {
	Name string `xml:"name"`
}

// fetchRSSFeed はRSSフィードから最新記事を取得
func FetchRSSFeed(feedURL string) (*RSSFeed, error) {
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

	return &feed, nil
}

// getLatestRSSArticles は指定した件数の最新記事を取得
func GetLatestRSSArticles(cfg config.Config) (*notifier.SlackMessage, error) {
	rss, err := FetchRSSFeed(cfg.RSSFeedURL)
	if err != nil {
		return nil, err
	}

	// 日付でソート（新しい順）
	sort.Slice(rss.Entries, func(i, j int) bool {
		ti, _ := time.Parse(time.RFC3339, rss.Entries[i].Published)
		tj, _ := time.Parse(time.RFC3339, rss.Entries[j].Published)
		return ti.After(tj)
	})

	return toSlackMessage(rss), nil
}

func toSlackMessage(rss *RSSFeed) *notifier.SlackMessage {
	latestArticle := rss.Entries[0]
	// URLはlink要素のhref属性またはurl要素から取得
	articleURL := latestArticle.Link.Href
	if articleURL == "" {
		articleURL = latestArticle.URL
	}

	message := fmt.Sprintf(`
## LTS グループ Qiita アドベントカレンダー 2025

%s
%s

%s
`, latestArticle.Title,
		latestArticle.Content,
		articleURL,
	)
	return &notifier.SlackMessage{
		Text: message,
	}
}
