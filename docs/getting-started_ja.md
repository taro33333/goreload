# はじめに

[English](getting-started.md) | [日本語](getting-started_ja.md)

## インストール

### Homebrew (macOS/Linux)

```bash
brew tap taro33333/tap
brew install --cask goreload
```

### Go Install

```bash
go install github.com/taro33333/goreload/cmd/goreload@latest
```

### バイナリのダウンロード

[GitHub Releases](https://github.com/taro33333/goreload/releases) から最新のリリースをダウンロードしてください。

### ソースからビルド

```bash
git clone https://github.com/taro33333/goreload.git
cd goreload
go build -o goreload ./cmd/goreload
```

## クイックスタート

### 1. 設定の初期化

Goプロジェクトのディレクトリに移動して実行します:

```bash
goreload init
```

これにより、デフォルト設定の `goreload.yaml` 設定ファイルが作成されます。

### 2. 監視の開始

```bash
goreload
```

goreload は以下の動作を行います:

1. アプリケーションをビルド
2. コンパイルされたバイナリを起動
3. ファイルの変更を監視
4. 変更が検知されると自動的に再ビルドして再起動

### 3. 停止

`Ctrl+C` を押して goreload を停止します。

## プロジェクト構成例

```
myapp/
├── main.go
├── go.mod
├── go.sum
└── goreload.yaml
```

### 最小限の `goreload.yaml`

```yaml
root: "."
tmp_dir: "tmp"

build:
  cmd: "go build -o ./tmp/main ."
  bin: "./tmp/main"

watch:
  extensions:
    - ".go"
  dirs:
    - "."
  exclude_dirs:
    - "tmp"
    - "vendor"

log:
  level: "info"
```

## 次のステップ

- [設定リファレンス](./configuration_ja.md) - すべての設定オプションについて学ぶ
- [CLI リファレンス](./cli_ja.md) - コマンドラインオプションを調べる
