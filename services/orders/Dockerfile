# Multi-stage build for orders service
FROM golang:1.20-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go env -w GOPROXY=https://proxy.golang.org
RUN go mod download
COPY . .
RUN go build -o /orders ./...

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=build /orders /orders
EXPOSE 8002
ENV ORDERS_PORT=8002
# default INVENTORY_URL assumes inventory is reachable at http://inventory:8001
ENV INVENTORY_URL=http://inventory:8001
ENTRYPOINT ["/orders"]

