package analysis

import (
	"context"
	"errors"
	"strings"
	"time"

	"clauseguard/backend/internal/contracts"
)

var ErrEmptyContract = errors.New("contract text must contain at least one clause")

type Service struct {
	analyzer Analyzer
}

func NewService(analyzer Analyzer) *Service {
	return &Service{analyzer: analyzer}
}

func (s *Service) Analyze(ctx context.Context, contractText string) (*Report, error) {
	if strings.TrimSpace(contractText) == "" {
		return nil, ErrEmptyContract
	}
	clauses := contracts.Split(contractText)
	if len(clauses) == 0 {
		return nil, ErrEmptyContract
	}

	result, err := s.analyzer.Analyze(ctx, clauses)
	if err != nil {
		return nil, err
	}
	return &Report{
		Summary:     result.Summary,
		OverallRisk: result.OverallRisk,
		Clauses:     result.Clauses,
		CreatedAt:   time.Now().UTC(),
	}, nil
}
