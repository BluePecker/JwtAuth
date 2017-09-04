package jwtauthd

import (
    "github.com/spf13/cobra"
    "github.com/spf13/pflag"
    "github.com/spf13/viper"
    "github.com/BluePecker/JwtAuth/daemon"
    "os"
)

type Storage struct {
    Driver     string
    Path       string
    Host       string
    Port       int
    MaxRetries int
    Username   string
    Password   string
    Database   string
}

type Security struct {
    TLS  bool
    Key  string
    Cert string
}

type Args struct {
    PidFile  string
    LogFile  string
    LogLevel string
    Version  bool
    Port     int
    Host     string
    Conf     string
    Daemon   bool
    
    Security Security
    Storage  Storage
}

type JwtAuthCommand struct {
    Args  Args
    Cmd   *cobra.Command
    Viper *viper.Viper
}

var JwtAuth *JwtAuthCommand = &JwtAuthCommand{}

func UsageTemplate() string {
    return `Usage:{{if .Runnable}}
  {{if .HasAvailableFlags}}
  {{appendIfNotPresent .UseLine "[OPTIONS] COMMAND [arg...]"}}{{else}}{{.UseLine}}{{end}}{{end}}{{if .HasAvailableSubCommands}}
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
            JwtAuth.Args.PidFile = JwtAuth.Viper.GetString("pid")
            JwtAuth.Args.LogLevel = JwtAuth.Viper.GetString("log-level")
            JwtAuth.Args.LogFile = JwtAuth.Viper.GetString("log")
            JwtAuth.Args.Version = JwtAuth.Viper.GetBool("version")
            JwtAuth.Args.Daemon = JwtAuth.Viper.GetBool("daemon")
            
            JwtAuth.Args.Storage.Driver = JwtAuth.Viper.GetString("storage.driver")
            JwtAuth.Args.Storage.Path = JwtAuth.Viper.GetString("storage.path")
            JwtAuth.Args.Storage.Host = JwtAuth.Viper.GetString("storage.host")
            JwtAuth.Args.Storage.Port = JwtAuth.Viper.GetInt("storage.port")
            JwtAuth.Args.Storage.MaxRetries = JwtAuth.Viper.GetInt("storage.max-retries")
            JwtAuth.Args.Storage.Username = JwtAuth.Viper.GetString("storage.username")
            JwtAuth.Args.Storage.Password = JwtAuth.Viper.GetString("storage.password")
            JwtAuth.Args.Storage.Database = JwtAuth.Viper.GetString("storage.database")
            JwtAuth.Args.Security.TLS = JwtAuth.Viper.GetBool("security.tls")
            JwtAuth.Args.Security.Cert = JwtAuth.Viper.GetString("security.cert")
            JwtAuth.Args.Security.Key = JwtAuth.Viper.GetString("security.key")
            
            // 开启SERVER服务
            daemon.NewStart(daemon.Options{
                PidFile: JwtAuth.Args.PidFile,
                LogLevel: JwtAuth.Args.LogLevel,
                LogFile: JwtAuth.Args.LogFile,
                Port: JwtAuth.Args.Port,
                Host: JwtAuth.Args.Host,
                Security: daemon.Security{
                    TLS: JwtAuth.Args.Security.TLS,
                    Key: JwtAuth.Args.Security.Key,
                    Cert: JwtAuth.Args.Security.Cert,
                },
                Version: JwtAuth.Args.Version,
                Daemon: JwtAuth.Args.Daemon,
                Storage: daemon.Storage{
                    Driver: JwtAuth.Args.Storage.Driver,
                    Path: JwtAuth.Args.Storage.Path,
                    Host: JwtAuth.Args.Storage.Host,
                    Port: JwtAuth.Args.Storage.Port,
                    MaxRetries: JwtAuth.Args.Storage.MaxRetries,
                    Username: JwtAuth.Args.Storage.Username,
                    Password: JwtAuth.Args.Storage.Password,
                    Database: JwtAuth.Args.Storage.Database,
                },
            })
            
            return nil
        },
    }
    JwtAuth.Cmd.SetUsageTemplate(UsageTemplate())
    
    var PFlags *pflag.FlagSet = JwtAuth.Cmd.Flags()
    
    PFlags.IntVarP(&JwtAuth.Args.Port, "port", "p", 6010, "set the server listening port")
    PFlags.StringVarP(&JwtAuth.Args.Host, "host", "", "127.0.0.1", "set the server bind host")
    PFlags.StringVarP(&JwtAuth.Args.Conf, "config", "c", "/etc/jwt_authd.json", "set configuration file")
    PFlags.BoolVarP(&JwtAuth.Args.Version, "version", "v", false, "print version information and quit")
    PFlags.BoolVarP(&JwtAuth.Args.Daemon, "daemon", "d", false, "enable daemon mode")
    PFlags.StringVarP(&JwtAuth.Args.PidFile, "pid", "", "/var/run/jwt-auth.pid", "path to use for daemon PID file")
    PFlags.StringVarP(&JwtAuth.Args.LogLevel, "log-level", "l", "info", "set the logging level")
    PFlags.StringVarP(&JwtAuth.Args.LogFile, "log", "", "/var/log/jwt-auth.log", "path to use for log file")
    PFlags.StringVarP(&JwtAuth.Args.Storage.Driver, "storage-driver", "", "redis", "specify the storage driver")
    PFlags.StringVarP(&JwtAuth.Args.Storage.Path, "storage-path", "", "", "specify the storage path")
    PFlags.StringVarP(&JwtAuth.Args.Storage.Host, "storage-host", "", "127.0.0.1", "specify the storage host")
    PFlags.IntVarP(&JwtAuth.Args.Storage.Port, "storage-port", "", 6379, "specify the storage port")
    PFlags.IntVarP(&JwtAuth.Args.Storage.MaxRetries, "storage-max-retries", "", 3, "specify the storage max retries")
    PFlags.StringVarP(&JwtAuth.Args.Storage.Username, "storage-username", "", "", "specify the storage username")
    PFlags.StringVarP(&JwtAuth.Args.Storage.Password, "storage-password", "", "", "specify the storage password")
    PFlags.StringVarP(&JwtAuth.Args.Storage.Database, "storage-database", "", "", "specify the storage database")
    PFlags.BoolVarP(&JwtAuth.Args.Security.TLS, "security-tls", "", false, "use TLS and verify the remote")
    PFlags.StringVarP(&JwtAuth.Args.Security.Cert, "security-cert", "", "", "path to TLS certificate file")
    PFlags.StringVarP(&JwtAuth.Args.Security.Key, "security-key", "", "", "path to TLS key file")
    
    JwtAuth.Viper.BindPFlag("port", PFlags.Lookup("port"));
    JwtAuth.Viper.BindPFlag("host", PFlags.Lookup("host"));
    JwtAuth.Viper.BindPFlag("version", PFlags.Lookup("version"));
    JwtAuth.Viper.BindPFlag("daemon", PFlags.Lookup("daemon"));
    JwtAuth.Viper.BindPFlag("pid", PFlags.Lookup("pid"));
    JwtAuth.Viper.BindPFlag("log", PFlags.Lookup("log"));
    JwtAuth.Viper.BindPFlag("log-level", PFlags.Lookup("log-level"));
    JwtAuth.Viper.BindPFlag("storage.driver", PFlags.Lookup("storage-driver"));
    JwtAuth.Viper.BindPFlag("storage.path", PFlags.Lookup("storage-path"));
    JwtAuth.Viper.BindPFlag("storage.host", PFlags.Lookup("storage-host"));
    JwtAuth.Viper.BindPFlag("storage.port", PFlags.Lookup("storage-port"));
    JwtAuth.Viper.BindPFlag("storage.max-retries", PFlags.Lookup("storage-max-retries"));
    JwtAuth.Viper.BindPFlag("storage.username", PFlags.Lookup("storage-username"));
    JwtAuth.Viper.BindPFlag("storage.password", PFlags.Lookup("storage-password"));
    JwtAuth.Viper.BindPFlag("storage.database", PFlags.Lookup("storage-database"));
    JwtAuth.Viper.BindPFlag("security.tls", PFlags.Lookup("security-tls"));
    JwtAuth.Viper.BindPFlag("security.cert", PFlags.Lookup("security-cert"));
    JwtAuth.Viper.BindPFlag("security.key", PFlags.Lookup("security-key"));
}