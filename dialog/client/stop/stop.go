package stop

import (
	"github.com/spf13/cobra"
	"github.com/BluePecker/JwtAuth/dialog/client"
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
			body, err := cli.Get("/v1.0/signal/stop")
			if err != nil {
				return err
			}
			fmt.Printf("result: %s", body)
			return nil
		},
	}

	return cmd
}
