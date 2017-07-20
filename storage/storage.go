package storage

import (
    "fmt"
    "github.com/BluePecker/JwtAuth/storage/driver"
    "github.com/BluePecker/JwtAuth/storage/header"
)

func New(name string, options header.Options) (header.Driver, error) {
    if storage, find := driver.Manager[name]; !find {
        return nil, fmt.Errorf("storage: unknown driver %q (forgotten import?)", name)
    } else {
        if err := storage.Initializer(options); err != nil {
            return nil, fmt.Errorf("storage: %q driver init failed", name);
        }
        return storage, nil;
    }
}