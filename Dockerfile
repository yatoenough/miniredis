FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o mini-redis ./cmd/miniredis

FROM alpine:latest

WORKDIR /app

RUN mkdir -p /data

COPY --from=builder /app/mini-redis .

EXPOSE 6379

WORKDIR /data

CMD ["/app/mini-redis"]
