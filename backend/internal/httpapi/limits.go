package httpapi

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type limitKind int

const (
	limitAllowed limitKind = iota
	limitDailyBudget
	limitHourlyIP
)

type requestLimiter struct {
	mu             sync.Mutex
	dailyBudget    int
	perIPPerHour   int
	currentUTCDate string
	dailyCount     int
	requestsByIP   map[string][]time.Time
}

func newRequestLimiter(dailyBudget, perIPPerHour int) *requestLimiter {
	if dailyBudget <= 0 {
		dailyBudget = 40
	}
	if perIPPerHour <= 0 {
		perIPPerHour = 10
	}
	return &requestLimiter{
		dailyBudget:  dailyBudget,
		perIPPerHour: perIPPerHour,
		requestsByIP: make(map[string][]time.Time),
	}
}

func (l *requestLimiter) allow(ip string, now time.Time) limitKind {
	l.mu.Lock()
	defer l.mu.Unlock()

	now = now.UTC()
	date := now.Format("2006-01-02")
	if l.currentUTCDate != date {
		l.currentUTCDate = date
		l.dailyCount = 0
	}
	if l.dailyCount >= l.dailyBudget {
		return limitDailyBudget
	}

	cutoff := now.Add(-time.Hour)
	recent := l.requestsByIP[ip][:0]
	for _, requestedAt := range l.requestsByIP[ip] {
		if requestedAt.After(cutoff) {
			recent = append(recent, requestedAt)
		}
	}
	if len(recent) >= l.perIPPerHour {
		l.requestsByIP[ip] = recent
		return limitHourlyIP
	}

	l.dailyCount++
	l.requestsByIP[ip] = append(recent, now)
	return limitAllowed
}

func clientIP(request *http.Request) string {
	if forwarded := request.Header.Get("X-Forwarded-For"); forwarded != "" {
		if first := strings.TrimSpace(strings.SplitN(forwarded, ",", 2)[0]); first != "" {
			return first
		}
	}
	if host, _, err := net.SplitHostPort(request.RemoteAddr); err == nil {
		return host
	}
	return request.RemoteAddr
}
