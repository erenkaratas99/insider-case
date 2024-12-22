# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o insider ./main.go

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/insider .

COPY --from=builder /app/configs ./configs

EXPOSE 3001
CMD ["./insider", "jobScheduler"]
