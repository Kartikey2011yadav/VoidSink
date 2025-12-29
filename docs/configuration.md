# Configuration

VoidSink uses a flexible configuration system powered by `koanf`. It supports configuration via a YAML file, environment variables, and command-line flags.

## Configuration File (`config.yaml`)

The default configuration is located at `configs/config.yaml`.

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: 10s
  write_timeout: 0s # 0 means infinite (required for tarpit)
  max_conns_per_ip: 10

metrics:
  enabled: true
  port: 9090
  path: "/metrics"

log:
  level: "info" # debug, info, warn, error
  format: "json" # json, console
```

## Environment Variables

Any configuration option can be overridden using environment variables. The prefix is `VOIDSINK_`. Use double underscores `__` to separate nested keys.

Examples:
- `VOIDSINK_SERVER__PORT=8081` overrides `server.port`.
- `VOIDSINK_LOG__LEVEL=debug` overrides `log.level`.

## Docker Configuration

When running in Docker, the configuration file is mounted at `/app/configs/config.yaml`.

To modify the configuration without rebuilding the image, you can:
1.  Edit `configs/config.yaml` on your host machine.
2.  Restart the container: `docker-compose restart voidsink`.

Alternatively, you can set environment variables in the `docker-compose.yml` file:

```yaml
services:
  voidsink:
    environment:
      - VOIDSINK_SERVER__MAX_CONNS_PER_IP=50
```
