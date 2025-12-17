# 開発ガイド

[English](development.md) | [日本語](development_ja.md)

## 前提条件

- Go 1.21 以上
- Git
- (オプション) リント用 staticcheck

## はじめに

### リポジトリのクローン

```bash
git clone https://github.com/taro33333/goreload.git
cd goreload
```

### 依存関係のインストール

```bash
go mod download
```

### ビルド

```bash
go build -o ./tmp/goreload ./cmd/goreload
```

### テスト実行

```bash
# 全テスト
go test ./...

# 詳細出力付き
go test -v ./...

# カバレッジ付き
go test -cover ./...

# カバレッジレポート生成
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### リント

```bash
# go vet
go vet ./...

# staticcheck (インストール済みの場合)
staticcheck ./...
```

## プロジェクト構造

```
goreload/
├── cmd/goreload/          # CLI エントリーポイント
├── internal/              # プライベートパッケージ
│   ├── config/            # 設定
│   ├── logger/            # ログ
│   ├── builder/           # ビルド実行
│   ├── runner/            # プロセス管理
│   ├── watcher/           # ファイル監視
│   └── engine/            # オーケストレーション
├── docs/                  # ドキュメント
├── .claude/               # Claude Code 設定
├── .github/workflows/     # CI/CD
├── goreload.yaml          # 設定サンプル
├── .goreleaser.yaml       # リリース設定
└── go.mod                 # Go モジュール
```

## コーディング規約

### Go イディオム

- [Effective Go](https://golang.org/doc/effective_go) に従う
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) に従う
- フォーマットには `gofmt` または `goimports` を使用する

### 命名規則

| タイプ | 規約 | 例 |
|------|------------|---------|
| パッケージ | 小文字、単語 | `config`, `logger` |
| エクスポート | PascalCase | `LoadConfig`, `Builder` |
| 非エクスポート | camelCase | `parseCommand`, `runLoop` |
| インターフェース | -er 接尾辞 | `Builder`, `Runner`, `Watcher` |
| レシーバ | 1-2 文字 | `func (w *Watcher)` |

### エラーハンドリング

```go
// コンテキスト付きでエラーをラップ
if err != nil {
    return fmt.Errorf("load config: %w", err)
}

// センチネルエラーの定義
var ErrNotFound = errors.New("not found")

// 特定のエラーのチェック
if errors.Is(err, ErrNotFound) {
    // 処理
}
```

### ドキュメント

```go
// Package config provides configuration loading and validation.
package config

// Config represents the complete goreload configuration.
type Config struct {
    // Root is the project root directory.
    Root string `yaml:"root"`
}

// LoadWithDefaults loads configuration from path and merges with defaults.
func LoadWithDefaults(path string) (*Config, error) {
    // ...
}
```

## テスト

### テーブル駆動テスト

```go
func TestConfig_Validate(t *testing.T) {
    tests := []struct {
        name    string
        cfg     Config
        wantErr error
    }{
        {
            name:    "valid config",
            cfg:     validConfig(),
            wantErr: nil,
        },
        {
            name: "empty build cmd",
            cfg: Config{
                Build: BuildConfig{Cmd: ""},
            },
            wantErr: ErrEmptyBuildCmd,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.cfg.Validate()
            if !errors.Is(err, tt.wantErr) {
                t.Errorf("got %v, want %v", err, tt.wantErr)
            }
        })
    }
}
```

### テストヘルパー

```go
func TestBuilder_Build(t *testing.T) {
    // 一時ファイルには t.TempDir() を使用
    tmpDir := t.TempDir()

    // テストファイルの作成
    goFile := filepath.Join(tmpDir, "main.go")
    os.WriteFile(goFile, []byte("package main\nfunc main(){}"), 0644)

    // クリーンアップには t.Cleanup() を使用
    t.Cleanup(func() {
        // クリーンアップコード
    })
}
```

### カバレッジ目標

| パッケージ | 目標 |
|---------|--------|
| config | 90%+ |
| logger | 95%+ |
| builder | 85%+ |
| watcher | 80%+ |
| runner | 75%+ |
| engine | 60%+ |

## 新機能の追加

### 1. 設計

1. インターフェースを定義
2. 動作を文書化
3. エラーケースを考慮
4. テスト容易性を計画

### 2. 実装

1. 実装を作成
2. GoDoc コメントを追加
3. エラーを適切に処理
4. キャンセルにはコンテキストを使用

### 3. テスト

1. テーブル駆動テストを作成
2. エラーケースをテスト
3. エッジケースをテスト
4. カバレッジ目標を達成

### 4. ドキュメント

1. `docs/` 内の関連ドキュメントを更新
2. アーキテクチャが変更された場合は `CLAUDE.md` を更新
3. ユーザー向けの場合は README を更新

## リリースプロセス

### バージョンタグ付け

```bash
# タグを作成してプッシュ
git tag v0.1.0
git push origin v0.1.0
```

### 自動リリース

GitHub Actions が自動的に行います:

1. テストとリントを実行
2. 全プラットフォーム向けのバイナリをビルド
3. GitHub Release を作成
4. Homebrew tap に公開

### ローカルテスト

```bash
# goreleaser 設定を確認
goreleaser check

# スナップショットビルド (公開なし)
goreleaser release --snapshot --clean
```

## デバッグ

### デバッグログの有効化

```yaml
# goreload.yaml
log:
  level: "debug"
```

### デバッグシンボル付きビルド

```bash
go build -gcflags="all=-N -l" -o ./tmp/goreload ./cmd/goreload
```

### Delve の使用

```bash
# delve のインストール
go install github.com/go-delve/delve/cmd/dlv@latest

# デバッグ
dlv debug ./cmd/goreload -- -c goreload.yaml
```

## CI/CD

### GitHub Actions

| ワークフロー | トリガー | アクション |
|----------|---------|---------|
| `ci.yaml` | main への Push/PR | テスト、リント、ビルド |
| `release.yaml` | タグプッシュ (v*) | リリース、公開 |

### 必要なシークレット

| シークレット | 目的 |
|--------|---------|
| `GITHUB_TOKEN` | リリースアセット |
| `HOMEBREW_TAP_TOKEN` | Homebrew tap プッシュ |

## トラブルシューティング

### 一般的な問題

**モジュールエラーでビルドが失敗する:**

```bash
go mod tidy
```

**権限エラーでテストが失敗する:**

```bash
# テスト一時ディレクトリのファイル権限を確認
chmod +x ./tmp/test-binary
```

**監視が変更を検知しない:**

- `watch.dirs` にディレクトリが含まれているか確認
- `watch.exclude_dirs` で除外されていないか確認
- ファイル拡張子が `watch.extensions` に含まれているか確認

## リソース

- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [fsnotify Documentation](https://github.com/fsnotify/fsnotify)
- [Cobra Documentation](https://cobra.dev/)
