package cmd

import (
	"context"
	"fmt"

	"github.com/vavilapallirahul/debug-cli/k8s"

	"github.com/spf13/cobra"
)

var cpuNamespace string
var topCount int

var cpuCmd = &cobra.Command{
	Use:   "cpu",
	Short: "Show top pods consuming CPU in a namespace",
	RunE: func(cmd *cobra.Command, args []string) error {
		pods, err := k8s.GetTopPodsByCPU(context.Background(), cpuNamespace, topCount)
		if err != nil {
			return fmt.Errorf("failed to get CPU usage (is metrics-server running?): %w", err)
		}

		fmt.Printf("Top %d pods in namespace %q by CPU usage:\n", topCount, cpuNamespace)
		for i, p := range pods {
			fmt.Printf("%d. %s\t%s\n", i+1, p.Name, p.CPU)
		}
		return nil
	},
}

func init() {
	cpuCmd.Flags().StringVarP(&cpuNamespace, "namespace", "n", "", "Namespace (required)")
	cpuCmd.Flags().IntVarP(&topCount, "count", "c", 5, "Number of pods to display")
	cpuCmd.MarkFlagRequired("namespace")
	rootCmd.AddCommand(cpuCmd)
}
