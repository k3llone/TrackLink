FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /tracklink-api ./cmd/api

FROM alpine:3.21

WORKDIR /srv
COPY --from=builder /tracklink-api /usr/local/bin/tracklink-api

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/tracklink-api"]
