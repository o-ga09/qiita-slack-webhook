package aggregatetags

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type QiitaItem struct {
	Title      string `json:"title"`
	URL        string `json:"url"`
	LikesCount int    `json:"likes_count"`
	User       struct {
		ID string `json:"id"`
	} `json:"user"`
}

// fetchQiitaItemsByTag はタグに基づいて記事を取得
func fetchQiitaItemsByTag(tag string, perPage, page int) ([]QiitaItem, error) {
	url := fmt.Sprintf("https://qiita.com/api/v2/tags/%s/items?page=%d&per_page=%d", tag, page, perPage)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Qiita APIトークンがあれば設定（オプション）
	if token := os.Getenv("QIITA_ACCESS_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var items []QiitaItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}

	return items, nil
}
