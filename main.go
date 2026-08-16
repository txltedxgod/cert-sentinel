package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/txltedxgod/cert-sentinel/pkg/checker"
	"github.com/txltedxgod/cert-sentinel/pkg/metrics"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to configuration yaml")
	listenAddr := flag.String("listen", ":9109", "Metrics server address")
	flag.Parse()

	log.Printf("[cert-sentinel] Starting exporter on %s\n", *listenAddr)

	cfg, err := checker.LoadConfig(*configPath)
	if err != nil {
		log.Printf("Warning: could not load %s, using defaults: %v\n", *configPath, err)
		cfg = &checker.Config{
			Targets:       []string{"google.com:443", "github.com:443"},
			CheckInterval: 60 * time.Second,
			WarnThresholdDays: 14,
		}
	}

	sentinel := checker.NewSentinel(cfg)
	go sentinel.Start()

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK\n"))
	})

	server := &http.Server{
		Addr:         *listenAddr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("[cert-sentinel] Shutting down...")
}
