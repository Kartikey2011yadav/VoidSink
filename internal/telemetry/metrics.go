package telemetry

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// BytesSent tracks the total number of bytes sent to trapped clients.
	BytesSent = promauto.NewCounter(prometheus.CounterOpts{
		Name: "voidsink_bytes_sent_total",
		Help: "The total number of bytes sent to trapped clients",
	})

	// ActiveConnections tracks the number of currently active trapped connections.
	ActiveConnections = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "voidsink_active_connections",
		Help: "The number of currently active trapped connections",
	})

	// TrapsTriggered tracks the total number of times a trap has been hit.
	TrapsTriggered = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "voidsink_traps_triggered_total",
		Help: "The total number of times a trap has been hit",
	}, []string{"trap_type", "path"})
)
