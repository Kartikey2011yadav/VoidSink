# Build Stage
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build static binary
RUN CGO_ENABLED=0 go build -trimpath -o VoidSink ./cmd/voidsink

# Runtime Stage
FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=builder /app/VoidSink /app/VoidSink
COPY docker_config.toml /config/config.toml

EXPOSE 8080

ENTRYPOINT ["/app/VoidSink", "-c", "/config/config.toml"]
