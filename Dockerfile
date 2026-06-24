# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS builder

WORKDIR /src

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /croquing ./cmd/server

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /croquing .

EXPOSE 8080

USER nobody:nobody

ENTRYPOINT ["./croquing"]
