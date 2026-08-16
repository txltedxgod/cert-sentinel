package checker

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"
)

type CertInfo struct {
	Target       string
	CommonName   string
	Issuer       string
	DNSNames     []string
	NotBefore    time.Time
	NotAfter     time.Time
	DaysRemaining float64
	IsValid      bool
	Error        error
}

func CheckTarget(target string, timeout time.Duration) *CertInfo {
	host, port, err := net.SplitHostPort(target)
	if err != nil {
		host = target
		port = "443"
		target = net.JoinHostPort(host, port)
	}

	dialer := &net.Dialer{
		Timeout: timeout,
	}

	conn, err := tls.DialWithDialer(dialer, "tcp", target, &tls.Config{
		ServerName: host,
		InsecureSkipVerify: false,
	})

	if err != nil {
		return &CertInfo{
			Target:  target,
			IsValid: false,
			Error:   err,
		}
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return &CertInfo{
			Target:  target,
			IsValid: false,
			Error:   fmt.Errorf("no peer certificates presented"),
		}
	}

	leaf := certs[0]
	now := time.Now()
	daysRemaining := leaf.NotAfter.Sub(now).Hours() / 24.0

	return &CertInfo{
		Target:        target,
		CommonName:    leaf.Subject.CommonName,
		Issuer:        leaf.Issuer.CommonName,
		DNSNames:      leaf.DNSNames,
		NotBefore:     leaf.NotBefore,
		NotAfter:      leaf.NotAfter,
		DaysRemaining: daysRemaining,
		IsValid:       daysRemaining > 0,
	}
}
