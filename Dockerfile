# Compile api
FROM golang:1.18-alpine AS builder
WORKDIR /app

RUN apk add --no-cache gcc musl-dev linux-headers

COPY ./ ./
RUN go mod download
RUN go build ./cmd/crawler


# Copy compiled stuff and run it
FROM golang:1.18-alpine

COPY --from=builder /app/crawler /app/crawler

ENTRYPOINT ["/app/crawler"]
