package application

import (
	"acai_travel/internal/chat/adapters/llm"
	"acai_travel/internal/chat/domain"
	"context"
)

type DestinationExpert struct {
	client domain.LLMClient
}

func NewDestinationExpert(client domain.LLMClient) *DestinationExpert {
	return &DestinationExpert{client: client}
}

func (u *DestinationExpert) Run(
	ctx context.Context,
	chat *domain.Chat,
	injections domain.PromptInjectable,
	model domain.LLMModel,
) (*domain.Chat, error) {
	agent := domain.DestinationExpert
	sessionChat, err := domain.NewAgentSessionFromInjection(agent, chat.UserID, injections)
	if err != nil {
		return nil, err
	}
	sessionChat.AppendMessagesFrom(chat)

	session := llm.NewLLMModelSession(u.client, string(model))

	response, err := session.Chat(ctx, sessionChat.Messages)
	if err != nil {
		return nil, err
	}

	sessionChat.AddMessage(response)

	return sessionChat, nil
}
