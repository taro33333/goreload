# Docker サポート

[English](docker.md) | [日本語](docker_ja.md)

Docker コンテナ内で goreload を使用することで、コンテナ化された開発環境でホットリロードを有効にすることができます。

## Dockerfile のセットアップ

プロジェクトで goreload を使用するには、`Dockerfile` 内で `go install` を使用してインストールします。

```dockerfile
FROM golang:1.25

# goreload のインストール
RUN go install github.com/taro33333/goreload/cmd/goreload@latest

WORKDIR /app

# go モジュールファイルのコピー
COPY go.mod go.sum ./
RUN go mod download

# ソースコードのコピー
COPY . .

# エントリーポイントとして goreload を設定
CMD ["goreload"]
```

## Docker Compose のセットアップ

ホットリロードを機能させるには、ローカルのソースコードディレクトリをコンテナに**マウントする必要があります**。これにより、goreload はホストマシン上で行われたファイル変更を検知できます。

```yaml
version: '3.8'

services:
  app:
    build: .
    # カレントディレクトリをコンテナ内の /app にマウント
    volumes:
      - .:/app
    # 必要に応じてポートを公開
    ports:
      - "8080:8080"
```

## 設定のヒント

### ビルド出力

デフォルトでは、goreload はバイナリを `./tmp/main` にビルドします。コンテナユーザーがこのディレクトリへの書き込み権限を持っていることを確認してください。

### Linux/Docker でのファイル監視

goreload は `fsnotify` を使用しており、Linux（したがって Docker コンテナ内）では `inotify` に依存しています。これは Docker for Mac/Windows のファイル共有機能とシームレスに連携します。

### マルチステージビルド

本番ビルドでは通常 goreload は不要です。マルチステージビルドを使用して、開発環境と本番イメージを分離することができます。

```dockerfile
# 開発用ステージ
FROM golang:1.25 AS dev
RUN go install github.com/taro33333/goreload/cmd/goreload@latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
CMD ["goreload"]

# 本番用ステージ
FROM golang:1.25 AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```
