package chat

import "context"

type ClauseReference struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Question string            `json:"question"`
	Clauses  []ClauseReference `json:"clauses"`
	History  []Message         `json:"history"`
}

type Response struct {
	Answer    string   `json:"answer"`
	Citations []string `json:"citations"`
}

type Client interface {
	Chat(context.Context, Request) (Response, error)
}
