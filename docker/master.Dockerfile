FROM golang:1.24 AS builder
WORKDIR /app
COPY MapReduce .
RUN go build -o bin/inverted-index ./cmd/invertedindex

FROM ubuntu:22.04

# Install certificates and any other needed dependencies in one RUN layer
RUN apt-get update && \
    apt-get install -y curl ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/bin/inverted-index /inverted-index

EXPOSE 9000
EXPOSE 9100

ENTRYPOINT ["/inverted-index"]