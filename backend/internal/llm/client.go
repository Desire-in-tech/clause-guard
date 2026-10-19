package llm

import (
	"context"

	"clauseguard/backend/internal/analysis"
	"clauseguard/backend/internal/chat"
	"clauseguard/backend/internal/contracts"
)

type Client interface {
	Analyze(context.Context, []contracts.Clause) (analysis.AnalysisResult, error)
	Chat(context.Context, chat.Request) (chat.Response, error)
}
