package kick

import (
	"github.com/spf13/cobra"
	"fmt"
	"github.com/BluePecker/JwtAuth/engine/client"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kick",
		Short: "force users to go offline",
		RunE: func(cmd *cobra.Command, args []string) error {
			unixSock, err := cmd.Parent().Flags().GetString("unix-sock")
			if err != nil {
				return err
			}
			cli := client.NewClient(unixSock)
			if body, err := cli.Get("/v1.0/token/kick"); err != nil {
				return err
			} else {
				fmt.Println(body)
				fmt.Printf("successfully kicked out the user.\n")
				return nil
			}
		},
	}
	return cmd
}
