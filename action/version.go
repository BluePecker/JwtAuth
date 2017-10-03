package action

import (
    "github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
    Use: "version",
    Short: "show the Jwt version information",
    Long: "show the Jwt version information",
    RunE: func(cmd *cobra.Command, args []string) error {
        
        return nil
    },
}