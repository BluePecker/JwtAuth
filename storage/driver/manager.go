package driver

import "github.com/BluePecker/JwtAuth/storage/header"

var Manager = make(map[string]header.Driver)

func Register(name string, driver header.Driver) {
    if driver == nil {
        panic("storage: register driver is nil")
    }
    if _, find := Manager[name]; find {
        panic("storage: register called twice for " + name)
    }
    
    Manager[name] = driver
}