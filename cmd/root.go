package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "debug",
    Short: "Kubernetes debugging CLI with LLM integration",
    Long:  `debug helps you debug Kubernetes pods by analyzing logs and showing resource usage with the help of an LLM.`,
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
