package checker

import (
	"log"
	"time"

	"github.com/txltedxgod/cert-sentinel/pkg/alerts"
	"github.com/txltedxgod/cert-sentinel/pkg/metrics"
)

type Sentinel struct {
	config     *Config
	dispatcher *alerts.Dispatcher
}

func NewSentinel(cfg *Config) *Sentinel {
	return &Sentinel{
		config:     cfg,
		dispatcher: alerts.NewDispatcher(cfg.WebhookURL),
	}
}

func (s *Sentinel) Start() {
	s.checkAll()

	ticker := time.NewTicker(s.config.CheckInterval)
	for range ticker.C {
		s.checkAll()
	}
}

func (s *Sentinel) checkAll() {
	log.Println("[cert-sentinel] Running TLS certificate probe checks...")

	for _, target := range s.config.Targets {
		info := CheckTarget(target, 5*time.Second)
		if !info.IsValid || info.Error != nil {
			log.Printf("[Warning] Failed or expired certificate on %s: %v\n", target, info.Error)
			metrics.CertValid.WithLabelValues(target).Set(0)
			metrics.CertExpiryDays.WithLabelValues(target, "unknown", "unknown").Set(0)
			continue
		}

		metrics.CertValid.WithLabelValues(target).Set(1)
		metrics.CertExpiryDays.WithLabelValues(target, info.CommonName, info.Issuer).Set(info.DaysRemaining)

		log.Printf("[OK] %s (%s) - %.1f days remaining\n", target, info.CommonName, info.DaysRemaining)

		if info.DaysRemaining <= s.config.WarnThresholdDays {
			s.dispatcher.SendAlert(target, info.CommonName, info.DaysRemaining)
		}
	}
}
