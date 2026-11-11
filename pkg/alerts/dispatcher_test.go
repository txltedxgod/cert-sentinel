package alerts

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestDispatcher_SendAlert(t *testing.T) {
	var received atomic.Int32
	var got AlertPayload

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received.Add(1)
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Errorf("failed to decode alert payload: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	noOp := NewDispatcher("")
	noOp.SendAlert("example.com:443", "example.com", 5.0)
	if received.Load() != 0 {
		t.Fatalf("expected empty webhookURL not to send requests, got %d", received.Load())
	}

	d := NewDispatcher(srv.URL)
	d.SendAlert("example.com:443", "example.com", 5.5)

	if received.Load() != 1 {
		t.Fatalf("expected 1 alert request, got %d", received.Load())
	}
	if got.Target != "example.com:443" || got.CommonName != "example.com" || got.DaysRemaining != 5.5 {
		t.Fatalf("unexpected alert payload: %+v", got)
	}
	if got.Timestamp == "" || got.Message == "" {
		t.Fatalf("expected populated timestamp and message, got %+v", got)
	}
}
