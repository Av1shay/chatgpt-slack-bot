package gpt

import (
	"context"
	"errors"
	gogpt "github.com/sashabaranov/go-gpt3"
)

var ErrNoChoices = errors.New("no choices found")

type Service struct {
	model       string
	maxTokens   int
	temperature float32
	client      *gogpt.Client
}

type Option func(s *Service)

func WithMaxTokens(m int) Option {
	return func(s *Service) {
		s.maxTokens = m
	}
}

func WithTemperature(t float32) Option {
	return func(s *Service) {
		s.temperature = t
	}
}

func NewService(authToken string, opts ...Option) *Service {
	client := gogpt.NewClient(authToken)
	service := &Service{
		model:       gogpt.GPT3TextDavinci003,
		maxTokens:   100,
		temperature: 0,
		client:      client,
	}
	for _, opt := range opts {
		opt(service)
	}
	return service
}

func (s *Service) Do(ctx context.Context, prompt string) (string, error) {
	req := gogpt.CompletionRequest{
		Model:       s.model,
		MaxTokens:   s.maxTokens,
		Temperature: s.temperature,
		TopP:        1,
		Prompt:      prompt,
	}
	resp, err := s.client.CreateCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", ErrNoChoices
	}
	return resp.Choices[0].Text, nil
}
