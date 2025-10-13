# --- Dev Stage for Local --- #
FROM golang:1.22.4-bullseye as dev

ENV TZ=Asia/Tokyo
ENV APP_DIR=/usr/src/app

COPY ./filmogophery-api/ ${APP_DIR}
WORKDIR ${APP_DIR}

# ローカルでテストする時は、gotest を使う
RUN go install github.com/rakyll/gotest@latest
# ハンドラのテスト用にmock生成パッケージをインストール（Go 1.22.4 互換バージョンを指定）
RUN go install go.uber.org/mock/mockgen@v0.4.0
RUN go install go.uber.org/mock/mockgen/model@v0.4.0
# ホットリロードのために、air を使う
RUN go install github.com/air-verse/air@v1.52.3

RUN go mod tidy



# --- Build Stage --- #
FROM golang:1.22.4-bullseye as builder

ENV TZ=Asia/Tokyo
ENV CGO_ENABLED=0
ENV APP_DIR=/usr/src/app

COPY ./filmogophery-api/ ${APP_DIR}
WORKDIR ${APP_DIR}

RUN go mod download
RUN go build \
    -o ./bin/server \
    -ldflags '-s -w' \
    ./cmd/server/main.go



# --- Deploy Stage --- #
FROM gcr.io/distroless/static-debian11 as deploy

ENV TZ=Asia/Tokyo

COPY --from=builder /usr/src/app/cmd/server/bin/server /
CMD ["/server"]
