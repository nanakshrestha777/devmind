package ai

import (
	"context"
	"fmt"

	"github.com/ollama/ollama/api"
)

func Explain(contextData string) (string, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return "", fmt.Errorf("failed to create Ollama client: %w", err)
	}

	prompt := "You are a Senior Software Engineer. Explain the architecture of the code based on these relationships: " + contextData
	req := &api.GenerateRequest{
		Model:  "llama2",
		Prompt: prompt,
	}

	var response string

	fn := func(resp api.GenerateResponse) error {
		response = resp.Response
		return nil
	}

	err = client.Generate(context.Background(), req, fn)
	return response, err
}
