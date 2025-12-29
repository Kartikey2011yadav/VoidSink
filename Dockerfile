# Build Stage
FROM golang:1.25.5 AS builder

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
COPY configs/config.yaml /app/configs/config.yaml
COPY assets /app/assets

EXPOSE 8080
EXPOSE 9090

ENTRYPOINT ["/app/VoidSink", "-c", "/app/configs/config.yaml"]
