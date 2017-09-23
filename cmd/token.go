package cmd

import (
    "github.com/spf13/cobra"
)

var TokenCmd = &cobra.Command{
    Use: "token",
    Short: "stop running server",
    Long: "stop running server",
    RunE: func(cmd *cobra.Command, args []string) error {
        
        return nil
    },
}