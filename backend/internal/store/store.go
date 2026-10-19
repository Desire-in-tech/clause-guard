package store

import (
	"context"

	"clauseguard/backend/internal/analysis"
)

type Store interface {
	SaveReport(context.Context, *analysis.Report) error
	GetReport(context.Context, string) (*analysis.Report, error)
}
