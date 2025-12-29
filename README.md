# VoidSink

VoidSink is a high-performance, modular honeypot framework written in Go.

## Features

- **Modular Architecture**: Easily extensible Trap interface.
- **High Performance**: Uses `fasthttp` and `zerolog`.
- **HellPot Logic**: Includes an HTTP Infinite Trap that streams infinite data to bots.

## Getting Started

### Prerequisites

- Go 1.21+

### Running

```bash
go mod tidy
go run cmd/voidsink/main.go
```

## Configuration

Configuration is located in `configs/config.yaml`.

```yaml
log_level: debug

traps:
  http_infinite:
    enabled: true
    addr: ":8080"
```
