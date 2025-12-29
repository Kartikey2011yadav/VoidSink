# Monitoring & Observability

VoidSink includes a complete observability stack to help you visualize attacks in real-time.

## The Stack

The monitoring stack consists of three components running in Docker containers:

1.  **VoidSink**: Exposes metrics at `http://localhost:9090/metrics`.
2.  **Prometheus**: Scrapes the metrics from VoidSink every 5 seconds.
3.  **Grafana**: Visualizes the data from Prometheus.

## Metrics

VoidSink exports standard Go runtime metrics and custom application metrics:

| Metric Name | Type | Description |
| :--- | :--- | :--- |
| `voidsink_active_traps` | Gauge | The number of currently open connections (attackers stuck in the tarpit). |
| `voidsink_traps_total` | Counter | The total number of connections accepted since startup. |
| `voidsink_bytes_sent_total` | Counter | The total amount of garbage data sent to attackers (in bytes). |

## Grafana Dashboard

A pre-configured dashboard is included in `deploy/grafana/dashboards/voidsink.json`.

### Accessing the Dashboard

1.  Ensure the stack is running: `docker-compose up -d`
2.  Open your browser to [http://localhost:3000](http://localhost:3000).
3.  Login with default credentials:
    - User: `admin`
    - Password: `admin`
4.  Navigate to **Dashboards** > **VoidSink**.

### Key Panels

- **Active Traps**: A real-time gauge of how many bots are currently connected. A sudden spike here indicates a new attack wave.
- **Data Egress Rate**: Shows the bandwidth being consumed by the tarpit.
- **Total Traps**: A cumulative count of all victims.

## Troubleshooting

- **"No Data" in Grafana**:
    - Check if Prometheus is running: `http://localhost:9091/targets`.
    - Ensure VoidSink is reachable by Prometheus (they share the `voidsink_default` network).
    - Verify metrics are being generated: `curl http://localhost:9090/metrics`.
