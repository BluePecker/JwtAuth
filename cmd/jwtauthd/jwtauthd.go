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
    TLS    bool
    Key    string
    Cert   string
    Verify bool
}

type Args struct {
    PidFile  string
    LogFile  string
    LogLevel string
    Version  bool
    Port     int
    Host     string
    Conf     string
    Secret   string
    Daemon   bool
    
    Security Security
    Storage  Storage
}

type RootCommand struct {
    Args  Args
    Cmd   *cobra.Command
    Viper *viper.Viper
}

var RootCmd *RootCommand = &RootCommand{}

func UsageTemplate() string {
    return `Usage:{{if .Runnable}}{{if .HasAvailableFlags}}
  {{appendIfNotPresent .UseLine "[OPTIONS] COMMAND [arg...]"}}{{else}}{{.UseLine}}{{end}}{{end}}{{if .HasAvailableSubCommands}}
  {{ .CommandPath}} [command]
  {{end}}{{if gt .Aliases 0}}
Aliases:{{.NameAndAliases}}
{{end}}{{if .HasExample}}
Examples:{{ .Example }}
{{end}}{{ if .HasAvailableLocalFlags}}
Options:
{{.LocalFlags.FlagUsages | trimRightSpace}}
{{end}}{{ if .HasAvailableSubCommands}}
Commands:{{range .Commands}}{{if .IsAvailableCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}
{{end}}{{ if .HasAvailableInheritedFlags}}
Global Flags:
{{.InheritedFlags.FlagUsages | trimRightSpace}}{{end}}{{if .HasHelpSubCommands}}
Additional help topics:{{range .Commands}}{{if .IsHelpCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasAvailableSubCommands }}
Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
}

func init() {
    RootCmd.Viper = viper.GetViper()
    
    RootCmd.Cmd = &cobra.Command{
        Use: "jwt-auth",
        Short: "Jwt auth server",
        Long: "User login information verification service",
        SilenceErrors: true,
        RunE: func(cmd *cobra.Command, args []string) error {
            if _, err := os.Stat(RootCmd.Args.Conf); err == nil {
                RootCmd.Viper.SetConfigFile(RootCmd.Args.Conf)
                if err := RootCmd.Viper.ReadInConfig(); err != nil {
                    return err
                }
            }
            
            RootCmd.Args.Port = RootCmd.Viper.GetInt("port")
            RootCmd.Args.Host = RootCmd.Viper.GetString("host")
            RootCmd.Args.PidFile = RootCmd.Viper.GetString("pid")
            RootCmd.Args.LogLevel = RootCmd.Viper.GetString("log-level")
            RootCmd.Args.LogFile = RootCmd.Viper.GetString("log")
            RootCmd.Args.Secret = RootCmd.Viper.GetString("secret")
            RootCmd.Args.Version = RootCmd.Viper.GetBool("version")
            RootCmd.Args.Daemon = RootCmd.Viper.GetBool("daemon")
            
            RootCmd.Args.Storage.Driver = RootCmd.Viper.GetString("storage.driver")
            RootCmd.Args.Storage.Path = RootCmd.Viper.GetString("storage.path")
            RootCmd.Args.Storage.Host = RootCmd.Viper.GetString("storage.host")
            RootCmd.Args.Storage.Port = RootCmd.Viper.GetInt("storage.port")
            RootCmd.Args.Storage.MaxRetries = RootCmd.Viper.GetInt("storage.max-retries")
            RootCmd.Args.Storage.Username = RootCmd.Viper.GetString("storage.username")
            RootCmd.Args.Storage.Password = RootCmd.Viper.GetString("storage.password")
            RootCmd.Args.Storage.Database = RootCmd.Viper.GetString("storage.database")
            RootCmd.Args.Security.TLS = RootCmd.Viper.GetBool("security.tls")
            RootCmd.Args.Security.Cert = RootCmd.Viper.GetString("security.cert")
            RootCmd.Args.Security.Key = RootCmd.Viper.GetString("security.key")
            
            // 开启SERVER服务
            daemon.NewStart(daemon.Options{
                PidFile: RootCmd.Args.PidFile,
                LogLevel: RootCmd.Args.LogLevel,
                LogFile: RootCmd.Args.LogFile,
                Port: RootCmd.Args.Port,
                Host: RootCmd.Args.Host,
                Security: daemon.Security{
                    TLS: RootCmd.Args.Security.TLS,
                    Verify: RootCmd.Args.Security.Verify,
                    Key: RootCmd.Args.Security.Key,
                    Cert: RootCmd.Args.Security.Cert,
                },
                Version: RootCmd.Args.Version,
                Daemon: RootCmd.Args.Daemon,
                Storage: daemon.Storage{
                    Driver: RootCmd.Args.Storage.Driver,
                    Path: RootCmd.Args.Storage.Path,
                    Host: RootCmd.Args.Storage.Host,
                    Port: RootCmd.Args.Storage.Port,
                    MaxRetries: RootCmd.Args.Storage.MaxRetries,
                    Username: RootCmd.Args.Storage.Username,
                    Password: RootCmd.Args.Storage.Password,
                    Database: RootCmd.Args.Storage.Database,
                },
                Secret: RootCmd.Args.Secret,
            })
            
            return nil
        },
    }
    RootCmd.Cmd.SetUsageTemplate(UsageTemplate())
    
    var PFlags *pflag.FlagSet = RootCmd.Cmd.Flags()
    
    PFlags.IntVarP(&RootCmd.Args.Port, "port", "p", 6010, "set the server listening port")
    PFlags.StringVarP(&RootCmd.Args.Host, "host", "", "127.0.0.1", "set the server bind host")
    PFlags.StringVarP(&RootCmd.Args.Conf, "config", "c", "/etc/jwt_authd.json", "set configuration file")
    PFlags.BoolVarP(&RootCmd.Args.Version, "version", "v", false, "print version information and quit")
    PFlags.BoolVarP(&RootCmd.Args.Daemon, "daemon", "d", false, "enable daemon mode")
    PFlags.StringVarP(&RootCmd.Args.Secret, "secret", "s", "", "specify secret for jwt encode")
    PFlags.StringVarP(&RootCmd.Args.PidFile, "pid", "", "/var/run/jwt-auth.pid", "path to use for daemon PID file")
    PFlags.StringVarP(&RootCmd.Args.LogLevel, "log-level", "l", "info", "set the logging level")
    PFlags.StringVarP(&RootCmd.Args.LogFile, "log", "", "/var/log/jwt-auth.log", "path to use for log file")
    PFlags.StringVarP(&RootCmd.Args.Storage.Driver, "storage-driver", "", "redis", "specify the storage driver")
    PFlags.StringVarP(&RootCmd.Args.Storage.Path, "storage-path", "", "", "specify the storage path")
    PFlags.StringVarP(&RootCmd.Args.Storage.Host, "storage-host", "", "127.0.0.1", "specify the storage host")
    PFlags.IntVarP(&RootCmd.Args.Storage.Port, "storage-port", "", 6379, "specify the storage port")
    PFlags.IntVarP(&RootCmd.Args.Storage.MaxRetries, "storage-max-retries", "", 3, "specify the storage max retries")
    PFlags.StringVarP(&RootCmd.Args.Storage.Username, "storage-username", "", "", "specify the storage username")
    PFlags.StringVarP(&RootCmd.Args.Storage.Password, "storage-password", "", "", "specify the storage password")
    PFlags.StringVarP(&RootCmd.Args.Storage.Database, "storage-database", "", "", "specify the storage database")
    PFlags.BoolVarP(&RootCmd.Args.Security.TLS, "tls", "", false, "use TLS; implied by --tlsverify")
    PFlags.StringVarP(&RootCmd.Args.Security.Cert, "tlscert", "", "", "path to TLS certificate file")
    PFlags.StringVarP(&RootCmd.Args.Security.Key, "tlskey", "", "", "path to TLS key file")
    PFlags.BoolVarP(&RootCmd.Args.Security.Verify, "tlsverify", "", false, "path to TLS key file")
    
    RootCmd.Viper.BindPFlag("port", PFlags.Lookup("port"))
    RootCmd.Viper.BindPFlag("host", PFlags.Lookup("host"))
    RootCmd.Viper.BindPFlag("version", PFlags.Lookup("version"))
    RootCmd.Viper.BindPFlag("secret", PFlags.Lookup("secret"))
    RootCmd.Viper.BindPFlag("daemon", PFlags.Lookup("daemon"))
    RootCmd.Viper.BindPFlag("pid", PFlags.Lookup("pid"))
    RootCmd.Viper.BindPFlag("log", PFlags.Lookup("log"))
    RootCmd.Viper.BindPFlag("log-level", PFlags.Lookup("log-level"))
    RootCmd.Viper.BindPFlag("storage.driver", PFlags.Lookup("storage-driver"))
    RootCmd.Viper.BindPFlag("storage.path", PFlags.Lookup("storage-path"))
    RootCmd.Viper.BindPFlag("storage.host", PFlags.Lookup("storage-host"))
    RootCmd.Viper.BindPFlag("storage.port", PFlags.Lookup("storage-port"))
    RootCmd.Viper.BindPFlag("storage.max-retries", PFlags.Lookup("storage-max-retries"))
    RootCmd.Viper.BindPFlag("storage.username", PFlags.Lookup("storage-username"))
    RootCmd.Viper.BindPFlag("storage.password", PFlags.Lookup("storage-password"))
    RootCmd.Viper.BindPFlag("storage.database", PFlags.Lookup("storage-database"))
    RootCmd.Viper.BindPFlag("security.cert", PFlags.Lookup("tlscert"))
    RootCmd.Viper.BindPFlag("security.key", PFlags.Lookup("tlskey"))
    RootCmd.Viper.BindPFlag("security.verify", PFlags.Lookup("tlsverify"))
    
    RootCmd.Cmd.AddCommand(StopCmd, TokenCmd, VersionCmd)
}