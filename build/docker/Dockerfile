# syntax=docker/dockerfile:1

FROM golang:1.20-alpine3.17 AS builder
WORKDIR /build
COPY . .
RUN go build -o /build/linux_amd64/app ./cmd/main.go

FROM alpine:3.17
WORKDIR /app
COPY --from=builder /build/linux_amd64/app ./
CMD ["./app"]