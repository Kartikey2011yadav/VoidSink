package spidertrap

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Kartikey2011yadav/voidsink/internal/heffalump"
	"github.com/Kartikey2011yadav/voidsink/internal/telemetry"
	"github.com/Kartikey2011yadav/voidsink/pkg/notifier"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

// SpiderTrap implements a recursive directory trap to confuse web crawlers.
type SpiderTrap struct {
	addr       string
	serverName string
	server     *fasthttp.Server
	heffalump  *heffalump.Heffalump
	notifier   *notifier.Notifier
}

// New creates a new instance of SpiderTrap.
func New(addr, serverName string, h *heffalump.Heffalump, n *notifier.Notifier) *SpiderTrap {
	return &SpiderTrap{
		addr:       addr,
		serverName: serverName,
		heffalump:  h,
		notifier:   n,
	}
}

// Start starts the HTTP server.
func (t *SpiderTrap) Start(ctx context.Context) error {
	t.server = &fasthttp.Server{
		Handler:      t.requestHandler,
		Name:         t.serverName,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second, // We want to send the page quickly
		IdleTimeout:  30 * time.Second,
	}

	log.Info().Str("address", t.addr).Msg("Starting Spider Trap")

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
func (t *SpiderTrap) Shutdown(ctx context.Context) error {
	log.Info().Msg("Shutting down Spider Trap")
	if t.server != nil {
		return t.server.Shutdown()
	}
	return nil
}

func (t *SpiderTrap) requestHandler(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())
	remoteIP := ctx.RemoteAddr().String()
	userAgent := string(ctx.UserAgent())

	log.Info().Str("path", path).Str("remote_addr", remoteIP).Msg("Spider Trap hit")
	telemetry.TrapsTriggered.WithLabelValues("spider_trap", path).Inc()

	// Send Alert
	if t.notifier != nil {
		t.notifier.SendAlert("SpiderTrap", remoteIP, userAgent)
	}

	ctx.SetContentType("text/html")

	// Generate a simple HTML page with recursive links
	fmt.Fprintf(ctx, "<!DOCTYPE html><html><head><title>Index of %s</title></head><body>", path)
	fmt.Fprintf(ctx, "<h1>Index of %s</h1><hr><pre>", path)

	// Link to parent
	fmt.Fprintf(ctx, "<a href=\"../\">../</a>\n")

	// Generate 5-10 random subdirectories
	count := rand.Intn(6) + 5 // 5 to 10

	// Seed heffalump for this request
	w1, w2 := t.heffalump.Seed()

	for i := 0; i < count; i++ {
		// Get a random word for the directory name
		word := t.heffalump.Next(w1, w2)

		// Clean up the word to make it URL-safe(ish) and look like a directory
		cleanWord := strings.Map(func(r rune) rune {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
				return r
			}
			return -1
		}, word)

		if len(cleanWord) < 3 {
			cleanWord = "folder" + cleanWord // Ensure it's not empty or too short
		}

		// Construct the link
		// We use relative paths to keep the crawler going deeper
		// If current path ends with /, append directly. If not, add /
		link := cleanWord + "/"

		fmt.Fprintf(ctx, "<a href=\"%s\">%s</a>\n", link, link)

		// Advance the chain
		w3 := t.heffalump.Next(w1, w2)
		w1, w2 = w2, w3
	}

	fmt.Fprintf(ctx, "</pre><hr></body></html>")
}
