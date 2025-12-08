# syntax=docker/dockerfile:1.4
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Allow controlling ldflags via build-arg. For dev we omit strips, for prod we strip symbols.
ARG BUILD_LDFLAGS="-s -w"

# Use build cache for go build and module cache to speed up repeated builds.
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="${BUILD_LDFLAGS}" -o aion-api ./cmd/aion-api

FROM alpine:3.19.1

WORKDIR /app

COPY --from=builder /app/aion-api .

COPY infrastructure/docker/scripts/entrypoint.sh .

RUN chmod +x ./entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]