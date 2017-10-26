package stop

import (
	"github.com/spf13/cobra"
	"github.com/BluePecker/JwtAuth/engine/client"
	"fmt"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "stop running server",
		RunE: func(cmd *cobra.Command, args []string) error {
			unixSock, err := cmd.Parent().Flags().GetString("unix-sock")
			if err != nil {
				return err
			}
			cli := client.NewClient(unixSock)
			if _, err := cli.Get("/v1.0/signal/stop"); err != nil {
				return err
			}
			fmt.Printf("successfully terminated the process.\n")
			return nil
		},
	}

	return cmd
}
