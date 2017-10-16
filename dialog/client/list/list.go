package list

import (
	"github.com/spf13/cobra"
	"github.com/BluePecker/JwtAuth/dialog/client"
	"fmt"
	"encoding/json"
	"github.com/BluePecker/JwtAuth/dialog/server/parameter/token/request"
	"bytes"
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
			fmt.Printf("post: %s\n", body)
			cli := client.NewClient(unixSock)
			if body, err := cli.Post("/v1.0/token/list",
				"application/json;charset=utf-8", bytes.NewBuffer([]byte(body))); err != nil {
				return err
			} else {
				fmt.Printf("result: %s\n", body)
			}

			return nil
		},
	}

	return cmd
}
