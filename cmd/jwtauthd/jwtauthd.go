package jwtauthd

import (
    "github.com/spf13/cobra"
    "github.com/spf13/pflag"
    "github.com/spf13/viper"
)

type Args struct {
    Port int
    Conf string
}

type JwtAuthCommand struct {
    Args  Args
    Cmd   *cobra.Command
    Viper *viper.Viper
}

var JwtAuth *JwtAuthCommand = &JwtAuthCommand{}

func UsageTemplate() string {
    return `Usage:{{if .Runnable}}
  {{if .HasAvailableFlags}}{{appendIfNotPresent .UseLine "[OPTIONS] COMMAND [arg...]"}}{{else}}{{.UseLine}}{{end}}{{end}}{{if .HasAvailableSubCommands}}
  {{ .CommandPath}} [command]{{end}}{{if gt .Aliases 0}}
Aliases:
  {{.NameAndAliases}}
{{end}}{{if .HasExample}}
Examples:
{{ .Example }}{{end}}{{ if .HasAvailableSubCommands}}
Commands:{{range .Commands}}{{if .IsAvailableCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasAvailableLocalFlags}}
Options:
{{.LocalFlags.FlagUsages | trimRightSpace}}{{end}}{{ if .HasAvailableInheritedFlags}}
Global Flags:
{{.InheritedFlags.FlagUsages | trimRightSpace}}{{end}}{{if .HasHelpSubCommands}}
Additional help topics:{{range .Commands}}{{if .IsHelpCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasAvailableSubCommands }}
Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
}

func init() {
    JwtAuth.Viper = viper.GetViper()
    JwtAuth.Cmd = &cobra.Command{
        Use: "jwt-authd",
        Short: "jwt auth server",
        Long: "a simple jwt auth server",
        SilenceErrors: true,
        RunE: func(cmd *cobra.Command, args []string) error {
            // todo
            return nil
        },
        PersistentPreRun: func(cmd *cobra.Command, args []string) {
            // todo
        },
    }
    
    var PFlags *pflag.FlagSet = JwtAuth.Cmd.Flags()
    
    PFlags.IntVarP(&JwtAuth.Args.Port, "port", "p", 6010, "Set the server listening port")
    PFlags.StringVarP(&JwtAuth.Args.Conf, "config", "c", "/etc/jwt_authd.json", "Set the config file path")
    
    // todo
    
    JwtAuth.Cmd.SetUsageTemplate(UsageTemplate())
}