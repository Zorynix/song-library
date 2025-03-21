FROM golang:1.23-alpine3.21 AS builder

WORKDIR /build

RUN apk update --no-cache && \
    apk add --no-cache gcc musl-dev tzdata  

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o /app/song-library ./cmd/app

FROM alpine:3.9.6

COPY --from=builder /usr/share/zoneinfo/Europe/Moscow /usr/share/zoneinfo/Europe/Moscow
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

WORKDIR /app
COPY --from=builder /app/song-library /app/song-library
COPY config/config.yaml /app/config/config.yaml
COPY migrations /app/migrations

ENV TZ=Europe/Moscow
ENV APP_PORT=8080

EXPOSE 8080

CMD ["./song-library"]