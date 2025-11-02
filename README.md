# Qiita アドベントカレンダーいいね集計 Slack Bot ツール

タグ毎の記事一覧を取得する API
https://qiita.com/api/v2/docs#get-apiv2tagstag_iditems

## セットアップ

1. Slack Webhook URL を取得

   - https://api.slack.com/apps でアプリを作成
   - Incoming Webhooks を有効化して Webhook URL を取得

2. 環境変数を設定

```bash
export SLACK_WEBHOOK_URL="https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
export QIITA_ACCESS_TOKEN="your-qiita-token"
```

3. 依存関係をインストール

```bash
go mod tidy
```

4. 実行

```bash
# シンプルなメッセージ送信
go run main.go -mode=message -message="Hello, Slack!"

# Qiitaタグのいいね集計
go run main.go -mode=aggregate -tag=Go -pages=5

# RSSフィードから最新記事取得
go run main.go -mode=rss -rss="https://qiita.com/tags/Go/feed" -limit=10
```

## 機能

### 1. 集計モード (`-mode=aggregate`)

- 指定した Qiita タグの記事を取得
- 記事のいいね数を集計
- トップ 10 記事を Slack に投稿
- 総記事数と総いいね数を表示

### 2. RSS モード (`-mode=rss`)

- RSS フィードから最新記事を取得
- 指定件数の記事を Slack に投稿

## オプション

```
-mode string
    動作モード: message, aggregate, rss (デフォルト: "message")
-tag string
    Qiitaのタグ名（aggregateモード用） (デフォルト: "Go")
-pages int
    取得する最大ページ数（aggregateモード用） (デフォルト: 5)
-rss string
    RSSフィードのURL（rssモード用）
-limit int
    取得する記事数（rssモード用） (デフォルト: 10)
-message string
    送信するメッセージ（messageモード用）
-help
    ヘルプを表示
```
