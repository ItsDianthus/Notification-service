FROM golang:1.24-alpine AS builder
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /bin/bot ./cmd/bot

FROM alpine:latest AS runtime

COPY --from=builder /bin/bot /bin/bot

CMD ["/bin/bot"]