package daemon

import (
	"syscall"
	"os"
)

func (d *Daemon) Stop() error {
	return syscall.Kill(os.Getpid(), syscall.SIGTERM)
}
