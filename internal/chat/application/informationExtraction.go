package application

import (
	"acai_travel/internal/chat/adapters/llm"
	"acai_travel/internal/chat/domain"
	"context"
)

type InformationExtractor struct {
	client domain.LLMClient
}

func NewInformationExtractor(client domain.LLMClient) *InformationExtractor {
	return &InformationExtractor{client: client}
}

func (u *InformationExtractor) Run(
	ctx context.Context,
	chat *domain.Chat,
	schema map[string]any,
	model domain.LLMModel,
) (map[string]string, error) {

	session := llm.NewLLMModelSession(u.client, string(model))

	result, err := session.StructuredOutput(ctx, chat.Messages, schema)
	if err != nil {
		return nil, err
	}

	return result, nil
}
