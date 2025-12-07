FROM golang:1.20-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go env -w GOPROXY=https://proxy.golang.org
RUN go mod download
COPY . .
RUN go build -o /client ./...

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=build /client /client
ENTRYPOINT ["/client"]

