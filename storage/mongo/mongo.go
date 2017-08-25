package mongo

import (
    "github.com/BluePecker/JwtAuth/storage"
)

type Mongo struct {
    
}

func init() {
    storage.Register("mongo", &Mongo{})
}