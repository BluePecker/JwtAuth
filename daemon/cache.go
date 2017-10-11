package daemon

import (
	_ "github.com/BluePecker/JwtAuth/pkg/storage/redis"
	"github.com/BluePecker/JwtAuth/pkg/storage"
)

func (d *Daemon) NewCache() (err error) {
	d.Cache, err = storage.New(d.Options.Cache.Driver, d.Options.Cache.Opts)
	return err
}
