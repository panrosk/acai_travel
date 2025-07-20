package config

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func NewOpenAiClient(apiKey string) openai.Client {

	client := openai.NewClient(
		option.WithAPIKey(apiKey), // defaults to os.LookupEnv("OPENAI_API_KEY")
	)
	return client
}
