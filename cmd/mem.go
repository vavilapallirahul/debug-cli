package cmd

import (
	"context"
	"fmt"

	"github.com/vavilapallirahul/debug-cli/k8s"

	"github.com/spf13/cobra"
)

var memNamespace string
var memCount int

var memCmd = &cobra.Command{
	Use:   "mem",
	Short: "Show top pods consuming memory in a namespace",
	RunE: func(cmd *cobra.Command, args []string) error {
		pods, err := k8s.GetTopPodsByMemory(context.Background(), memNamespace, memCount)
		if err != nil {
			return fmt.Errorf("failed to get memory usage (is metrics-server running?): %w", err)
		}

		fmt.Printf("Top %d pods in namespace %q by memory usage:\n", memCount, memNamespace)
		for i, p := range pods {
			fmt.Printf("%d. %s\t%s\n", i+1, p.Name, p.Memory)
		}
		return nil
	},
}

func init() {
	memCmd.Flags().StringVarP(&memNamespace, "namespace", "n", "", "Namespace (required)")
	memCmd.Flags().IntVarP(&memCount, "count", "c", 5, "Number of pods to display")
	memCmd.MarkFlagRequired("namespace")
	rootCmd.AddCommand(memCmd)
}
