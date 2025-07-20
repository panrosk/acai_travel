package tests

import (
	"acai_travel/internal/chat/adapters/llm"
	"acai_travel/internal/chat/domain"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func createClient(t *testing.T) *llm.OpenAIClient {
	godotenv.Load("../.env")
	apiKey := os.Getenv("OPENAI_API_KEY")
	fmt.Println(apiKey)
	if apiKey == "" {
		t.Skip("OPENAI_API_KEY not set, skipping integration test")
	}
	client := llm.NewOpenAIClient(apiKey)
	return client
}

func TestIntegration_Chat(t *testing.T) {
	client := createClient(t)
	ctx := context.Background()

	chatID := uuid.New()
	msg := domain.NewUserMessage(chatID, "Hello, who won the World Cup in 2018?")
	t.Logf("Sending message to model gpt-4: %s", msg.Content)

	response, err := client.Chat(ctx, []domain.Message{msg}, "gpt-4")
	if err != nil {
		t.Fatalf("Chat failed: %v", err)
	}

	t.Logf("Received response: %s", response.Content)

	assert.Equal(t, domain.SenderAI, response.Sender)
	assert.NotEmpty(t, response.Content)
}

func TestIntegration_StructuredOutput(t *testing.T) {
	client := createClient(t)
	ctx := context.Background()

	chatID := uuid.New()
	msg := domain.NewUserMessage(chatID, "I want to visit Italy. I'm on a budget and I love nature.")
	t.Logf("Sending structured message to model gpt-4: %s", msg.Content)

	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"Destinations": map[string]string{"type": "string"},
			"Preferences":  map[string]string{"type": "string"},
			"Interest":     map[string]string{"type": "string"},
		},
		"required":             []string{"Destinations", "Preferences", "Interest"},
		"additionalProperties": false,
	}

	response, err := client.StructuredOutput(ctx, []domain.Message{msg}, "gpt-4o", schema)
	if err != nil {
		t.Fatalf("StructuredOutput failed: %v", err)
	}

	t.Logf("Structured response content: %s", response.Content)

	var decoded map[string]string
	err = json.Unmarshal([]byte(response.Content), &decoded)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	t.Logf("Parsed JSON: %+v", decoded)

	assert.NotEmpty(t, decoded["Destinations"])
	assert.NotEmpty(t, decoded["Preferences"])
	assert.NotEmpty(t, decoded["Interest"])
}

func TestIntegration_StreamChat(t *testing.T) {
	client := createClient(t)
	ctx := context.Background()

	chatID := uuid.New()
	msg := domain.NewUserMessage(chatID, "Can you tell me a short story about a robot and a cat?")
	t.Logf("Streaming message to model gpt-3.5: %s", msg.Content)

	var fullResponse string

	err := client.StreamChat(ctx, []domain.Message{msg}, func(content string) error {
		fullResponse += content
		return nil
	}, "gpt-4")

	if err != nil {
		t.Fatalf("StreamChat failed: %v", err)
	}

	t.Logf("Full streamed response: %s", fullResponse)

	assert.NotEmpty(t, fullResponse)
}
