FROM ubuntu:latest
LABEL authors="Dianazavr"
# build/docker/bot/Dockerfile

# --- stage: build binary ---
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download \
  && go mod tidy
RUN go mod download

COPY . .
WORKDIR /app/cmd/bot
RUN CGO_ENABLED=0 go build -o /bin/bot

# --- stage: final image ---
FROM alpine:3.17
RUN apk update && \
    apk add --no-cache ca-certificates && \
    rm -rf /var/cache/apk/*
COPY --from=builder /bin/bot /bin/bot
COPY config/bot.yaml /etc/bot.yaml
ENTRYPOINT ["/bin/bot", "--config", "/etc/bot.yaml"]