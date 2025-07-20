package application

import (
	"acai_travel/internal/chat/adapters/llm"
	"acai_travel/internal/chat/domain"
	"context"
)

type TripSynthesizer struct {
	client domain.LLMClient
}

func NewTripSynthesizer(client domain.LLMClient) *TripSynthesizer {
	return &TripSynthesizer{client: client}
}

func (u *TripSynthesizer) Stream(
	ctx context.Context,
	chat *domain.Chat,
	injections domain.PromptInjectable,
	model domain.LLMModel,
	streamFn func(eventType, data string) error,
) error {
	agent := domain.TripSynthesizer
	sessionChat, err := domain.NewAgentSessionFromInjection(agent, chat.UserID, injections)
	if err != nil {
		return err
	}
	sessionChat.AppendMessagesFrom(chat)

	session := llm.NewLLMModelSession(u.client, string(model))

	err = session.StreamChat(ctx, chat.Messages, WrapMessageStreamer(streamFn))
	if err != nil {
		return err
	}

	return nil
}
