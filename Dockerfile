FROM golang:1.25.0 AS builder
WORKDIR /app
COPY . .
WORKDIR /app/cmd
RUN go build -o urbioxe main.go

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/cmd/urbioxe /urbioxe

ENTRYPOINT ["/urbioxe"]
