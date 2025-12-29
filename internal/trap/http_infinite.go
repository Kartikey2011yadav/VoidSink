package trap

import (
	"bufio"
	"context"
	"time"

	"github.com/Kartikey2011yadav/voidsink/internal/heffalump"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

// HTTPInfiniteTrap implements the Trap interface for an HTTP honeypot
// that streams infinite data to clients.
type HTTPInfiniteTrap struct {
	addr      string
	server    *fasthttp.Server
	heffalump *heffalump.Heffalump
	pool      *heffalump.BufferPool
}

// NewHTTPInfiniteTrap creates a new instance of HTTPInfiniteTrap.
func NewHTTPInfiniteTrap(addr string, h *heffalump.Heffalump) *HTTPInfiniteTrap {
	return &HTTPInfiniteTrap{
		addr:      addr,
		heffalump: h,
		pool:      heffalump.NewBufferPool(),
	}
}

// Start starts the HTTP server.
func (t *HTTPInfiniteTrap) Start(ctx context.Context) error {
	t.server = &fasthttp.Server{
		Handler: t.requestHandler,
		Name:    "VoidSink/1.0",
		// Increase timeouts to allow long-running connections (it's a trap after all)
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 0, // Infinite
		IdleTimeout:  30 * time.Second,
	}

	log.Info().Str("address", t.addr).Msg("Starting HTTP Infinite Trap")

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
func (t *HTTPInfiniteTrap) Shutdown(ctx context.Context) error {
	log.Info().Msg("Shutting down HTTP Infinite Trap")
	if t.server != nil {
		return t.server.Shutdown()
	}
	return nil
}

func (t *HTTPInfiniteTrap) requestHandler(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())
	log.Info().Str("path", path).Str("remote_addr", ctx.RemoteAddr().String()).Msg("Trap hit")

	switch path {
	case "/robots.txt":
		ctx.SetContentType("text/plain")
		ctx.WriteString("User-agent: *\nDisallow: /")
	default:
		// HellPot logic: Stream infinite Markov chain data
		ctx.SetContentType("text/html")
		ctx.SetBodyStreamWriter(func(w *bufio.Writer) {
			// Get a buffer from the pool
			buf := t.pool.Get()
			defer t.pool.Put(buf)

			// Seed the generator for this connection
			w1, w2 := t.heffalump.Seed()

			for {
				// Fill the buffer with ~4KB of data
				// We check buf.Len() < 4000 to leave a little room for the last word
				for buf.Len() < 4000 {
					w3 := t.heffalump.Next(w1, w2)
					buf.WriteString(w3)
					buf.WriteByte(' ') // Add space
					w1, w2 = w2, w3
				}

				// Write the buffer to the network
				if _, err := w.Write(buf.Bytes()); err != nil {
					// Connection closed by client or error
					log.Debug().Err(err).Msg("Connection closed during streaming")
					return
				}

				// Flush to ensure data is sent
				if err := w.Flush(); err != nil {
					return
				}

				// Reset buffer for next iteration
				buf.Reset()

				// Optional: Small sleep to prevent 100% CPU usage if the network is too fast
				// In a real honeypot, you might want to limit bitrate to keep them hooked longer
				// without burning your own bandwidth.
				// time.Sleep(50 * time.Millisecond)
			}
		})
	}
}
