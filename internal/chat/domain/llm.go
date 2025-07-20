// This implemenation can be extended to use tools, or structured outputs or simple completions
package domain

import "context"

type LLMModel string

// LLMProvider defines the expected behavior from a Large Language Model provider.
type LLMClient interface {
	StructuredOutput(ctx context.Context, messages []Message, model string, schema any) (map[string]string, error)
	Chat(ctx context.Context, messages []Message, model string) (Message, error)
	StreamChat(ctx context.Context, messages []Message, streamFn func(string) error, model string) error
}
