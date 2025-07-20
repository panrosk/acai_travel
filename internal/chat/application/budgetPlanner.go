package application

import (
	"acai_travel/internal/chat/adapters/llm"
	"acai_travel/internal/chat/domain"
	"context"
)

type BudgetPlanner struct {
	client domain.LLMClient
}

func NewBudgetPlanner(client domain.LLMClient) *BudgetPlanner {
	return &BudgetPlanner{client: client}
}

func (u *BudgetPlanner) Run(
	ctx context.Context,
	chat *domain.Chat,
	injections domain.PromptInjectable,
	model domain.LLMModel,
) (*domain.Chat, error) {
	agent := domain.BudgetPlanner

	sessionChat, err := domain.NewAgentSessionFromInjection(agent, chat.UserID, injections)
	if err != nil {
		return nil, err
	}

	if err := sessionChat.AppendMessagesFrom(chat); err != nil {
		return nil, err
	}

	session := llm.NewLLMModelSession(u.client, string(model))

	response, err := session.Chat(ctx, sessionChat.Messages)
	if err != nil {
		return nil, err
	}

	if err := sessionChat.AddMessage(response); err != nil {
		return nil, err
	}

	return sessionChat, nil
}
