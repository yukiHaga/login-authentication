# デプロイ用コンテナに含めるバイナリを作成するコンテナ
FROM golang:1.18.2-bullseye AS deploy-builder

WORKDIR /app

COPY go.mod go.sum ./

# go mod downloadは、引数にモジュールの指定がない場合は、go.modファイルに記載されたすべてのモジュールをダウンロードする
RUN go mod download

# アプリケーションのコードをコピー
COPY . .

WORKDIR /app/src

RUN go build -o main

# ローカル開発環境で利用するホットリロード環境
FROM golang:1.18.2 AS dev
# go installは指定したパッケージをダウンロード後にビルドし、実行可能なファイルを$GOBINへ格納する
RUN go install github.com/cosmtrek/air@latest

WORKDIR /app/src

EXPOSE 8080

# air -c [tomlファイル名] // 設定ファイルを指定してair実行(WORKDIRに.air.tomlを配置しておくこと)
# -c .ari.toml のオプションは省略可能。その場合、air はカレントディレクトリから .air.toml ファイルを探し起動します。
CMD ["air"]