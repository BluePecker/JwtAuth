package redis

import (
    "sync"
    "time"
    "github.com/BluePecker/JwtAuth/storage"
)

type Redis struct {
    mu       sync.RWMutex
    createAt time.Time
    values   storage.MemStore
}

func init() {
    storage.Register("redis", &Redis{})
}