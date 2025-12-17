# goreload

[![CI](https://github.com/taro33333/goreload/actions/workflows/ci.yaml/badge.svg)](https://github.com/taro33333/goreload/actions/workflows/ci.yaml)
[![Release](https://github.com/taro33333/goreload/actions/workflows/release.yaml/badge.svg)](https://github.com/taro33333/goreload/actions/workflows/release.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/taro33333/goreload)](https://goreportcard.com/report/github.com/taro33333/goreload)

[English](README.md) | [日本語](README_ja.md)

Goアプリケーション用のホットリロードツールです。ソースファイルを監視し、変更が検知されると自動的にアプリケーションを再ビルドして再起動します。

## 機能

- 設定可能な拡張子によるファイル監視
- 頻繁な再ビルドを防ぐデバウンスビルド
- 設定可能なタイムアウトによるグレースフルなプロセス終了
- カラーログ出力
- ファイル除外のためのGlobパターンサポート
- 再帰的なディレクトリ監視
- クロスプラットフォームサポート (Linux, macOS, Windows)

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

## 使い方

### クイックスタート

```bash
# 設定ファイルの初期化
goreload init

# デフォルト設定で実行
goreload
```

### 設定ファイルの指定

```bash
goreload -c ./config.yaml
```

### バージョンの表示

```bash
goreload version
```

## 設定

プロジェクトルートに `goreload.yaml` ファイルを作成します:

```yaml
# プロジェクトルートディレクトリ
root: "."

# ビルド成果物用の一時ディレクトリ
tmp_dir: "tmp"

# ビルド設定
build:
  # ビルドコマンド
  cmd: "go build -o ./tmp/main ."
  # 実行するバイナリ
  bin: "./tmp/main"
  # バイナリに渡す引数
  args: []
  # ファイル変更後、ビルドまでの遅延時間 (デバウンス)
  delay: "200ms"
  # プロセス終了の猶予時間
  kill_delay: "500ms"

# ファイル監視設定
watch:
  # 監視するファイル拡張子
  extensions:
    - ".go"
  # 監視するディレクトリ
  dirs:
    - "."
  # 除外するディレクトリ
  exclude_dirs:
    - "tmp"
    - "vendor"
    - ".git"
    - "node_modules"
  # 除外するファイル (Globパターン)
  exclude_files:
    - "*_test.go"

# ログ設定
log:
  # カラー出力を有効化
  color: true
  # タイムスタンプを表示
  time: true
  # ログレベル: debug, info, warn, error
  level: "info"
```

## 出力例

```
   __ _  ___  _ __ ___| | ___   __ _  __| |
  / _` |/ _ \| '__/ _ \ |/ _ \ / _` |/ _` |
 | (_| | (_) | | |  __/ | (_) | (_| | (_| |
  \__, |\___/|_|  \___|_|\___/ \__,_|\__,_|
  |___/                            v0.1.2

[INFO] watching: .
[INFO] excluding: [tmp vendor .git node_modules]
[INFO] building...
[INFO] ✓ build completed (1.23s)
[INFO] ✓ running ./tmp/main

[INFO] main.go changed
[INFO] building...
[INFO] ✓ build completed (0.45s)
[INFO] ✓ running ./tmp/main
```

## アーキテクチャ

```
goreload/
├── cmd/goreload/main.go     # CLI エントリーポイント
├── internal/
│   ├── config/              # 設定の読み込みとバリデーション
│   ├── watcher/             # ファイルシステム監視
│   ├── builder/             # ビルドコマンド実行
│   ├── runner/              # プロセス管理
│   ├── engine/              # オーケストレーター
│   └── logger/              # 構造化ログ
├── goreload.yaml            # 設定サンプル
└── README.md
```

## 依存関係

- [fsnotify](https://github.com/fsnotify/fsnotify) - ファイルシステム通知
- [yaml.v3](https://gopkg.in/yaml.v3) - YAML パース
- [cobra](https://github.com/spf13/cobra) - CLI フレームワーク
- [color](https://github.com/fatih/color) - カラーターミナル出力

## 貢献

貢献は大歓迎です！プルリクエストを送ってください。

## ライセンス

MIT
