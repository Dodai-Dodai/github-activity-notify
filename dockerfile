# ビルドステージ
FROM golang:latest AS builder

WORKDIR /work

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main main.go && \
    go clean -modcache

# ランタイムステージ
FROM gcr.io/distroless/static-debian12:latest

WORKDIR /app

COPY --from=builder /work/main /app/

ENTRYPOINT ["./main"]