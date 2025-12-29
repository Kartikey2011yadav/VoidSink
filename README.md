# VoidSink

<div align="center">

![VoidSink Logo](https://via.placeholder.com/150?text=VoidSink) <!-- Placeholder for a logo if you have one -->

**The High-Performance, Modular HTTP Tarpit & Honeypot Framework**

[![Go Version](https://img.shields.io/github/go-mod/go-version/Kartikey2011yadav/voidsink?style=flat-square&logo=go)](https://go.dev/)
[![Build Status](https://img.shields.io/github/actions/workflow/status/Kartikey2011yadav/voidsink/release.yml?branch=main&style=flat-square&logo=github)](https://github.com/Kartikey2011yadav/voidsink/actions)
[![Docker Image](https://img.shields.io/badge/docker-ready-blue?style=flat-square&logo=docker)](https://hub.docker.com/r/kartikey2011yadav/voidsink)
[![License](https://img.shields.io/badge/license-MIT-green?style=flat-square)](LICENSE)

</div>

---

## What is VoidSink?

**VoidSink** is a next-generation honeypot designed to trap malicious bots, scanners, and scrapers in an infinite loop of generated garbage data. Inspired by the legendary [HellPot](https://github.com/yunginnanet/HellPot), VoidSink takes the concept further with a **modular, interface-driven architecture** that prioritizes extensibility, observability, and stealth.

When a bot connects to VoidSink, it is fed an endless stream of Markov-chain-generated text (based on classic literature or any corpus you provide). This consumes the attacker's bandwidth and resources while keeping your actual infrastructure safe.

## Key Features

- **Infinite Tarpit**: Streams gigabytes of unique, non-repeating garbage data to keep connections open indefinitely.
- **Heffalump Engine**: Advanced Markov Chain generator that produces realistic-looking HTML/Text to fool sophisticated scrapers.
- **Stealth Mode**: Configurable `Server` header masquerading (e.g., mimics Nginx, Apache) to avoid fingerprinting.
- **Built-in Observability**: Native Prometheus metrics exporter (`/metrics`) to track active traps, bytes sent, and attack vectors.
- **High Performance**: Built on `valyala/fasthttp` and `sync.Pool` for zero-allocation hot paths and massive concurrency.
- **Modular Design**: Interface-based architecture (`Trap` interface) allows easy addition of new protocols (SSH, SMTP, TCP) without refactoring the core.
- **Cloud Native**: Multi-stage Docker builds, structured JSON logging (`zerolog`), and 12-factor app configuration.

## Architecture: VoidSink vs. HellPot

VoidSink is an architectural evolution of the HellPot concept:

| Feature | HellPot | VoidSink |
| :--- | :--- | :--- |
| **Design** | Monolithic HTTP Server | **Modular Framework** (Supports multiple Trap types) |
| **Engine** | Stream-oriented `io.Reader` | **Logic-oriented State Machine** (Decoupled generation) |
| **Config** | TOML | **Koanf** (YAML, JSON, Env Vars, Flags) |
| **Metrics** | External Exporters | **Native Prometheus Integration** |
| **Stealth** | Basic | **Header Masquerading & HTML Tokenization** |

## Installation & Usage

### Option 1: Docker (Recommended)

```bash
# Pull the latest image
docker pull kartikey2011yadav/voidsink:latest

# Run with default settings
docker run -d -p 8080:8080 -p 9090:9090 --name voidsink kartikey2011yadav/voidsink:latest
```

### Option 2: Docker Compose

```yaml
version: '3.8'
services:
  voidsink:
    image: kartikey2011yadav/voidsink:latest
    ports:
      - "8080:8080" # Trap Port
      - "9090:9090" # Metrics Port
    volumes:
      - ./configs/config.yaml:/app/configs/config.yaml
    restart: always
```

### Option 3: Build from Source

**Prerequisites**: Go 1.25+

```bash
# Clone the repository
git clone https://github.com/Kartikey2011yadav/voidsink.git
cd voidsink

# Build the binary
go build -o voidsink ./cmd/voidsink

# Run
./voidsink -c configs/config.yaml
```

## Configuration

VoidSink uses a flexible configuration system. You can configure it via `configs/config.yaml`:

```yaml
# Logging Configuration
log_level: debug # debug, info, warn, error
log_format: json # json, console

# Prometheus Metrics
metrics:
  enabled: true
  addr: ":9090"

# Trap Configuration
traps:
  http_infinite:
    enabled: true
    addr: ":8080"
    # Masquerade as a real server to fool scanners
    server_name: "nginx" 
```

## Monitoring

VoidSink exposes Prometheus metrics at `http://localhost:9090/metrics`.

**Key Metrics:**

- `voidsink_active_connections`: Number of bots currently trapped.
- `voidsink_bytes_sent_total`: Total garbage data sent to attackers.
- `voidsink_traps_triggered_total`: Count of hits per endpoint (e.g., `/wp-login.php`).

## Contributing

Contributions are welcome! Whether it's adding a new Trap type (SSH, Telnet), improving the Heffalump engine, or fixing bugs.

1. Fork the repository.
2. Create your feature branch (`git checkout -b feature/amazing-feature`).
3. Commit your changes (`git commit -m 'Add some amazing feature'`).
4. Push to the branch (`git push origin feature/amazing-feature`).
5. Open a Pull Request.

## License

Distributed under the MIT License. See `LICENSE` for more information.
