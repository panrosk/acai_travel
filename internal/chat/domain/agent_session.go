package domain

import (
	"github.com/google/uuid"
)

// NewAgentSessionFromInjection creates a new Chat initialized with the agent's system prompt
// generated from the provided PromptInjectable.
func NewAgentSessionFromInjection(agent Agent, userID uuid.UUID, injection PromptInjectable) (*Chat, error) {
	systemPrompt, err := injection.ToPrompt(agent)
	if err != nil {
		return nil, err
	}

	chat := NewChat(userID)

	systemMessage := NewSystemMessage(chat.ID, systemPrompt)

	chat.AddMessage(systemMessage)

	return chat, nil
}
