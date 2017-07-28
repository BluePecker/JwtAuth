package redis

import (
    "sync"
    "time"
    "strconv"
    "github.com/BluePecker/JwtAuth/storage"
    "github.com/go-redis/redis"
)

type Redis struct {
    create time.Time
    mu     sync.RWMutex
    mem    storage.MemStore
    client *redis.Client
}

func (er *Redis) Initializer(opt storage.Options) error {
    er.client = redis.NewClient(&redis.Options{
        Network: "tcp",
        Addr: opt.Host + ":" + strconv.Itoa(opt.Port),
        PoolSize: opt.PoolSize,
    })
    err := er.client.Ping().Err()
    if err != nil {
        defer er.client.Close()
    }
    return err
}

func (er *Redis) TTL(key string) int {
    er.mu.RLock()
    defer er.mu.RUnlock()
    if !er.mem.Exist(key) {
        return int(er.client.TTL(key).Val().Seconds())
    }
    return er.mem.TTL(key)
}

func (er *Redis) Read(key string) (interface{}, error) {
    return nil, nil
}

func (er *Redis) ReadInt(key string) (int, error) {
    return 0, nil
}

func (er *Redis) ReadString(key string) string {
    return ""
}

func (er *Redis) Upgrade(key string, expire int) {
    
}

func (er *Redis) Write(key string, value interface{}, expire int) {
    
}

func (er *Redis) WriteImmutable(key string, value interface{}, expire int) {
    
}

func init() {
    storage.Register("redis", &Redis{})
}