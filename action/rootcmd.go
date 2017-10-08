package action

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	//"github.com/BluePecker/JwtAuth/daemon"
	"os"
	"fmt"
	"strings"
)

type Storage struct {
	Driver string
	Opts   string
}

type TLS struct {
	Key  string
	Cert string
}

type Args struct {
	PidFile  string
	LogFile  string
	LogLevel string
	Version  bool
	SockFile string
	Port     int
	Host     string
	Conf     string
	Secret   string
	Daemon   bool

	TLS     TLS
	Storage Storage
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
		SilenceErrors: true,
		Use:           "jwt",
		Short:         "Jwt auth server",
		Long:          "User login information verification service",
		RunE: func(cmd *cobra.Command, args []string) error {

			fmt.Printf("%s\n", strings.Join(args, " "))

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
			RootCmd.Args.SockFile = RootCmd.Viper.GetString("unix-sock")
			RootCmd.Args.Secret = RootCmd.Viper.GetString("secret")
			RootCmd.Args.Version = RootCmd.Viper.GetBool("version")
			RootCmd.Args.Daemon = RootCmd.Viper.GetBool("daemon")

			RootCmd.Args.Storage.Driver = RootCmd.Viper.GetString("storage.driver")
			RootCmd.Args.Storage.Opts = RootCmd.Viper.GetString("storage.opts")
			RootCmd.Args.TLS.Key = RootCmd.Viper.GetString("tls.key")
			RootCmd.Args.TLS.Cert = RootCmd.Viper.GetString("tls.cert")

			// 开启SERVER服务
			//daemon.NewStart(daemon.Options{
			//	PidFile:  RootCmd.Args.PidFile,
			//	LogLevel: RootCmd.Args.LogLevel,
			//	LogFile:  RootCmd.Args.LogFile,
			//	SockFile: RootCmd.Args.SockFile,
			//	Port:     RootCmd.Args.Port,
			//	Host:     RootCmd.Args.Host,
			//	TLS: daemon.TLS{
			//		Cert: RootCmd.Args.TLS.Cert,
			//		Key:  RootCmd.Args.TLS.Key,
			//	},
			//	Version: RootCmd.Args.Version,
			//	Daemon:  RootCmd.Args.Daemon,
			//	Storage: daemon.Storage{
			//		Driver: RootCmd.Args.Storage.Driver,
			//		Opts:   RootCmd.Args.Storage.Opts,
			//	},
			//	Secret: RootCmd.Args.Secret,
			//})

			return nil
		},
	}
	RootCmd.Cmd.SetUsageTemplate(UsageTemplate())

	var PFlags *pflag.FlagSet = RootCmd.Cmd.Flags()

	PFlags.IntVarP(&RootCmd.Args.Port, "port", "p", 6010, "set the server listening port")
	PFlags.StringVarP(&RootCmd.Args.Host, "host", "", "127.0.0.1", "set the server bind host")
	PFlags.StringVarP(&RootCmd.Args.Conf, "config", "c", "/etc/jwt.json", "set configuration file")
	PFlags.BoolVarP(&RootCmd.Args.Version, "version", "v", false, "print version information and quit")
	PFlags.BoolVarP(&RootCmd.Args.Daemon, "daemon", "d", false, "enable daemon mode")
	PFlags.StringVarP(&RootCmd.Args.Secret, "secret", "s", "", "specify secret for jwt encode")
	PFlags.StringVarP(&RootCmd.Args.PidFile, "pid", "", "/var/run/jwt.pid", "path to use for daemon PID file")
	PFlags.StringVarP(&RootCmd.Args.LogLevel, "log-level", "l", "info", "set the logging level")
	PFlags.StringVarP(&RootCmd.Args.LogFile, "log", "", "/var/log/jwt.log", "path to use for log file")
	PFlags.StringVarP(&RootCmd.Args.SockFile, "unix-sock", "u", "/var/run/jwt.sock", "communication between the client and the daemon")
	PFlags.StringVarP(&RootCmd.Args.Storage.Driver, "storage-driver", "", "redis", "specify the storage driver")
	PFlags.StringVarP(&RootCmd.Args.Storage.Opts, "storage-opts", "", "redis://127.0.0.1:6379/1?PoolSize=20&MaxRetries=3&PoolTimeout=1000", "specify the storage uri")
	PFlags.StringVarP(&RootCmd.Args.TLS.Cert, "tlscert", "", "", "path to TLS certificate file")
	PFlags.StringVarP(&RootCmd.Args.TLS.Key, "tlskey", "", "", "path to TLS key file")

	RootCmd.Viper.BindPFlag("port", PFlags.Lookup("port"))
	RootCmd.Viper.BindPFlag("host", PFlags.Lookup("host"))
	RootCmd.Viper.BindPFlag("version", PFlags.Lookup("version"))
	RootCmd.Viper.BindPFlag("secret", PFlags.Lookup("secret"))
	RootCmd.Viper.BindPFlag("daemon", PFlags.Lookup("daemon"))
	RootCmd.Viper.BindPFlag("pid", PFlags.Lookup("pid"))
	RootCmd.Viper.BindPFlag("log", PFlags.Lookup("log"))
	RootCmd.Viper.BindPFlag("unix-sock", PFlags.Lookup("unix-sock"))
	RootCmd.Viper.BindPFlag("log-level", PFlags.Lookup("log-level"))
	RootCmd.Viper.BindPFlag("storage.driver", PFlags.Lookup("storage-driver"))
	RootCmd.Viper.BindPFlag("storage.opts", PFlags.Lookup("storage-opts"))
	RootCmd.Viper.BindPFlag("tls.cert", PFlags.Lookup("tlscert"))
	RootCmd.Viper.BindPFlag("tls.key", PFlags.Lookup("tlskey"))

	RootCmd.Cmd.AddCommand(StopCmd, TokenCmd, VersionCmd)
}
