# goreload - Go Hot Reload Tool

## Project Overview

goreload は Go アプリケーション用のホットリロードツールです。Air ライクな機能を提供し、ファイル変更を検知して自動でビルド・再起動を行います。

## Documentation

詳細なドキュメントは `docs/` ディレクトリを参照:

- [docs/README.md](docs/README.md) - ドキュメント目次
- [docs/getting-started.md](docs/getting-started.md) - インストールとクイックスタート
- [docs/configuration.md](docs/configuration.md) - 設定リファレンス（全オプション）
- [docs/cli.md](docs/cli.md) - CLI リファレンス
- [docs/architecture.md](docs/architecture.md) - 内部設計とコンポーネント
- [docs/development.md](docs/development.md) - 開発ガイド

## IMPORTANT: Documentation Sync Rules

**仕様変更時は必ず関連ドキュメントを更新すること！**

| 変更内容 | 更新が必要なファイル |
|---------|---------------------|
| 設定オプションの追加/変更/削除 | `docs/configuration.md`, `internal/config/config.go` |
| CLI コマンド/フラグの追加/変更 | `docs/cli.md`, `cmd/goreload/main.go` |
| アーキテクチャ/インターフェース変更 | `docs/architecture.md`, `CLAUDE.md` |
| インストール方法の変更 | `docs/getting-started.md`, `README.md` |
| 新機能追加 | 関連する全ドキュメント |

### ドキュメント更新チェックリスト

コード変更時に以下を確認:

1. [ ] 設定構造体 (`Config`) を変更した → `docs/configuration.md` を更新
2. [ ] CLI コマンドを追加/変更した → `docs/cli.md` を更新
3. [ ] インターフェースを変更した → `docs/architecture.md` を更新
4. [ ] デフォルト値を変更した → `docs/configuration.md` を更新
5. [ ] エラーメッセージを変更した → 関連ドキュメントを更新

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
├── docs/                    # ドキュメント
├── .claude/                 # Claude Code 設定
├── goreload.yaml            # サンプル設定ファイル
└── go.mod
```

## Quick Commands

```bash
# ビルド
go build -o ./tmp/goreload ./cmd/goreload

# テスト
go test ./... -v -cover

# 静的解析
go vet ./...

# 実行
./tmp/goreload
./tmp/goreload init
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

## Configuration Structure

現在の設定オプション（詳細は `docs/configuration.md` 参照）:

```yaml
root: "."           # プロジェクトルート
tmp_dir: "tmp"      # ビルド成果物ディレクトリ

build:
  cmd: "go build -o ./tmp/main ."  # ビルドコマンド
  bin: "./tmp/main"                 # 実行バイナリ
  args: []                          # 実行時引数
  delay: "200ms"                    # デバウンス遅延
  kill_delay: "500ms"               # 終了猶予時間

watch:
  extensions: [".go"]               # 監視拡張子
  dirs: ["."]                       # 監視ディレクトリ
  exclude_dirs: ["tmp", "vendor"]   # 除外ディレクトリ
  exclude_files: []                 # 除外ファイル（glob）

log:
  color: true                       # カラー出力
  time: true                        # タイムスタンプ
  level: "info"                     # ログレベル
```

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
