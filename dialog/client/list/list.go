package list

import (
	"github.com/spf13/cobra"
	"github.com/BluePecker/JwtAuth/dialog/client"
	"fmt"
	"encoding/json"
	"github.com/BluePecker/JwtAuth/dialog/server/parameter/jwt/request"
	"bytes"
	"github.com/BluePecker/JwtAuth/dialog/server/parameter"
	"github.com/BluePecker/JwtAuth/dialog/server/parameter/jwt/response"
	"github.com/BluePecker/JwtAuth/dialog/client/formatter"
	"github.com/BluePecker/JwtAuth/dialog/client/formatter/context"
	"github.com/BluePecker/JwtAuth/pkg/term"
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
				"application/json;charset=utf-8", bytes.NewBuffer(body)); err != nil {
				return err
			} else {
				var res parameter.Response
				if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&res); err != nil {
					fmt.Println(err)
				} else {
					if buffer, err := json.Marshal(res.Data); err != nil {

					} else {
						var list []response.JsonWebToken
						_, stdout, _ := term.StdStreams()
						json.NewDecoder(bytes.NewBuffer(buffer)).Decode(&list)
						(&formatter.JsonWebTokenContext{
							Context: context.Context{
								Writer:   stdout,
								Template: "table",
							},
							JsonWebTokens: list,
						}).Write()
					}
					//fmt.Printf("result: %s\n", body)
				}
			}

			return nil
		},
	}

	return cmd
}
