package llm

import (
	"acai_travel/internal/chat/domain"
	"fmt"
	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
)

func convertToOpenAIMessages(messages []domain.Message) []openai.ChatCompletionMessageParamUnion {
	var converted []openai.ChatCompletionMessageParamUnion
	for _, m := range messages {
		switch m.Sender {
		case domain.SenderUser:
			converted = append(converted, openai.UserMessage(m.Content))
		case domain.SenderSystem:
			converted = append(converted, openai.SystemMessage(m.Content))
		case domain.SenderAI:
			converted = append(converted, openai.AssistantMessage(m.Content))
		}
	}
	return converted
}

func mapModel(model string) (string, error) {
	switch model {
	case "gpt-4":
		return openai.ChatModelGPT4, nil
	case "gpt-4o":
		return openai.ChatModelGPT4o, nil
	case "gpt-3.5":
		return openai.ChatModelGPT3_5Turbo, nil

	default:
		return "", fmt.Errorf("unsupported OpenAI model: %s", model)
	}
}

func GenerateSchema[T any]() interface{} {
	// Structured Outputs uses a subset of JSON schema
	// These flags are necessary to comply with the subset
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}
