package checker

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCheckTarget_UntrustedLocalTLSServer(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	addr := strings.TrimPrefix(srv.URL, "https://")
	info := CheckTarget(addr, 2*time.Second)

	if info.Target != addr {
		t.Fatalf("expected target %q, got %q", addr, info.Target)
	}
	if info.IsValid {
		t.Fatal("expected self-signed test TLS server not to be treated as trusted")
	}
	if info.Error == nil {
		t.Fatal("expected untrusted self-signed certificate verification error")
	}
}

func TestCheckTarget_DefaultsPortOnUnreachableHost(t *testing.T) {
	info := CheckTarget("127.0.0.1", 100*time.Millisecond)
	if info.Target != "127.0.0.1:443" {
		t.Fatalf("expected default port 443 added, got %q", info.Target)
	}
	if info.IsValid || info.Error == nil {
		t.Fatalf("expected failure for closed local target, got valid=%v err=%v", info.IsValid, info.Error)
	}
}
