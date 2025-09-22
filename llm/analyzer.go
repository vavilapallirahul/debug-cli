package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// func AnalyzeLogs(logs string) (string, error) {
// 	apiKey := os.Getenv("OPENAI_API_KEY")
// 	if apiKey == "" {
// 		return "", fmt.Errorf("OPENAI_API_KEY not set")
// 	}

// 	client := openai.NewClient(apiKey)

// 	prompt := fmt.Sprintf(`You are a Kubernetes debugging assistant.
// Analyze the following pod logs and:
// 1. Identify possible reasons why the pod is failing.
// 2. Suggest actionable steps to fix the issue.

// Pod logs:
// %s
// `, logs)

// 	resp, err := client.CreateChatCompletion(
// 		context.Background(),
// 		openai.ChatCompletionRequest{
// 			Model: openai.GPT3Dot5Turbo, // ChatGPT 3.5
// 			Messages: []openai.ChatCompletionMessage{
// 				{
// 					Role:    "user",
// 					Content: prompt,
// 				},
// 			},
// 		},
// 	)
// 	if err != nil {
// 		return "", err
// 	}

// 	return resp.Choices[0].Message.Content, nil
// }

func AnalyzeLogs(logs string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY not set")
	}

	prompt := fmt.Sprintf(`You are a Kubernetes debugging assistant.
Analyze the following pod logs and:
1. Identify possible reasons why the pod is failing.
2. Suggest actionable steps to fix the issue.

Pod logs:
%s
`, logs)

	body := map[string]interface{}{
		"model":    "sonar-pro",
		"messages": []map[string]string{{"role": "user", "content": prompt}},
		"stream":   false,
	}
	b, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(context.Background(), "POST", "https://api.perplexity.ai/chat/completions", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from LLM")
	}
	return result.Choices[0].Message.Content, nil
}

func summarize(logs string) string {
	result, err := AnalyzeLogs(logs)
	if err != nil {
		return "Failed to analyze logs: " + err.Error()
	}
	return result
}
