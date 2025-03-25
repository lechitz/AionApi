FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o aion-api ./cmd/aion-api

FROM alpine:3.19.1

WORKDIR /app

COPY --from=builder /app/aion-api .

ENTRYPOINT ["./aion-api"]