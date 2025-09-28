# syntax=docker/dockerfile:1.7

FROM golang:1.24 AS builder
WORKDIR /app

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOTOOLCHAIN=local

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags "-s -w" -o /build/server ./cmd/server

FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app

COPY --from=builder /build/server ./server
COPY configs ./configs

ENV WONDER_SERVER_HOST=0.0.0.0 \
    WONDER_LOG_OUTPUT=stdout \
    WONDER_LOG_ENABLE_FILE=false

EXPOSE 8080

ENTRYPOINT ["./server"]
