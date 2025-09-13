package llm

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

func AnalyzeLogs(logs string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY not set")
	}

	client := openai.NewClient(apiKey)

	prompt := fmt.Sprintf(`You are a Kubernetes debugging assistant.
Analyze the following pod logs and:
1. Identify possible reasons why the pod is failing.
2. Suggest actionable steps to fix the issue.

Pod logs:
%s
`, logs)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo, // ChatGPT 3.5
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "user",
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func summarize(logs string) string {
	result, err := AnalyzeLogs(logs)
	if err != nil {
		return "Failed to analyze logs: " + err.Error()
	}
	return result
}
