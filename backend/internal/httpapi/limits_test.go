package httpapi

import (
	"net/http/httptest"
	"testing"
	"time"
)

func TestRequestLimiterDailyBudgetResetsAtUTCMidnight(t *testing.T) {
	limiter := newRequestLimiter(2, 10)
	first := time.Date(2026, 10, 21, 23, 59, 0, 0, time.UTC)

	if got := limiter.allow("192.0.2.1", first); got != limitAllowed {
		t.Fatalf("first request = %v, want allowed", got)
	}
	if got := limiter.allow("192.0.2.2", first); got != limitAllowed {
		t.Fatalf("second request = %v, want allowed", got)
	}
	if got := limiter.allow("192.0.2.3", first); got != limitDailyBudget {
		t.Fatalf("third request = %v, want daily limit", got)
	}
	if got := limiter.allow("192.0.2.3", first.Add(2*time.Minute)); got != limitAllowed {
		t.Fatalf("request after UTC midnight = %v, want allowed", got)
	}
}

func TestRequestLimiterPerIPHourlyLimit(t *testing.T) {
	limiter := newRequestLimiter(10, 1)
	now := time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC)
	if got := limiter.allow("192.0.2.1", now); got != limitAllowed {
		t.Fatalf("first request = %v, want allowed", got)
	}
	if got := limiter.allow("192.0.2.1", now.Add(time.Minute)); got != limitHourlyIP {
		t.Fatalf("second request = %v, want hourly limit", got)
	}
	if got := limiter.allow("192.0.2.1", now.Add(time.Hour+time.Second)); got != limitAllowed {
		t.Fatalf("request after rolling hour = %v, want allowed", got)
	}
	if got := limiter.allow("192.0.2.2", now.Add(2*time.Minute)); got != limitAllowed {
		t.Fatalf("different IP request = %v, want allowed", got)
	}
}

func TestClientIPUsesFirstForwardedAddress(t *testing.T) {
	request := httptest.NewRequest("POST", "/api/analyze", nil)
	request.RemoteAddr = "192.0.2.10:43123"
	request.Header.Set("X-Forwarded-For", "198.51.100.4, 198.51.100.5")

	if got := clientIP(request); got != "198.51.100.4" {
		t.Fatalf("clientIP() = %q, want first forwarded address", got)
	}
}

func TestClientIPFallsBackToRemoteAddress(t *testing.T) {
	request := httptest.NewRequest("POST", "/api/analyze", nil)
	request.RemoteAddr = "192.0.2.10:43123"

	if got := clientIP(request); got != "192.0.2.10" {
		t.Fatalf("clientIP() = %q, want remote IP", got)
	}
}
