# アーキテクチャ

[English](architecture.md) | [日本語](architecture_ja.md)

## 概要

goreload は、関心の分離が明確なモジュラーアーキテクチャに従っています。各コンポーネントは単一の責任を持ち、明確に定義されたインターフェースを通じて通信します。

## ディレクトリ構造

```
goreload/
├── cmd/goreload/
│   └── main.go              # CLI エントリーポイント、シグナル処理
├── internal/
│   ├── config/
│   │   ├── config.go        # 設定構造体、バリデーション
│   │   └── loader.go        # YAML 読み込み、デフォルト値
│   ├── logger/
│   │   └── logger.go        # 構造化ログ、カラー出力
│   ├── builder/
│   │   └── builder.go       # ビルドコマンド実行
│   ├── runner/
│   │   └── runner.go        # プロセスライフサイクル管理
│   ├── watcher/
│   │   ├── watcher.go       # ファイルシステム監視
│   │   └── filter.go        # パス/拡張子フィルタリング
│   └── engine/
│       └── engine.go        # オーケストレーション、メインループ
├── docs/                    # ドキュメント
├── .claude/                 # Claude Code 設定
└── .github/workflows/       # CI/CD
```

## コンポーネント図

```
┌─────────────────────────────────────────────────────────────┐
│                         CLI (main.go)                       │
│  - フラグのパース                                           │
│  - シグナル処理 (SIGINT, SIGTERM)                           │
│  - コンテキスト管理                                         │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Engine (engine.go)                     │
│  - watch-build-run サイクルのオーケストレーション           │
│  - 全コンポーネントの調整                                   │
│  - メインイベントループ                                     │
└─────────────────────────────────────────────────────────────┘
          │              │              │              │
          ▼              ▼              ▼              ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│    Config    │ │   Watcher    │ │   Builder    │ │    Runner    │
│              │ │              │ │              │ │              │
│ - YAML読込   │ │ - fsnotify   │ │ - ビルド     │ │ - 開始       │
│ - 検証       │ │ - デバウンス │ │   実行       │ │ - 停止       │
│ - デフォルト │ │ - フィルタ   │ │ - 出力       │ │ - 再起動     │
│              │ │              │ │   キャプチャ │ │ - SIGINT/    │
│              │ │              │ │              │ │   SIGKILL    │
└──────────────┘ └──────────────┘ └──────────────┘ └──────────────┘
                       │
                       ▼
               ┌──────────────┐
               │    Filter    │
               │              │
               │ - 拡張子     │
               │ - Globマッチ │
               │ - 除外       │
               └──────────────┘
```

## コアコンポーネント

### Config (`internal/config/`)

**責任:** 設定の読み込み、検証、提供。

**主要な型:**

```go
type Config struct {
    Root   string      // プロジェクトルートディレクトリ
    TmpDir string      // ビルド成果物ディレクトリ
    Build  BuildConfig // ビルド設定
    Watch  WatchConfig // 監視設定
    Log    LogConfig   // ログ設定
}

type BuildConfig struct {
    Cmd       string        // ビルドコマンド
    Bin       string        // バイナリパス
    Args      []string      // 実行時引数
    Delay     time.Duration // デバウンス遅延
    KillDelay time.Duration // シャットダウン猶予期間
}
```

**主要な関数:**

| 関数 | 説明 |
|----------|-------------|
| `LoadWithDefaults(path)` | YAML を読み込みデフォルトとマージ |
| `Default()` | デフォルト設定を取得 |
| `Validate()` | 設定値を検証 |
| `WriteDefault(path)` | デフォルト設定ファイルを生成 |

### Logger (`internal/logger/`)

**責任:** カラーとレベルをサポートする構造化ログ。

**インターフェース:**

```go
type Logger interface {
    Debug(msg string, args ...any)
    Info(msg string, args ...any)
    Warn(msg string, args ...any)
    Error(msg string, args ...any)
    SetOutput(w io.Writer)
    SetLevel(level Level)
}
```

**機能:**

- カラーコード出力 (設定可能)
- タイムスタンププレフィックス (設定可能)
- ログレベルフィルタリング
- スレッドセーフな出力

### Builder (`internal/builder/`)

**責任:** ビルドコマンドの実行と成果物の管理。

**インターフェース:**

```go
type Builder interface {
    Build(ctx context.Context) Result
    Clean() error
}

type Result struct {
    Success  bool
    Output   string
    Duration time.Duration
    Error    error
}
```

**機能:**

- コンテキスト対応 (キャンセル可能なビルド)
- Stdout/stderr キャプチャ
- 所要時間追跡
- 自動的な一時ディレクトリ作成

### Runner (`internal/runner/`)

**責任:** アプリケーションプロセスのライフサイクル管理。

**インターフェース:**

```go
type Runner interface {
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    Restart(ctx context.Context) error
    Running() bool
}
```

**グレースフルシャットダウン:**

1. プロセスグループに `SIGINT` を送信
2. `KillDelay` 待機
3. まだ実行中の場合は `SIGKILL` を送信
4. プロセスリソースのクリーンアップ

**プロセスグループ:**

`Setpgid: true` を使用してプロセスグループを作成し、子プロセスも確実に終了させます。

### Watcher (`internal/watcher/`)

**責任:** ファイルシステムの変更監視。

**インターフェース:**

```go
type Watcher interface {
    Start(ctx context.Context) error
    Events() <-chan Event
    Errors() <-chan error
    Close() error
}

type Event struct {
    Path string
    Op   Op
    Time time.Time
}
```

**機能:**

- 再帰的なディレクトリ監視
- デバウンス (頻繁な変更の結合)
- 新規ディレクトリの自動検出
- 除外パターン

### Filter (`internal/watcher/`)

**責任:** どのファイルが再ビルドをトリガーするかを決定。

**インターフェース:**

```go
type Filter interface {
    Match(path string) bool
}
```

**フィルタリングロジック:**

1. パスが除外ディレクトリにあるか確認
2. ファイル名が除外パターン (glob) に一致するか確認
3. 拡張子が監視リストにあるか確認
4. すべてのチェックを通過した場合に true を返す

### Engine (`internal/engine/`)

**責任:** watch-build-run サイクルのオーケストレーション。

**メインループ:**

```
┌─────────────────────────────────────────────┐
│                   開始                       │
└─────────────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────┐
│            初回ビルド & 実行                 │
└─────────────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────┐
│              イベント待機                    │◄────┐
│  - ファイル変更イベント                      │     │
│  - エラーイベント                            │     │
│  - コンテキストキャンセル                    │     │
└─────────────────────────────────────────────┘     │
                      │                              │
          ┌──────────┴──────────┐                   │
          ▼                     ▼                   │
   ┌────────────┐        ┌────────────┐            │
   │ ファイル   │        │ シャット   │            │
   │ イベント   │        │ ダウン     │            │
   └────────────┘        └────────────┘            │
          │                     │                   │
          ▼                     ▼                   │
   ┌────────────┐        ┌────────────┐            │
   │ プロセス   │        │ プロセス   │            │
   │ 停止       │        │ 停止       │            │
   │ ビルド     │        │ クリーン   │            │
   │ プロセス   │        │ アップ     │            │
   │ 開始       │        │ 終了       │            │
   └────────────┘        └────────────┘            │
          │                                         │
          └─────────────────────────────────────────┘
```

## データフロー

### 設定フロー

```
goreload.yaml → LoadWithDefaults() → Validate() → Config struct
                      ↑
               デフォルト値
```

### イベントフロー

```
ファイルシステム変更
       │
       ▼
   fsnotify
       │
       ▼
   Filter.Match()
       │
   ┌───┴───┐
   │ false │ → (無視)
   └───────┘
       │ true
       ▼
   デバウンスタイマー
       │
       ▼
   Engine イベントチャネル
       │
       ▼
   停止 → ビルド → 開始
```

## 並行処理モデル

### ゴルーチン

| ゴルーチン | 所有者 | 目的 |
|-----------|-------|---------|
| Main | CLI | シグナル処理、コンテキスト |
| Watcher loop | Watcher | fsnotify イベント処理 |
| Process wait | Runner | プロセス終了待機 |

### 同期

| コンポーネント | メカニズム | 目的 |
|-----------|-----------|---------|
| Logger | `sync.Mutex` | スレッドセーフな出力 |
| Runner | `sync.Mutex` | プロセス状態保護 |
| Watcher | Channels | イベント通信 |
| Engine | `sync.Mutex` | 実行状態 |

### コンテキストの使用

すべての長時間実行操作はキャンセルのために `context.Context` を受け取ります:

```go
func (e *Engine) Run(ctx context.Context) error
func (b *Builder) Build(ctx context.Context) Result
func (r *Runner) Start(ctx context.Context) error
func (r *Runner) Stop(ctx context.Context) error
func (w *Watcher) Start(ctx context.Context) error
```

## エラーハンドリング

### エラータイプ

| タイプ | 例 | 処理 |
|------|---------|----------|
| 設定 | 無効な YAML | 致命的、終了 |
| ビルド | コンパイルエラー | ログ出力、監視継続 |
| ランタイム | プロセスクラッシュ | ログ出力、次の変更を待機 |
| システム | ファイル権限 | エラーログ出力、継続 |

### センチネルエラー

```go
var (
    ErrEmptyBuildCmd    = errors.New("build command cannot be empty")
    ErrEmptyBin         = errors.New("binary path cannot be empty")
    ErrInvalidDelay     = errors.New("delay must be positive")
    ErrInvalidKillDelay = errors.New("kill_delay must be positive")
    ErrInvalidLogLevel  = errors.New("invalid log level")
    ErrNoExtensions     = errors.New("at least one extension required")
    ErrNoDirs           = errors.New("at least one directory required")
)
```

## 依存関係

| パッケージ | 目的 |
|---------|---------|
| `github.com/fsnotify/fsnotify` | クロスプラットフォームファイル監視 |
| `gopkg.in/yaml.v3` | YAML 設定パース |
| `github.com/spf13/cobra` | CLI フレームワーク |
| `github.com/fatih/color` | ターミナルカラー出力 |
