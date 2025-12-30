package gziptrap

import (
	"bufio"
	"compress/gzip"
	"context"
	"time"

	"github.com/Kartikey2011yadav/voidsink/internal/telemetry"
	"github.com/Kartikey2011yadav/voidsink/pkg/notifier"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

// GzipTrap implements the Trap interface for a Gzip honeypot.
type GzipTrap struct {
	addr       string
	serverName string
	server     *fasthttp.Server
	notifier   *notifier.Notifier
	zeroBuf    []byte
}

// New creates a new instance of GzipTrap.
func New(addr, serverName string, n *notifier.Notifier) *GzipTrap {
	// Pre-allocate a 32KB buffer of zeros to minimize allocations during the loop
	return &GzipTrap{
		addr:       addr,
		serverName: serverName,
		notifier:   n,
		zeroBuf:    make([]byte, 32*1024),
	}
}

// Start starts the HTTP server.
func (t *GzipTrap) Start(ctx context.Context) error {
	t.server = &fasthttp.Server{
		Handler:      t.requestHandler,
		Name:         t.serverName,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 0, // Infinite
		IdleTimeout:  30 * time.Second,
	}

	log.Info().Str("address", t.addr).Msg("Starting Gzip Infinite Trap")

	errChan := make(chan error, 1)
	go func() {
		errChan <- t.server.ListenAndServe(t.addr)
	}()

	select {
	case <-ctx.Done():
		return t.Shutdown(context.Background())
	case err := <-errChan:
		return err
	}
}

// Shutdown gracefully shuts down the server.
func (t *GzipTrap) Shutdown(ctx context.Context) error {
	log.Info().Msg("Shutting down Gzip Infinite Trap")
	if t.server != nil {
		return t.server.Shutdown()
	}
	return nil
}

func (t *GzipTrap) requestHandler(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())
	remoteIP := ctx.RemoteAddr().String()
	userAgent := string(ctx.UserAgent())

	log.Info().Str("path", path).Str("remote_addr", remoteIP).Msg("Gzip Trap hit")

	telemetry.TrapsTriggered.WithLabelValues("gzip_infinite", path).Inc()

	// Send Alert
	if t.notifier != nil {
		t.notifier.SendAlert("GzipInfinite", remoteIP, userAgent)
	}

	ctx.SetContentType("text/plain")
	ctx.Response.Header.Set("Content-Encoding", "gzip")

	ctx.SetBodyStreamWriter(func(w *bufio.Writer) {
		telemetry.ActiveConnections.Inc()
		defer telemetry.ActiveConnections.Dec()

		// Use BestCompression to maximize the expansion ratio (Gzip Bomb effect)
		// This makes the client work hard to decompress while we send very little data.
		gw, err := gzip.NewWriterLevel(w, gzip.BestCompression)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create gzip writer")
			return
		}
		defer gw.Close()

		for {
			// Write the pre-allocated zeros
			// We write uncompressed zeros to the gzip writer.
			// The gzip writer compresses them and writes to 'w'.
			n, err := gw.Write(t.zeroBuf)
			if err != nil {
				return // Client disconnected
			}

			// We track the uncompressed bytes "sent" (generated)
			// This represents the size of the data the attacker has to process.
			telemetry.BytesSent.Add(float64(n))

			// Flush the gzip writer to ensure data is pushed to the underlying writer
			if err := gw.Flush(); err != nil {
				return
			}

			// Flush the underlying writer to send data over the network
			if err := w.Flush(); err != nil {
				return
			}
		}
	})
}
