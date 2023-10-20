# Compile api
FROM golang:1.18-alpine AS builder
WORKDIR /app

RUN apk add --no-cache make gcc musl-dev linux-headers git

COPY ./ ./
RUN go mod download
RUN go build ./cmd/crawler


# Copy compiled stuff and run it
FROM golang:1.18-alpine

COPY --from=builder /app/crawler /app/crawler

ENTRYPOINT ["/app/crawler"]
