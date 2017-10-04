package action

import (
    "github.com/spf13/cobra"
    "net/http"
    "net"
    "context"
    "fmt"
    "io/ioutil"
    "log"
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
    
        resp, err := client.Get("http://unix/v1/version")
        if err != nil {
            log.Fatal(err)
        }
        version, err := ioutil.ReadAll(resp.Body)
        resp.Body.Close()
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("%s", version)
        return nil
    },
}