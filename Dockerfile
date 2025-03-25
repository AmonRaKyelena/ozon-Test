FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o github.com/AmonRaKyelena/ozon-Test ./cmd

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/github.com/AmonRaKyelena/ozon-Test .
COPY config.json .

EXPOSE 8080

ENTRYPOINT ["./github.com/AmonRaKyelena/ozon-Test"]
