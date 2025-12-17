# CLI リファレンス

[English](cli.md) | [日本語](cli_ja.md)

## 概要

```
goreload [flags]
goreload [command]
```

## グローバルフラグ

| フラグ | 短縮形 | デフォルト | 説明 |
|------|-------|---------|-------------|
| `--config` | `-c` | `goreload.yaml` | 設定ファイルへのパス |
| `--help` | `-h` | | コマンドのヘルプを表示 |

## コマンド

### `goreload` (デフォルト)

ホットリロード監視を実行します。

```bash
goreload
goreload -c ./custom-config.yaml
```

**動作:**

1. `goreload.yaml`（または指定されたファイル）から設定を読み込みます
2. 必要に応じて一時ディレクトリを作成します
3. 初回ビルドを実行します
4. コンパイルされたバイナリを起動します
5. ファイルの変更を監視します
6. 変更時: プロセス停止 → 再ビルド → 再起動
7. 中断されるまで継続します (Ctrl+C)

**終了コード:**

| コード | 説明 |
|------|-------------|
| 0 | 正常終了 (Ctrl+C) |
| 1 | 設定エラーまたは致命的なエラー |

### `goreload init`

デフォルトの設定ファイルを生成します。

```bash
goreload init
```

**動作:**

- カレントディレクトリに `goreload.yaml` を作成します
- ファイルが既に存在する場合は失敗します（誤って上書きするのを防ぐため）

**出力:**

```
Created goreload.yaml
```

### `goreload version`

バージョン情報を表示します。

```bash
goreload version
```

**出力:**

```
goreload v0.1.0
  commit: abc1234
  built:  2024-01-01T00:00:00Z
```

### `goreload help`

ヘルプ情報を表示します。

```bash
goreload help
goreload help init
goreload --help
```

## 使用例

### 基本的な使用法

```bash
# デフォルト設定を使用
goreload

# カスタム設定を使用
goreload -c ./config/dev.yaml
```

### Web アプリケーション

```yaml
# goreload.yaml
build:
  cmd: "go build -o ./tmp/server ./cmd/server"
  bin: "./tmp/server"
  args:
    - "-port=8080"
    - "-env=development"

watch:
  extensions:
    - ".go"
    - ".html"
    - ".css"
  dirs:
    - "."
    - "templates"
  exclude_dirs:
    - "tmp"
    - "node_modules"
```

### Make を使用したマイクロサービス

```yaml
# goreload.yaml
build:
  cmd: "make build"
  bin: "./bin/service"

watch:
  extensions:
    - ".go"
    - ".proto"
  exclude_files:
    - "*.pb.go"
```

### デバッグモード

```yaml
# goreload.yaml
build:
  cmd: "go build -gcflags='all=-N -l' -o ./tmp/debug ."
  bin: "./tmp/debug"

log:
  level: "debug"
```

## シグナル処理

goreload は以下のシグナルを処理します:

| シグナル | 動作 |
|--------|----------|
| `SIGINT` (Ctrl+C) | グレースフルシャットダウン |
| `SIGTERM` | グレースフルシャットダウン |

**グレースフルシャットダウンプロセス:**

1. ファイル変更の監視を停止
2. 実行中のプロセスに SIGINT を送信
3. `kill_delay` の期間待機
4. プロセスがまだ実行中の場合は SIGKILL を送信
5. 終了

## 出力フォーマット

### ログメッセージ

```
[TIMESTAMP] [LEVEL] message
```

例:

```
15:04:05 [INFO] watching: /path/to/project
15:04:05 [INFO] excluding: [tmp vendor .git]
15:04:05 [INFO] building...
15:04:07 [INFO] ✓ build completed (2.15s)
15:04:07 [INFO] ✓ running ./tmp/main
```

### ログレベル

| レベル | 色 | 説明 |
|-------|-------|-------------|
| `[DEBUG]` | グレー | 詳細なデバッグ情報 |
| `[INFO]` | シアン | 通常の操作メッセージ |
| `[WARN]` | 黄色 | 警告 |
| `[ERROR]` | 赤 | エラー |

### ステータスインジケーター

| インジケーター | 意味 |
|-----------|---------|
| `✓` | 成功 (ビルド完了、プロセス開始) |
| `✗` | 失敗 (ビルド失敗、プロセスエラー) |
