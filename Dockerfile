# build stage
FROM golang:1.20 AS builder
# working directory
WORKDIR /app
COPY ./ /app

RUN go get -d /app/cmd/aion-api

# rebuilt built in libraries and disabled cgo
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app /app/cmd/aion-api
# final stage
FROM alpine:3.19.1

# working directory
WORKDIR /app

# copy the binary file into working directory
COPY --from=builder /app .
# http server listens on port 5001
EXPOSE 5001
# Run the docker_imgs command when the container starts.
CMD ["/app/aion-api"]