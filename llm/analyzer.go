package llm

import "fmt"

// Dummy LLM call — replace with real API call (Ollama, OpenAI, etc.)
func AnalyzeLogs(logs string) (string, error) {
	if len(logs) == 0 {
		return "No logs available to analyze.", nil
	}

	// TODO: connect to LLM provider
	return fmt.Sprintf("Simulated LLM analysis: pod failed due to '%s'", summarize(logs)), nil
}

func summarize(logs string) string {
	if len(logs) > 80 {
		return logs[:80] + "..."
	}
	return logs
}
