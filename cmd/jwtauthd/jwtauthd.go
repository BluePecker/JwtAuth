package jwtauthd

import (
    "github.com/spf13/cobra"
    "github.com/spf13/pflag"
    "github.com/spf13/viper"
    "github.com/BluePecker/JwtAuth/daemon"
    "os"
)

type Args struct {
    PidFile string
    LogFile string
    Port    int
    Host    string
    Driver  string
    Conf    string
    Daemon  bool
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
        Use: "jwt-auth",
        Short: "Jwt auth server",
        Long: "User login information verification service",
        SilenceErrors: true,
        RunE: func(cmd *cobra.Command, args []string) error {
            if _, err := os.Stat(JwtAuth.Args.Conf); err == nil {
                JwtAuth.Viper.SetConfigFile(JwtAuth.Args.Conf)
                if err := JwtAuth.Viper.ReadInConfig(); err != nil {
                    return err
                }
            }
            
            JwtAuth.Args.Port = JwtAuth.Viper.GetInt("port")
            JwtAuth.Args.Host = JwtAuth.Viper.GetString("host")
            JwtAuth.Args.PidFile = JwtAuth.Viper.GetString("pidfile")
            JwtAuth.Args.LogFile = JwtAuth.Viper.GetString("logfile")
            JwtAuth.Args.Daemon = JwtAuth.Viper.GetBool("daemon")
            JwtAuth.Args.Driver = JwtAuth.Viper.GetString("driver")
            
            // 开启SERVER服务
            daemon.NewStart(daemon.Conf{
                PidFile: JwtAuth.Args.PidFile,
                LogFile: JwtAuth.Args.LogFile,
                Port: JwtAuth.Args.Port,
                Host: JwtAuth.Args.Host,
                Daemon: JwtAuth.Args.Daemon,
            })
            
            return nil
        },
    }
    JwtAuth.Cmd.SetUsageTemplate(UsageTemplate())
    
    var PFlags *pflag.FlagSet = JwtAuth.Cmd.Flags()
    
    PFlags.IntVarP(&JwtAuth.Args.Port, "port", "p", 6010, "set the server listening port")
    PFlags.StringVarP(&JwtAuth.Args.Host, "host", "", "127.0.0.1", "set the server bind host")
    PFlags.StringVarP(&JwtAuth.Args.Driver, "driver", "", "redis", "specify the storage driver")
    PFlags.StringVarP(&JwtAuth.Args.Conf, "config", "c", "/etc/jwt_authd.json", "configuration file specifying")
    PFlags.BoolVarP(&JwtAuth.Args.Daemon, "daemon", "d", false, "enable daemon mode")
    PFlags.StringVarP(&JwtAuth.Args.PidFile, "pid", "", "/var/run/jwt-auth.pid", "path to use for daemon PID file")
    PFlags.StringVarP(&JwtAuth.Args.LogFile, "log", "", "/var/log/jwt-auth.log", "path to use for log file")
    
    JwtAuth.Viper.BindPFlag("port", PFlags.Lookup("port"));
    JwtAuth.Viper.BindPFlag("host", PFlags.Lookup("host"));
    JwtAuth.Viper.BindPFlag("pid", PFlags.Lookup("pid"));
    JwtAuth.Viper.BindPFlag("log", PFlags.Lookup("log"));
    JwtAuth.Viper.BindPFlag("daemon", PFlags.Lookup("daemon"));
    JwtAuth.Viper.BindPFlag("driver", PFlags.Lookup("driver"));
}