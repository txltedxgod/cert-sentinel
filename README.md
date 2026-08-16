# cert-sentinel

> Automated SSL/TLS certificate expiration monitor and Prometheus exporter with Webhook alerting written in **Go**.

[![Go](https://img.shields.io/badge/Go-1.22-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![Prometheus](https://img.shields.io/badge/Metrics-Prometheus-E6522C?style=flat-square&logo=prometheus)](https://prometheus.io)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat-square&logo=docker)](https://docker.com)
[![License](https://img.shields.io/badge/License-MIT-blue?style=flat-square)](LICENSE)

`#ssl-certificate` `#tls-monitoring` `#prometheus-exporter` `#devops` `#observability` `#golang` `#alerting`

---

## Features

- **Automated TLS Probing:** Continuously checks remote hosts and parses leaf certificate validity and expiration dates.
- **Prometheus Exporter:** Serves `/metrics` formatted for easy ingestion into Prometheus & Grafana alerting dashboards.
- **Webhook Alerts:** Triggers configurable JSON alert payloads to Slack, Discord, Telegram or custom webhooks when expiration is within threshold (e.g. < 14 days).
- **Multi-Target Polling:** Probes dozens of domains concurrently.

## Quick Start

### Run with Go

```bash
# Copy sample config
cp config.example.yaml config.yaml

# Run exporter
go run main.go -config=config.yaml -listen=:9109
```

### Docker

```bash
docker build -t cert-sentinel .
docker run -d -p 9109:9109 -v $(pwd)/config.yaml:/etc/cert-sentinel/config.yaml cert-sentinel
```

## Metrics Exported

- `cert_sentinel_expiry_days{target="...", common_name="...", issuer="..."}`: Days remaining until expiration.
- `cert_sentinel_valid{target="..."}`: Gauge (1 = Valid, 0 = Expired or connection error).
