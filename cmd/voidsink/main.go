package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Kartikey2011yadav/voidsink/internal/config"
	"github.com/Kartikey2011yadav/voidsink/internal/heffalump"
	"github.com/Kartikey2011yadav/voidsink/internal/logger"
	"github.com/Kartikey2011yadav/voidsink/internal/trap"
	gziptrap "github.com/Kartikey2011yadav/voidsink/internal/traps/gzip"
	httptrap "github.com/Kartikey2011yadav/voidsink/internal/traps/http"
	jsontrap "github.com/Kartikey2011yadav/voidsink/internal/traps/json"
	logintrap "github.com/Kartikey2011yadav/voidsink/internal/traps/login"
	spidertrap "github.com/Kartikey2011yadav/voidsink/internal/traps/spider"
	"github.com/Kartikey2011yadav/voidsink/pkg/notifier"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

func main() {
	configPath := flag.String("c", "configs/config.yaml", "Path to configuration file")
	flag.Parse()

	// 1. Load Configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		panic(err)
	}

	// 2. Setup Logger
	logger.Setup(cfg.LogLevel, cfg.LogFile, cfg.LogFormat)
	log.Info().Msg("VoidSink starting up...")

	// 3. Initialize Core Components
	heffalumpEngine, err := heffalump.New("assets/corpus.txt")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize Heffalump engine")
	}

	alertNotifier := notifier.New(cfg.Notification.WebhookURL)

	// 4. Start Metrics Server
	if cfg.Metrics.Enabled {
		go startMetricsServer(cfg.Metrics.Addr)
	}

	// 5. Initialize Traps
	traps := initializeTraps(cfg, heffalumpEngine, alertNotifier)
	if len(traps) == 0 {
		log.Warn().Msg("No traps enabled. Exiting.")
		return
	}

	// 6. Start Traps
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

	// 7. Wait for Interrupt Signal
	waitForShutdown(cancel, traps, &wg)
}

func startMetricsServer(addr string) {
	http.Handle("/metrics", promhttp.Handler())
	log.Info().Str("addr", addr).Msg("Starting Metrics Server")
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Error().Err(err).Msg("Metrics server failed")
	}
}

func initializeTraps(cfg *config.Config, h *heffalump.Heffalump, n *notifier.Notifier) []trap.Trap {
	var traps []trap.Trap

	if cfg.Traps.HTTPInfinite.Enabled {
		t := httptrap.New(cfg.Traps.HTTPInfinite.Addr, cfg.Traps.HTTPInfinite.ServerName, h, n)
		traps = append(traps, t)
		log.Info().Str("type", "HTTPInfinite").Str("addr", cfg.Traps.HTTPInfinite.Addr).Msg("Trap enabled")
	}

	if cfg.Traps.JSONInfinite.Enabled {
		t := jsontrap.New(cfg.Traps.JSONInfinite.Addr, cfg.Traps.JSONInfinite.ServerName, h, n)
		traps = append(traps, t)
		log.Info().Str("type", "JSONInfinite").Str("addr", cfg.Traps.JSONInfinite.Addr).Msg("Trap enabled")
	}

	if cfg.Traps.SpiderTrap.Enabled {
		t := spidertrap.New(cfg.Traps.SpiderTrap.Addr, cfg.Traps.SpiderTrap.ServerName, h, n)
		traps = append(traps, t)
		log.Info().Str("type", "SpiderTrap").Str("addr", cfg.Traps.SpiderTrap.Addr).Msg("Trap enabled")
	}

	if cfg.Traps.GzipInfinite.Enabled {
		t := gziptrap.New(cfg.Traps.GzipInfinite.Addr, cfg.Traps.GzipInfinite.ServerName, n)
		traps = append(traps, t)
		log.Info().Str("type", "GzipInfinite").Str("addr", cfg.Traps.GzipInfinite.Addr).Msg("Trap enabled")
	}

	if cfg.Traps.LoginTrap.Enabled {
		t := logintrap.New(cfg.Traps.LoginTrap.Addr, cfg.Traps.LoginTrap.ServerName, n)
		traps = append(traps, t)
		log.Info().Str("type", "LoginTrap").Str("addr", cfg.Traps.LoginTrap.Addr).Msg("Trap enabled")
	}

	return traps
}

func waitForShutdown(cancel context.CancelFunc, traps []trap.Trap, wg *sync.WaitGroup) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down VoidSink...")

	cancel() // Signal traps to stop

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

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
	wg.Wait()
	log.Info().Msg("VoidSink shutdown complete")
}
