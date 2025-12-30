package jsontrap

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Kartikey2011yadav/voidsink/internal/heffalump"
	"github.com/Kartikey2011yadav/voidsink/internal/telemetry"
	"github.com/Kartikey2011yadav/voidsink/pkg/notifier"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

// JSONInfiniteTrap implements the Trap interface for a JSON honeypot.
type JSONInfiniteTrap struct {
	addr       string
	serverName string
	server     *fasthttp.Server
	heffalump  *heffalump.Heffalump
	pool       *heffalump.BufferPool
	notifier   *notifier.Notifier
}

// New creates a new instance of JSONInfiniteTrap.
func New(addr, serverName string, h *heffalump.Heffalump, n *notifier.Notifier) *JSONInfiniteTrap {
	return &JSONInfiniteTrap{
		addr:       addr,
		serverName: serverName,
		heffalump:  h,
		pool:       heffalump.NewBufferPool(),
		notifier:   n,
	}
}

// Start starts the HTTP server.
func (t *JSONInfiniteTrap) Start(ctx context.Context) error {
	t.server = &fasthttp.Server{
		Handler:      t.requestHandler,
		Name:         t.serverName,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 0, // Infinite
		IdleTimeout:  30 * time.Second,
	}

	log.Info().Str("address", t.addr).Msg("Starting JSON Infinite Trap")

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
func (t *JSONInfiniteTrap) Shutdown(ctx context.Context) error {
	log.Info().Msg("Shutting down JSON Infinite Trap")
	if t.server != nil {
		return t.server.Shutdown()
	}
	return nil
}

func (t *JSONInfiniteTrap) requestHandler(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())
	remoteIP := ctx.RemoteAddr().String()
	userAgent := string(ctx.UserAgent())

	log.Info().Str("path", path).Str("remote_addr", remoteIP).Msg("JSON Trap hit")

	telemetry.TrapsTriggered.WithLabelValues("json_infinite", path).Inc()

	// Send Alert
	if t.notifier != nil {
		t.notifier.SendAlert("JSONInfinite", remoteIP, userAgent)
	}

	ctx.SetContentType("application/json")
	ctx.SetBodyStreamWriter(func(w *bufio.Writer) {
		telemetry.ActiveConnections.Inc()
		defer telemetry.ActiveConnections.Dec()

		// Start JSON Array
		w.WriteString("[\n")
		if err := w.Flush(); err != nil {
			return
		}

		// Get a buffer from the pool for constructing the JSON object
		buf := t.pool.Get()
		defer t.pool.Put(buf)

		var id uint64 = 1

		// Seed the generator for this connection
		w1, w2 := t.heffalump.Seed()

		for {
			// Reset buffer for the new object
			buf.Reset()

			// Generate Heffalump text
			// We use a separate buffer for the raw text to avoid mixing it with the JSON structure
			textBuf := t.pool.Get()

			// Generate ~500 bytes of text
			for textBuf.Len() < 500 {
				w3 := t.heffalump.Next(w1, w2)
				textBuf.WriteString(w3)
				textBuf.WriteByte(' ')
				w1, w2 = w2, w3
			}

			// Marshal just the string content to get proper JSON escaping (quotes included)
			escapedJSONString, _ := json.Marshal(textBuf.String())

			// We are done with textBuf
			t.pool.Put(textBuf)

			// Construct the final JSON object
			// {"id": <id>, "timestamp": "<time>", "data": <escaped_text>}

			buf.WriteString(`{"id":`)
			fmt.Fprintf(buf, "%d", id)
			buf.WriteString(`,"timestamp":"`)
			buf.WriteString(time.Now().Format(time.RFC3339))
			buf.WriteString(`","data":`)
			buf.Write(escapedJSONString)
			buf.WriteString(`},` + "\n")

			// Write to the response stream
			if _, err := w.Write(buf.Bytes()); err != nil {
				return // Client disconnected
			}

			if err := w.Flush(); err != nil {
				return
			}

			id++
		}
	})
}
