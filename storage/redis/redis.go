package redis

import (
    "github.com/BluePecker/JwtAuth/storage/driver"
    "sync"
    "time"
    "github.com/BluePecker/JwtAuth/storage/header"
)

type Redis struct {
    mu       sync.RWMutex
    createAt time.Time
    values   header.MemStore
}

func init() {
    driver.Register("redis", &Redis{})
}