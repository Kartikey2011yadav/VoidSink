# VoidSink Architecture

VoidSink is designed as a high-performance, modular framework for deploying network tarpits and honeypots. Unlike traditional monolithic honeypots, VoidSink uses a decoupled architecture that separates the core server logic from the specific "trap" implementations and the content generation engine.

## High-Level Overview

The application is structured around three main pillars:
1.  **The Trap Handler**: Manages incoming network connections and protocol specifics.
2.  **The Heffalump Engine**: Generates the content (payload) sent to the attacker.
3.  **Telemetry & Observability**: Tracks the performance and effectiveness of the traps.

## Directory Structure

The project follows the Standard Go Project Layout:

- `cmd/`: Main entry points for the application.
- `configs/`: Default configuration files.
- `deploy/`: Deployment assets (Dockerfiles, Compose files, Grafana dashboards).
- `internal/`: Private application and library code.
    - `config/`: Configuration loading and validation using `koanf`.
    - `heffalump/`: The Markov Chain text generation engine.
    - `logger/`: Structured logging setup using `zerolog`.
    - `telemetry/`: Prometheus metrics definitions and exporters.
    - `trap/`: The core interface and implementations for different trap types (e.g., HTTP).

## Core Components

### 1. The Trap Interface

At the heart of VoidSink is the `Trap` interface. This allows the application to support multiple protocols (HTTP, TCP, SMTP, etc.) without changing the main application loop.

Currently, the primary implementation is the **HTTP Trap**, built on top of `valyala/fasthttp`.

**Why `fasthttp`?**
Standard `net/http` in Go creates a new goroutine for every request. While lightweight, this can still lead to resource exhaustion under heavy load (which is the goal of a tarpit). `fasthttp` uses a worker pool model and zero-allocation practices, allowing VoidSink to handle tens of thousands of concurrent "stuck" connections with minimal memory footprint.

### 2. The Heffalump Engine

The Heffalump Engine is responsible for generating the "bait". It is a streaming Markov Chain generator.

- **Input**: It ingests a corpus of text (e.g., classic literature).
- **Processing**: It builds a frequency map of word transitions.
- **Output**: It produces an infinite stream of coherent-looking but nonsensical text.

This engine is crucial because static files are easily fingerprinted by attackers. The Heffalump engine ensures that every byte stream is unique, making it difficult for automated scanners to detect they are in a tarpit based on content signatures alone.

### 3. Telemetry

VoidSink is "Observability First". It doesn't just log errors; it exposes internal state via Prometheus metrics.

- **`voidsink_active_traps`**: Gauge showing how many attackers are currently stuck.
- **`voidsink_bytes_sent_total`**: Counter showing how much garbage data has been sent.
- **`voidsink_traps_total`**: Counter of total connections handled.

These metrics are exposed on a dedicated port (default `:9090`) to separate monitoring traffic from the trap traffic.
