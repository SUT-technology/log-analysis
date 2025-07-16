# مرحله Build
FROM golang:1.23.2 AS builder

ENV GOPROXY=https://goproxy.io,direct

WORKDIR /app

COPY go.mod go.sum ./
# RUN go mod download

COPY . .

RUN go build -o /app/serve ./cmd/logAnalysis/main.go

# مرحله اجرای نهایی
FROM debian:bookworm-slim

WORKDIR /app

# باینری و فایل کانفیگ را کپی کن
COPY --from=builder /app/serve /usr/local/bin/
COPY assets ./assets

# اجرای اپ با config دیفالت
ENTRYPOINT ["/usr/local/bin/serve"]
