# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN swag init -g ./cmd/jobScheduler.go -o ./docs/jobScheduler -instanceName jobScheduler --exclude ./internal/apps/messengerApi

RUN go build -o insider ./main.go

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/insider .

COPY --from=builder /app/configs ./configs

COPY --from=builder /app/docs/jobScheduler /docs/jobScheduler

EXPOSE 3001

CMD ["./insider", "jobScheduler"]