FROM ubuntu:latest
LABEL authors="Dianazavr"
# build/docker/scrapper/Dockerfile

# --- stage: build binary ---
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download \
  && go mod tidy
RUN go mod download

COPY . .
WORKDIR /app/cmd/scrapper
RUN CGO_ENABLED=0 go build -o /bin/scrapper

# --- stage: final image ---
FROM alpine:3.17
RUN apk update && \
    apk add --no-cache ca-certificates && \
    rm -rf /var/cache/apk/*
COPY --from=builder /bin/scrapper /bin/scrapper
COPY config/scrapper.yaml /etc/scrapper.yaml
ENTRYPOINT ["/bin/scrapper", "--config", "/etc/scrapper.yaml"]