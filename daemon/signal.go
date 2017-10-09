package daemon

import (
	"syscall"
	"os"
)

func (d *Daemon) Stop() {
	syscall.Kill(os.Getpid(), syscall.SIGTERM)
}
