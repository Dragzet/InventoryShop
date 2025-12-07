# Multi-stage build for inventory service
FROM golang:1.20-alpine AS build
WORKDIR /src
# Скопировать только модули для более быстрого кэширования
COPY go.mod ./
RUN go env -w GOPROXY=https://proxy.golang.org
RUN go mod download
COPY . .
RUN go build -o /inventory ./...

# Final minimal image
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=build /inventory /inventory
EXPOSE 8001
ENV INVENTORY_PORT=8001
ENTRYPOINT ["/inventory"]

