package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/vavilapallirahul/debug-cli/k8s"
	"github.com/vavilapallirahul/debug-cli/llm"

	"github.com/spf13/cobra"
)

var namespace string
var maxLines int

func init() {
	logsCmd := &cobra.Command{
		Use:   "logs [pod]",
		Short: "Analyze logs of a pod with LLM suggestions",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			podName := args[0]

			// Fetch pod logs
			logs, err := k8s.GetPodLogs(context.Background(), podName, namespace)
			if err != nil {
				return fmt.Errorf("failed to get logs for pod %q in namespace %q: %w", podName, namespace, err)
			}

			// Take last N lines
			// Convert logs from []byte to string
			logStr := string(logs)

			// Take last N lines
			lines := strings.Split(logStr, "\n")

			if len(lines) > maxLines {
				lines = lines[len(lines)-maxLines:]
			}
			logSnippet := strings.Join(lines, "\n")

			// Send to LLM for analysis
			analysis, err := llm.AnalyzeLogs(logSnippet)
			if err != nil {
				return fmt.Errorf("LLM analysis failed: %w", err)
			}

			fmt.Println("LLM Analysis:\n", analysis)
			return nil
		},
	}

	logsCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Pod namespace (required)")
	logsCmd.MarkFlagRequired("namespace")

	logsCmd.Flags().IntVarP(&maxLines, "lines", "l", 100, "Number of last log lines to analyze")

	rootCmd.AddCommand(logsCmd)
}
