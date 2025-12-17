# goreload - Go Hot Reload Tool

## Project Overview

goreload は Go アプリケーション用のホットリロードツールです。Air ライクな機能を提供し、ファイル変更を検知して自動でビルド・再起動を行います。

## Architecture

```
goreload/
├── cmd/goreload/main.go     # CLI エントリーポイント (Cobra)
├── internal/
│   ├── config/              # 設定読み込み・バリデーション
│   │   ├── config.go        # Config 構造体、Validate()
│   │   └── loader.go        # YAML 読み込み、デフォルト値
│   ├── logger/              # 構造化ログ出力
│   │   └── logger.go        # カラー出力、レベル制御
│   ├── builder/             # ビルドコマンド実行
│   │   └── builder.go       # Build(), Clean()
│   ├── runner/              # プロセス管理
│   │   └── runner.go        # Start(), Stop(), Restart()
│   ├── watcher/             # ファイル監視
│   │   ├── watcher.go       # fsnotify ラッパー、デバウンス
│   │   └── filter.go        # 拡張子・パスフィルタ
│   └── engine/              # オーケストレーター
│       └── engine.go        # watch-build-run サイクル
├── goreload.yaml            # サンプル設定ファイル
├── go.mod
└── go.sum
```

## Coding Standards

### Go イディオム

- Effective Go と Go Code Review Comments に準拠
- エクスポートされた型・関数には GoDoc コメント必須
- エラーは `fmt.Errorf("context: %w", err)` でラップ
- `context.Context` は第一引数
- 早期リターンでネストを浅く

### 命名規則

- 短く明確な変数名: `cfg`, `ctx`, `err`
- レシーバ名は 1-2 文字: `func (w *Watcher)`
- インターフェース名は `-er` サフィックス: `Builder`, `Runner`

### エラーハンドリング

- センチネルエラー: `var ErrXxx = errors.New("xxx")`
- panic は使わない（init での致命的エラー以外）

## Build & Test Commands

```bash
# ビルド
go build -o ./tmp/goreload ./cmd/goreload

# テスト実行
go test ./... -v

# テスト（カバレッジ付き）
go test ./... -cover

# 静的解析
go vet ./...
staticcheck ./...

# goreload 実行
./tmp/goreload

# 設定ファイル生成
./tmp/goreload init

# バージョン確認
./tmp/goreload version
```

## Key Interfaces

```go
// Builder - ビルドコマンド実行
type Builder interface {
    Build(ctx context.Context) Result
    Clean() error
}

// Runner - プロセス管理
type Runner interface {
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    Restart(ctx context.Context) error
    Running() bool
}

// Watcher - ファイル監視
type Watcher interface {
    Start(ctx context.Context) error
    Events() <-chan Event
    Errors() <-chan error
    Close() error
}

// Logger - ログ出力
type Logger interface {
    Debug(msg string, args ...any)
    Info(msg string, args ...any)
    Warn(msg string, args ...any)
    Error(msg string, args ...any)
}
```

## Important Patterns

### グレースフルシャットダウン (runner.go)

1. SIGINT をプロセスグループに送信
2. KillDelay 待機
3. SIGKILL で強制終了

### デバウンス (watcher.go)

- ファイル変更イベントを一定時間バッファリング
- 連続した変更を1回のビルドにまとめる

### 設定のマージ (loader.go)

- デフォルト値 → YAML ファイル → 検証

## Dependencies

- `github.com/fsnotify/fsnotify` - ファイル監視
- `gopkg.in/yaml.v3` - YAML パース
- `github.com/spf13/cobra` - CLI フレームワーク
- `github.com/fatih/color` - カラー出力

## Test Coverage Goals

| パッケージ | 目標 |
|-----------|------|
| config    | 90%+ |
| logger    | 95%+ |
| builder   | 85%+ |
| watcher   | 80%+ |
| runner    | 75%+ |
| engine    | 60%+ |

## Common Issues

### ビルドが失敗する場合

- `go.mod` がプロジェクトルートにあるか確認
- 依存関係: `go mod tidy`

### プロセスが終了しない場合

- `kill_delay` を増やす
- 子プロセスが SIGINT を適切に処理しているか確認

### ファイル変更が検知されない場合

- `watch.dirs` に対象ディレクトリが含まれているか
- `watch.exclude_dirs` で除外されていないか
- `watch.extensions` に拡張子が含まれているか
