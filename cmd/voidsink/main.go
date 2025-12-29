package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Kartikey2011yadav/voidsink/internal/config"
	"github.com/Kartikey2011yadav/voidsink/internal/logger"
	"github.com/Kartikey2011yadav/voidsink/internal/trap"

	"github.com/rs/zerolog/log"
)

func main() {
	configPath := flag.String("c", "configs/config.yaml", "Path to configuration file")
	flag.Parse()

	// 1. Load Configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		// Fallback or fatal? Let's just log and exit if we can't load config.
		// But config.Load handles missing file gracefully (returns empty config).
		// If it returns error it's a parsing error.
		panic(err)
	}

	// 2. Setup Logger
	logger.Setup(cfg.LogLevel, cfg.LogFile, cfg.LogFormat)
	log.Info().Msg("VoidSink starting up...")

	// 3. Initialize Traps
	var traps []trap.Trap

	if cfg.Traps.HTTPInfinite.Enabled {
		t := trap.NewHTTPInfiniteTrap(cfg.Traps.HTTPInfinite.Addr)
		traps = append(traps, t)
		log.Info().Str("type", "HTTPInfinite").Str("addr", cfg.Traps.HTTPInfinite.Addr).Msg("Trap enabled")
	}

	if len(traps) == 0 {
		log.Warn().Msg("No traps enabled. Exiting.")
		return
	}

	// 4. Start Traps
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	for _, t := range traps {
		wg.Add(1)
		go func(tr trap.Trap) {
			defer wg.Done()
			if err := tr.Start(ctx); err != nil {
				log.Error().Err(err).Msg("Trap stopped with error")
			}
		}(t)
	}

	// 5. Wait for Interrupt Signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down VoidSink...")

	// 6. Graceful Shutdown
	cancel() // Signal traps to stop

	// Give traps some time to cleanup
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	// We can also call Shutdown() explicitly on traps if they don't respect context cancellation in Start()
	// But our interface design implies Start(ctx) should handle it.
	// However, for completeness, let's iterate and call Shutdown just in case Start is blocking on something else.
	// But wait, if Start blocks, we can't call Shutdown on the same object easily unless we stored them.
	// We have them in `traps` slice.

	// Let's do explicit Shutdown calls in parallel or sequence
	var shutdownWg sync.WaitGroup
	for _, t := range traps {
		shutdownWg.Add(1)
		go func(tr trap.Trap) {
			defer shutdownWg.Done()
			if err := tr.Shutdown(shutdownCtx); err != nil {
				log.Error().Err(err).Msg("Error during trap shutdown")
			}
		}(t)
	}
	shutdownWg.Wait()

	wg.Wait() // Wait for Start() goroutines to return
	log.Info().Msg("VoidSink shutdown complete")
}
