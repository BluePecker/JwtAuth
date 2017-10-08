package action

import (
    "github.com/spf13/cobra"
    "net/http"
    "log"
    "io/ioutil"
    "fmt"
    "context"
    "net"
)

var StopCmd = &cobra.Command{
    Use: "stop",
    Short: "stop running server",
    Long: "stop running server",
    RunE: func(cmd *cobra.Command, args []string) error {
        client := &http.Client{
            Transport: &http.Transport{
                DialContext:func(ctx context.Context, network, addr string) (net.Conn, error) {
                    return net.Dial("unix", RootCmd.Args.SockFile)
                },
            },
        }

        resp, err := client.Get("http://unix/v1/stop")
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
