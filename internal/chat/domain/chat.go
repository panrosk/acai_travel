package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Chat represents a conversation between a user and the system.
type Chat struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	CreatedAt time.Time
	Messages  []Message
}

// Errors related to chat validation
var (
	ErrChatDoesNotBelong        = errors.New("message does not belong to this chat")
	ErrFirstMessageMustBeSystem = errors.New("first message in a chat must be sent by the system")
)

// NewChat creates a new chat instance with a given user ID.
func NewChat(userID uuid.UUID) *Chat {
	return &Chat{
		ID:        uuid.New(),
		UserID:    userID,
		CreatedAt: time.Now().UTC(),
		Messages:  []Message{},
	}
}

// AddMessage validates and adds the message to the chat in place.
// Returns an error if the message doesn't belong to the chat or if the first message is not system-generated.
func (c *Chat) AddMessage(message Message) error {
	if message.ChatID != c.ID {
		return ErrChatDoesNotBelong
	}

	c.Messages = append(c.Messages, message)
	return nil
}

// AppendMessagesFrom copies all messages from another chat into the current one,
// generating new message IDs and assigning the current chat's ID to each copied message.
// Returns an error if validation fails.
func (c *Chat) AppendMessagesFrom(other *Chat) error {
	if len(other.Messages) == 0 {
		return nil // Nothing to append
	}

	for _, m := range other.Messages {
		newMsg := Message{
			ID:        uuid.New(),
			ChatID:    c.ID,
			Sender:    m.Sender,
			Content:   m.Content,
			Timestamp: time.Now().UTC(),
		}
		c.Messages = append(c.Messages, newMsg)
	}

	return nil
}
