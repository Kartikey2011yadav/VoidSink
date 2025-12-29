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

## Documentation

Detailed documentation is available in the `docs/` directory:

- [Architecture Overview](docs/architecture.md): Learn about the modular design and core components.
- [The Heffalump Engine](docs/heffalump_engine.md): Understand how the infinite text generation works.
- [Monitoring & Metrics](docs/monitoring.md): Guide to the Prometheus and Grafana setup.
- [Configuration](docs/configuration.md): Reference for config files and environment variables.

## Installation & Usage

### Option 1: Full Stack (Recommended)

The easiest way to run VoidSink with full monitoring (Grafana + Prometheus) is using Docker Compose.

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/Kartikey2011yadav/voidsink.git
    cd voidsink
    ```

2.  **Start the stack**:
    ```bash
    docker-compose up -d
    ```

3.  **Access the services**:
    - **VoidSink (Tarpit)**: `http://localhost:8080` (This is the trap!)
    - **Grafana**: `http://localhost:3000` (Login: `admin`/`admin`)
    - **Prometheus**: `http://localhost:9091`

### Option 2: Standalone Docker

If you only want the tarpit without the monitoring stack:

```bash
docker build -t voidsink .
docker run -p 8080:8080 -p 9090:9090 voidsink
```

### Option 3: Build from Source

Requires Go 1.21+.

```bash
go build -o voidsink cmd/voidsink/main.go
./voidsink
```
## Contributing

Contributions are welcome! Whether it's adding a new Trap type (SSH, Telnet), improving the Heffalump engine, or fixing bugs.

1. Fork the repository.
2. Create your feature branch (`git checkout -b feature/amazing-feature`).
3. Commit your changes (`git commit -m 'Add some amazing feature'`).
4. Push to the branch (`git push origin feature/amazing-feature`).
5. Open a Pull Request.

## License

Distributed under the MIT License. See `LICENSE` for more information.
