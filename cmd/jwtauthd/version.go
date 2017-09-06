package jwtauthd

import (
    "github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
    Use: "version",
    Short: "stop running server",
    Long: "stop running server",
    RunE: func(cmd *cobra.Command, args []string) error {
        
        return nil
    },
}