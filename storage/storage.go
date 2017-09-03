package storage

import (
    "fmt"
)

type (
    Option struct {
        Path     string
        Host     string
        Port     int
        Username string
        Password string
        PoolSize int
    }
    
    Driver interface {
        Read(key string) (v interface{}, err error)
        
        ReadInt(key string) (v int, err error)
        
        ReadString(key string) (string, error)
        
        Upgrade(key string, expire int)
        
        Initializer(options Option) error
        
        Write(key string, value interface{}, expire int)
        
        TTL(key string) int
        
        WriteImmutable(key string, value interface{}, expire int)
        
        Remove(key string)
    }
)

var provider = make(map[string]Driver)

func Register(name string, driver Driver) {
    if driver == nil {
        panic("storage: register driver is nil")
    }
    if _, find := provider[name]; find {
        panic("storage: register called twice for " + name)
    }
    
    provider[name] = driver
}

func New(name string, options Option) (Driver, error) {
    if storage, find := provider[name]; !find {
        return nil, fmt.Errorf("storage: unknown driver %q (forgotten import?)", name)
    } else {
        if err := storage.Initializer(options); err != nil {
            return nil, fmt.Errorf("storage: %q driver init failed", name);
        }
        return storage, nil;
    }
}