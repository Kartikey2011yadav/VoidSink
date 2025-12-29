package trap

import (
	"bufio"
	"context"
	"time"
	"github.com/Kartikey2011yadav/voidsink/pkg/markov"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

// HTTPInfiniteTrap implements the Trap interface for an HTTP honeypot
// that streams infinite data to clients.
type HTTPInfiniteTrap struct {
	addr      string
	server    *fasthttp.Server
	generator *markov.Generator
}

// NewHTTPInfiniteTrap creates a new instance of HTTPInfiniteTrap.
func NewHTTPInfiniteTrap(addr string) *HTTPInfiniteTrap {
	return &HTTPInfiniteTrap{
		addr:      addr,
		generator: markov.NewGenerator(),
	}
}

// Start starts the HTTP server.
func (t *HTTPInfiniteTrap) Start(ctx context.Context) error {
	t.server = &fasthttp.Server{
		Handler: t.requestHandler,
		Name:    "VoidSink/1.0",
	}

	log.Info().Str("address", t.addr).Msg("Starting HTTP Infinite Trap")

	// ListenAndServe blocks. We can use a goroutine to handle context cancellation if needed,
	// but usually Shutdown is called from another goroutine (e.g. signal handler).
	// However, to respect the context passed in Start, we should monitor it.

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
		// HellPot logic: Stream infinite data
		ctx.SetContentType("text/html")
		ctx.SetBodyStreamWriter(func(w *bufio.Writer) {
			for {
				// Generate a chunk of data
				chunk := t.generator.Generate(1024) // 1KB chunk
				if _, err := w.Write(chunk); err != nil {
					// Connection closed by client or error
					log.Debug().Err(err).Msg("Connection closed during streaming")
					return
				}
				if _, err := w.Write([]byte("\n")); err != nil {
					return
				}
				if err := w.Flush(); err != nil {
					return
				}
				// Small sleep to control rate if necessary, or just blast it.
				// HellPot usually blasts it, but let's be nice to CPU for now.
				time.Sleep(10 * time.Millisecond)
			}
		})
	}
}
