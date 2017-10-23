package daemon

import (
	"os"
	"fmt"
	"github.com/sevlyar/go-daemon"
	"github.com/Sirupsen/logrus"
	"github.com/BluePecker/JwtAuth/pkg/storage"
	"syscall"
)

type Cache struct {
	Driver string
	Opts   string
}

type TLS struct {
	Key  string
	Cert string
}

type Options struct {
	PidFile  string
	LogFile  string
	LogLevel string
	Port     int
	Host     string
	SockFile string
	Daemon   bool
	Version  bool
	TLS      TLS
	Cache    Cache
	Secret   string
}

type Daemon struct {
	Options *Options

	Cache *storage.Engine
}

func setLogger(level string) {
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	Level, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.Error(err)
		os.Exit(0)
	}
	logrus.SetLevel(Level)
}

func NewDaemon(background bool, args Options) (*Daemon, *daemon.Context) {
	if background {
		ctx := &daemon.Context{
			PidFileName: args.PidFile,
			PidFilePerm: 0644,
			LogFilePerm: 0640,
			Umask:       027,
			WorkDir:     "/",
			LogFileName: args.LogFile,
		}
		if process, err := ctx.Reborn(); err == nil {
			if process != nil {
				return nil, nil
			} else {
				return &Daemon{Options: &args}, ctx
			}
		} else {
			if err == daemon.ErrWouldBlock {
				logrus.Error("daemon already exists.")
			} else {
				logrus.Errorf("Unable to run: ", err)
			}
			os.Exit(0)
		}
	}
	return &Daemon{Options: &args}, nil
}

func NewStart(args Options) {

	setLogger(args.LogLevel)

	if args.Secret == "" {
		fmt.Println("please specify secret for jwt encode.")
		os.Exit(0)
	}

	if progress, ctx := NewDaemon(args.Daemon, args); progress == nil {
		return
	} else {
		if ctx != nil {
			defer ctx.Release()
		}
		if err := progress.NewCache(); err != nil {
			logrus.Error(err)
			os.Exit(0)
		}

		sigterm := make(chan struct{})
		go func() {
			logrus.Infof("ready to listen: http://%s:%d", args.Host, args.Port)
			if err := progress.WebServer(sigterm); err != nil {
				logrus.Error(err)
			} else {
				logrus.Infof("ready to listen: http://%s:%d", args.Host, args.Port)
			}
		}()
		daemon.SetSigHandler(func(sig os.Signal) error {
			close(sigterm)
			return daemon.ErrStop
		}, syscall.SIGTERM)

		if err := daemon.ServeSignals(); err != nil {
			logrus.Error(err)
		}
		logrus.Info("daemon terminated")
	}
}
