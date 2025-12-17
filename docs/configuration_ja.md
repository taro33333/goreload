# 設定リファレンス

[English](configuration.md) | [日本語](configuration_ja.md)

goreload は YAML 設定ファイル（デフォルト: `goreload.yaml`）を使用して動作を制御します。

## 設定ファイルの場所

goreload は以下の順序で設定ファイルを探します:

1. `-c` または `--config` フラグで指定されたパス
2. カレントディレクトリの `goreload.yaml`

## 完全な設定例

```yaml
# プロジェクトルートディレクトリ
root: "."

# ビルド成果物用の一時ディレクトリ
tmp_dir: "tmp"

# ビルド設定
build:
  cmd: "go build -o ./tmp/main ."
  bin: "./tmp/main"
  args: []
  delay: "200ms"
  kill_delay: "500ms"

# ファイル監視設定
watch:
  extensions:
    - ".go"
  dirs:
    - "."
  exclude_dirs:
    - "tmp"
    - "vendor"
    - ".git"
    - "node_modules"
  exclude_files:
    - "*_test.go"

# ログ設定
log:
  color: true
  time: true
  level: "info"
```

## 設定オプション

### ルートレベル

| オプション | 型 | デフォルト | 説明 |
|--------|------|---------|-------------|
| `root` | string | `"."` | プロジェクトルートディレクトリ。すべての相対パスはここから解決されます。 |
| `tmp_dir` | string | `"tmp"` | ビルド成果物用のディレクトリ。存在しない場合は自動的に作成されます。 |

### ビルド設定 (`build`)

| オプション | 型 | デフォルト | 説明 |
|--------|------|---------|-------------|
| `cmd` | string | `"go build -o ./tmp/main ."` | 実行するビルドコマンド。シェルスタイルのクォートをサポートします。 |
| `bin` | string | `"./tmp/main"` | 実行するコンパイル済みバイナリのパス。 |
| `args` | []string | `[]` | 実行時にバイナリに渡す引数。 |
| `delay` | duration | `"200ms"` | ファイル変更後、ビルドをトリガーするまでのデバウンス遅延時間。 |
| `kill_delay` | duration | `"500ms"` | SIGKILL 前のプロセス終了の猶予時間。 |

#### Duration フォーマット

Duration の値は Go の duration フォーマットをサポートします:

- `"100ms"` - 100ミリ秒
- `"1s"` - 1秒
- `"1m30s"` - 1分30秒

#### ビルドコマンドの例

```yaml
# 標準的な Go ビルド
build:
  cmd: "go build -o ./tmp/main ."

# ビルドタグ付き
build:
  cmd: "go build -tags=dev -o ./tmp/main ."

# ldflags 付き
build:
  cmd: "go build -ldflags='-s -w' -o ./tmp/main ."

# make を使用
build:
  cmd: "make build"
  bin: "./bin/app"
```

### 監視設定 (`watch`)

| オプション | 型 | デフォルト | 説明 |
|--------|------|---------|-------------|
| `extensions` | []string | `[".go"]` | 監視するファイル拡張子。少なくとも1つ必要です。 |
| `dirs` | []string | `["."]` | 再帰的に監視するディレクトリ。少なくとも1つ必要です。 |
| `exclude_dirs` | []string | `["tmp", "vendor", ".git", "node_modules"]` | 監視から除外するディレクトリ。 |
| `exclude_files` | []string | `[]` | 除外するファイルパターン（Globパターンをサポート）。 |

#### 拡張子のフォーマット

拡張子は先頭のドットを含める必要があります:

```yaml
watch:
  extensions:
    - ".go"
    - ".html"
    - ".tmpl"
```

#### Glob パターン

`exclude_files` オプションは Glob パターンをサポートします:

| パターン | 説明 |
|---------|-------------|
| `*_test.go` | すべてのテストファイル |
| `*.gen.go` | すべての生成ファイル |
| `mock_*.go` | すべてのモックファイル |
| `*.pb.go` | すべての protobuf 生成ファイル |

### ログ設定 (`log`)

| オプション | 型 | デフォルト | 説明 |
|--------|------|---------|-------------|
| `color` | bool | `true` | カラー出力を有効にします。 |
| `time` | bool | `true` | ログ出力にタイムスタンプを表示します。 |
| `level` | string | `"info"` | ログレベル: `debug`, `info`, `warn`, `error`。 |

#### ログレベル

| レベル | 説明 |
|-------|-------------|
| `debug` | デバッグ用の詳細出力 |
| `info` | 通常の操作メッセージ |
| `warn` | 動作を停止しない警告 |
| `error` | エラーのみ |

## 環境変数

goreload は標準的な Go 環境変数を尊重します:

| 変数 | 説明 |
|----------|-------------|
| `GO111MODULE` | Go モジュールモード |
| `GOOS` | ターゲット OS |
| `GOARCH` | ターゲットアーキテクチャ |
| `CGO_ENABLED` | CGO の有効化/無効化 |

## バリデーションルール

設定は読み込み時に検証されます。以下のルールが適用されます:

1. `build.cmd` - 空であってはなりません
2. `build.bin` - 空であってはなりません
3. `build.delay` - 負の値であってはなりません
4. `build.kill_delay` - 負の値であってはなりません
5. `watch.extensions` - 少なくとも1つの拡張子が必要です
6. `watch.dirs` - 少なくとも1つのディレクトリが必要です
7. `log.level` - 次のいずれかでなければなりません: `debug`, `info`, `warn`, `error`

## デフォルト設定

設定ファイルが存在しない場合、goreload は以下のデフォルト値を使用します:

```yaml
root: "."
tmp_dir: "tmp"
build:
  cmd: "go build -o ./tmp/main ."
  bin: "./tmp/main"
  args: []
  delay: "200ms"
  kill_delay: "500ms"
watch:
  extensions:
    - ".go"
  dirs:
    - "."
  exclude_dirs:
    - "tmp"
    - "vendor"
    - ".git"
    - "node_modules"
  exclude_files: []
log:
  color: true
  time: true
  level: "info"
```
