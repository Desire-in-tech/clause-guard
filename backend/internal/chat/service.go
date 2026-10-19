package chat

import (
	"context"
	"strings"
)

type Service struct {
	client Client
}

func NewService(client Client) *Service {
	return &Service{client: client}
}

func (s *Service) Ask(ctx context.Context, request Request) (Response, error) {
	request.Question = strings.TrimSpace(request.Question)
	return s.client.Chat(ctx, request)
}
