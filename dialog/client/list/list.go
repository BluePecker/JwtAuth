package list

import (
	"github.com/spf13/cobra"
	"github.com/BluePecker/JwtAuth/dialog/client"
	"fmt"
	"strings"
	"encoding/json"
	"github.com/BluePecker/JwtAuth/dialog/server/parameter/token/request"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "list user's json-web-token",
		RunE: func(cmd *cobra.Command, args []string) error {
			unixSock, err := cmd.Parent().Flags().GetString("unix-sock")
			if err != nil {
				return err
			}
			body, _ := json.Marshal(request.List{Unique: args[0]})
			cli := client.NewClient(unixSock)
			if body, err := cli.Post("/v1.0/token/list",
				"application/json", strings.NewReader(string(body))); err != nil {
				return err
			} else {
				fmt.Printf("result: %s\n", body)
			}

			return nil
		},
	}

	return cmd
}
