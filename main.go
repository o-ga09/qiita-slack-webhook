package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	aggregatetags "github.com/o-ga09/qiita-slack-webhook/internal/aggregateTags"
	"github.com/o-ga09/qiita-slack-webhook/internal/config"
	"github.com/o-ga09/qiita-slack-webhook/internal/notifier"
	"github.com/o-ga09/qiita-slack-webhook/internal/rss"
)

const (
	RETRY_COUNT = 3
)

var (
	mode       = flag.String("mode", "message", "動作モード: message, aggregate, rss")
	tag        = flag.String("tag", "Go", "Qiitaのタグ名（aggregateモード用）")
	maxPages   = flag.Int("pages", 5, "取得する最大ページ数（aggregateモード用）")
	rssFeedURL = flag.String("rss", "", "RSSフィードのURL（rssモード用）")
	rssLimit   = flag.Int("limit", 10, "取得する記事数（rssモード用）")
	message    = flag.String("message", "", "送信するメッセージ（messageモード用）")
	help       = flag.Bool("help", false, "ヘルプを表示")
)

func init() {
	flag.Parse()
}

func main() {
	if *help {
		printHelp()
		return
	}

	config := config.Config{
		Mode:       *mode,
		Tag:        *tag,
		MaxPages:   *maxPages,
		RSSFeedURL: *rssFeedURL,
		RSSLimit:   *rssLimit,
		Message:    *message,
	}

	ctx := context.Background()

	var slackMessage *notifier.SlackMessage
	var err error

	switch config.Mode {
	case "message":
		slackMessage = handleMessageMode(config)
	case "aggregate":
		slackMessage, err = aggregatetags.AggregateLikes(config)
		if err != nil {
			log.Fatalf("Failed to aggregate likes: %v", err)
		}
	case "rss":
		slackMessage, err = rss.GetLatestRSSArticles(config)
		if err != nil {
			log.Fatalf("Failed to aggregate likes: %v", err)
		}
	default:
		log.Fatalf("Unknown mode: %s. Use -help for usage information.", config.Mode)
	}

	// HTTPステータスが200以外の場合のみ3回リトライする
	for range RETRY_COUNT {
		err := notifier.SendSlackNotification(ctx, *slackMessage)
		if err == nil || !errors.Is(err, notifier.ErrHTTPStatusNotOK) {
			if err != nil {
				log.Fatalf("Failed to send notification: %v", err)
			}
			log.Println("Notification sent successfully!")
			break
		}
	}
}

func printHelp() {
	fmt.Println("Qiita Slack Webhook Tool")
	fmt.Println("\n使い方:")
	fmt.Println("  main [options]")
	fmt.Println("\nオプション:")
	flag.PrintDefaults()
	fmt.Println("\n例:")
	fmt.Println("  # シンプルなメッセージ送信")
	fmt.Println("  main -mode=message -message=\"Hello, Slack!\"")
	fmt.Println("\n  # Qiitaタグのいいね集計")
	fmt.Println("  main -mode=aggregate -tag=Go -pages=5")
	fmt.Println("\n  # RSSフィードから最新記事取得")
	fmt.Println("  main -mode=rss -rss=\"https://qiita.com/tags/Go/feed\" -limit=10")
}

func handleMessageMode(config config.Config) *notifier.SlackMessage {
	msg := config.Message
	if msg == "" {
		msg = os.Getenv("SLACK_MESSAGE")
		if msg == "" {
			msg = "Hello from Qiita Slack Bot!"
		}
	}
	return &notifier.SlackMessage{Text: msg}
}
