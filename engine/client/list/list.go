package list

import (
	"github.com/spf13/cobra"
	"github.com/BluePecker/JwtAuth/engine/client"
	"encoding/json"
	"github.com/BluePecker/JwtAuth/engine/server/parameter/jwt/request"
	"bytes"
	"github.com/BluePecker/JwtAuth/engine/server/parameter"
	"github.com/BluePecker/JwtAuth/engine/server/parameter/jwt/response"
	"github.com/BluePecker/JwtAuth/engine/client/formatter"
	"github.com/BluePecker/JwtAuth/engine/client/formatter/context"
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
			cliAPI := client.NewClient(unixSock)
			if body, err := cliAPI.Post("/v1.0/token/list",
				request.List{Unique: args[0]}); err != nil {
				return err
			} else {
				var res parameter.Response
				if err := json.NewDecoder(body).Decode(&res); err != nil {
					return err
				} else {
					if buffer, err := json.Marshal(res.Data); err != nil {
						return err
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
				}
			}

			return nil
		},
	}

	return cmd
}
