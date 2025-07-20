package llm

import (
	"acai_travel/internal/chat/domain"
	"context"
)

type LLMModelSession[T domain.LLMClient] struct {
	client T
	model  string
}

func NewLLMModelSession[T domain.LLMClient](client T, model string) *LLMModelSession[T] {
	return &LLMModelSession[T]{
		client: client,
		model:  model,
	}
}

func (s *LLMModelSession[T]) Chat(ctx context.Context, messages []domain.Message) (domain.Message, error) {
	return s.client.Chat(ctx, messages, s.model)
}

func (s *LLMModelSession[T]) StreamChat(ctx context.Context, messages []domain.Message, streamFn func(string) error) error {
	return s.client.StreamChat(ctx, messages, streamFn, s.model)
}

func (s *LLMModelSession[T]) StructuredOutput(
	ctx context.Context,
	messages []domain.Message,
	schema any,
) (map[string]string, error) {
	return s.client.StructuredOutput(ctx, messages, s.model, schema)
}
