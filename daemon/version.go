package daemon

import "fmt"

const (
	VERSION = "1.0.0"
)

func (d *Daemon) Version() (string, error) {
	return fmt.Sprintf("version %s.", VERSION), nil
}
