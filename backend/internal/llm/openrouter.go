package llm

import (
	"context"
	"errors"

	"clauseguard/backend/internal/analysis"
	"clauseguard/backend/internal/chat"
	"clauseguard/backend/internal/contracts"
)

var ErrNotImplemented = errors.New("OpenRouter is not implemented in this step")

type OpenRouterClient struct{}

func NewOpenRouterClient() *OpenRouterClient {
	return &OpenRouterClient{}
}

func (c *OpenRouterClient) Analyze(context.Context, []contracts.Clause) (analysis.AnalysisResult, error) {
	return analysis.AnalysisResult{}, ErrNotImplemented
}

func (c *OpenRouterClient) Chat(context.Context, chat.Request) (chat.Response, error) {
	return chat.Response{}, ErrNotImplemented
}
