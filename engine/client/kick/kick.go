package kick

import (
	"github.com/spf13/cobra"
	"fmt"
	"github.com/BluePecker/JwtAuth/engine/client"
	"github.com/BluePecker/JwtAuth/engine/server/parameter/jwt/request"
	"github.com/BluePecker/JwtAuth/engine/server/parameter"
	"encoding/json"
	"errors"
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
			if body, err := cli.Post("/v1.0/token/kick", request.Kick{Unique: args[0]}); err != nil {
				return err
			} else {
				defer body.Close()
				var res parameter.Response
				if err := json.NewDecoder(body).Decode(&res); err != nil {
					return err
				} else {
					if res.Code != 200 {
						return errors.New(res.Message)
					}
					fmt.Printf("successfully kicked out the user.\n")
					return nil
				}
			}
		},
	}
	return cmd
}
