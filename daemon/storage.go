package daemon

import "github.com/BluePecker/JwtAuth/pkg/storage"

func (d *Daemon) Storage() (err error) {
    d.StorageE, err = storage.New(d.Options.Storage.Driver, d.Options.Storage.Opts)
    return err
}
