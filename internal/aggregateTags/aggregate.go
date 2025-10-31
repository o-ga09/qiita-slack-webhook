package aggregatetags

import (
	"sort"

	"github.com/o-ga09/qiita-slack-webhook/internal/config"
	"github.com/o-ga09/qiita-slack-webhook/internal/notifier"
)

type LikeSummary struct {
	TotalLikes  int
	TotalItems  int
	TopArticles []QiitaItem
}

// aggregateLikes はタグの記事のいいね数を集計
func AggregateLikes(cfg config.Config) (*notifier.SlackMessage, error) {
	var allItems []QiitaItem
	perPage := 100

	for page := 1; page <= cfg.MaxPages; page++ {
		items, err := fetchQiitaItemsByTag(cfg.Tag, perPage, page)
		if err != nil {
			return nil, err
		}

		if len(items) == 0 {
			break
		}

		allItems = append(allItems, items...)

		// 取得した記事数が perPage より少ない場合は最後のページ
		if len(items) < perPage {
			break
		}
	}

	// いいね数でソート
	sort.Slice(allItems, func(i, j int) bool {
		return allItems[i].LikesCount > allItems[j].LikesCount
	})

	totalLikes := 0
	for _, item := range allItems {
		totalLikes += item.LikesCount
	}

	// トップ10を取得
	topCount := min(10, len(allItems))

	return &notifier.SlackMessage{}, nil
}
