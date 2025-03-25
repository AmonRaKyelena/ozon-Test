FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ozon-test-project ./cmd

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/ozon-test-project .
COPY config.json .

EXPOSE 8080

ENTRYPOINT ["./ozon-test-project"]
