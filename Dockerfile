FROM golang:1.21 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o aion-api ./cmd/aion-api

FROM alpine:3.19.1
WORKDIR /app
COPY --from=builder /app/aion-api /app/
EXPOSE 5001
CMD ["./aion-api"]