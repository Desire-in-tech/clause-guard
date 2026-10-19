package store

import (
	"context"
	"errors"
	"sync"
	"time"

	"clauseguard/backend/internal/analysis"
)

var ErrReportNotFound = errors.New("report not found")

type MemoryStore struct {
	mu      sync.RWMutex
	reports map[string]analysis.Report
	ttl     time.Duration
	now     func() time.Time
}

func NewMemoryStore(ttl time.Duration) *MemoryStore {
	return &MemoryStore{
		reports: make(map[string]analysis.Report),
		ttl:     ttl,
		now:     time.Now,
	}
}

func (s *MemoryStore) SaveReport(_ context.Context, report *analysis.Report) error {
	id, err := newReportID()
	if err != nil {
		return err
	}
	now := s.now().UTC()
	report.ID = id
	if report.CreatedAt.IsZero() {
		report.CreatedAt = now
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.removeExpired(now)
	s.reports[id] = cloneReport(*report)
	return nil
}

func (s *MemoryStore) GetReport(_ context.Context, id string) (*analysis.Report, error) {
	now := s.now().UTC()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.removeExpired(now)
	report, ok := s.reports[id]
	if !ok {
		return nil, ErrReportNotFound
	}
	copy := cloneReport(report)
	return &copy, nil
}

func (s *MemoryStore) removeExpired(now time.Time) {
	for id, report := range s.reports {
		if !report.CreatedAt.IsZero() && now.Sub(report.CreatedAt) >= s.ttl {
			delete(s.reports, id)
		}
	}
}

func cloneReport(report analysis.Report) analysis.Report {
	report.Clauses = append([]analysis.Clause(nil), report.Clauses...)
	return report
}
