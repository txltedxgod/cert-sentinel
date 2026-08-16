package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	CertExpiryDays = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cert_sentinel_expiry_days",
			Help: "Number of days remaining until certificate expiration",
		},
		[]string{"target", "common_name", "issuer"},
	)

	CertValid = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cert_sentinel_valid",
			Help: "1 if SSL certificate is currently valid and unexpired, 0 otherwise",
		},
		[]string{"target"},
	)
)
