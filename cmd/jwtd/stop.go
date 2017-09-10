package jwtd

import (
    "github.com/spf13/cobra"
)

var StopCmd = &cobra.Command{
    Use: "stop",
    Short: "stop running server",
    Long: "stop running server",
    RunE: func(cmd *cobra.Command, args []string) error {
        
        return nil
    },
}
