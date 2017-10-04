package action

import (
    "github.com/spf13/cobra"
    "net/http"
    "net"
    "context"
    "fmt"
)

var VersionCmd = &cobra.Command{
    Use: "version",
    Short: "show the Jwt version information",
    Long: "show the Jwt version information",
    RunE: func(cmd *cobra.Command, args []string) error {
        client := &http.Client{
            Transport: &http.Transport{
                DialContext:func(ctx context.Context, network, addr string) (net.Conn, error) {
                    return net.Dial("unix", RootCmd.Args.SockFile)
                },
            },
        }
        
        reps, _ := client.Get("http://unix/v1/version")
        fmt.Println(reps.Body)
        return nil
    },
}