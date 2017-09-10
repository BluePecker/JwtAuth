package daemon

import (
    "github.com/BluePecker/JwtAuth/daemon/store"
)

func (d *Daemon) Storage() (err error) {
    Opts := d.Options.Storage
    d.Store, err = (&store.Driver{
        Name:Opts.Driver,
        Opts:Opts.Opts,
    }).New()
    return err
}