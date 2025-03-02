FROM golang:1.24-alpine AS builder
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /bin/scrapper ./cmd/scrapper

FROM alpine:latest AS runtime

COPY --from=builder /bin/scrapper /bin/scrapper

CMD ["/bin/scrapper"]