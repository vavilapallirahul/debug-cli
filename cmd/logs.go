package cmd

import (
	"context"
	"fmt"

	"github.com/vavilapallirahul/debug-cli/k8s"
	"github.com/vavilapallirahul/debug-cli/llm"

	"github.com/spf13/cobra"
)

var namespace string

var logsCmd = &cobra.Command{
	Use:   "logs [pod]",
	Short: "Analyze logs of a pod with LLM",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		podName := args[0]

		logs, err := k8s.GetPodLogs(context.Background(), podName, namespace)
		if err != nil {
			if k8s.IsNotFoundErr(err) {
				return fmt.Errorf("pod %q not found in namespace %q", podName, namespace)
			}
			if k8s.IsNamespaceNotFoundErr(err) {
				return fmt.Errorf("namespace %q not found", namespace)
			}
			return fmt.Errorf("failed to get logs: %w", err)
		}

		analysis, err := llm.AnalyzeLogs(string(logs))
		if err != nil {
			return fmt.Errorf("LLM analysis failed: %w", err)
		}

		fmt.Println("LLM Analysis:\n", analysis)
		return nil
	},
}

func init() {
	logsCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Pod namespace (required)")
	logsCmd.MarkFlagRequired("namespace")
	rootCmd.AddCommand(logsCmd)
}
