package list

import (
	"github.com/spf13/cobra"
	"github.com/BluePecker/JwtAuth/dialog/client"
	"fmt"
	"net/http"
	"strings"
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
			var r http.Request
			r.ParseForm()
			fmt.Println(args)
			r.Form.Add("unique", args[0])
			body := strings.TrimSpace(r.Form.Encode())
			cli := client.NewClient(unixSock)
			if body, err := cli.Post("/v1.0/token/list",
				"application/json", strings.NewReader(body)); err != nil {
				return err
			} else {
				fmt.Printf("result: %s\n", body)
			}

			return nil
		},
	}

	return cmd
}
