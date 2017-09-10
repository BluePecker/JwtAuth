package store

import (
    //_ "github.com/BluePecker/JwtAuth/storage/ram"
    _ "github.com/BluePecker/JwtAuth/pkg/storage/redis"
    "github.com/BluePecker/JwtAuth/pkg/storage"
)

type (
    Driver struct {
        Name string
        Opts string
    }
)

func (d *Driver) New() (*storage.Driver, error) {
    return storage.New(d.Name, d.Opts)
}