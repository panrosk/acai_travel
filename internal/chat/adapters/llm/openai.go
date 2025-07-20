package llm

import (
	"acai_travel/internal/chat/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAIClient struct {
	client *openai.Client
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	client := openai.NewClient(option.WithAPIKey(apiKey))
	return &OpenAIClient{
		client: &client,
	}
}

func (o *OpenAIClient) Chat(ctx context.Context, messages []domain.Message, model string) (domain.Message, error) {
	model, err := mapModel(model)
	if err != nil {
		return domain.Message{}, err
	}

	resp, err := o.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: convertToOpenAIMessages(messages),
		Model:    model,
	})
	if err != nil {
		fmt.Println(err)
		return domain.Message{}, err
	}

	if len(resp.Choices) == 0 {
		return domain.Message{}, errors.New("no choices returned by OpenAI")
	}

	return domain.NewAIMessage(messages[0].ChatID, resp.Choices[0].Message.Content), nil
}

func (o *OpenAIClient) StreamChat(
	ctx context.Context,
	messages []domain.Message,
	streamFn func(string) error,
	model string,
) error {
	model, err := mapModel(model)
	if err != nil {
		return err
	}

	stream := o.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Messages: convertToOpenAIMessages(messages),
		Model:    model,
	})
	defer stream.Close()

	for stream.Next() {
		chunk := stream.Current()
		if len(chunk.Choices) == 0 {
			continue
		}

		content := chunk.Choices[0].Delta.Content
		if content == "" {
			continue
		}

		if err := streamFn(content); err != nil {
			return err
		}
	}

	return stream.Err()
}

func (o *OpenAIClient) StructuredOutput(
	ctx context.Context,
	messages []domain.Message,
	model string,
	schema any,
) (map[string]string, error) {
	if schema == nil {
		return nil, errors.New("structured output: schema is nil")
	}

	mappedModel, err := mapModel(model)
	if err != nil {
		return nil, err
	}

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "structured_response",
		Description: openai.String("Structured JSON response"),
		Schema:      schema,
		Strict:      openai.Bool(true),
	}

	oaMessages := convertToOpenAIMessages(messages)

	resp, err := o.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: oaMessages,
		Model:    mappedModel,
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: schemaParam,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("structured output request failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, errors.New("structured output: no choices returned")
	}

	content := resp.Choices[0].Message.Content
	if content == "" {
		return nil, errors.New("structured output: empty content")
	}

	var result map[string]string
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("structured output: invalid JSON returned (unmarshal failed): %w; raw=%s", err, content)
	}

	return result, nil
}
