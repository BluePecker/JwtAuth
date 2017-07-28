package redis

import (
    "sync"
    "time"
    "strconv"
    "github.com/BluePecker/JwtAuth/storage"
    "github.com/go-redis/redis"
    "fmt"
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
    er.mu.RLock()
    defer er.mu.RUnlock()
    if v, err := er.mem.Get(key); err != nil {
        status := er.client.Get(key)
        return status.Val(), status.Err()
    } else {
        return v, nil
    }
}

func (er *Redis) ReadInt(key string) (int, error) {
    er.mu.RLock()
    defer er.mu.RUnlock()
    v, err := er.mem.GetInt(key)
    if err != nil {
        status := er.client.Get(key)
        if status.Err() != nil {
            return 0, status.Err()
        }
        return strconv.Atoi(status.Val())
    }
    return v, nil
}

func (er *Redis) ReadString(key string) (string, error) {
    er.mu.RLock()
    defer er.mu.RUnlock()
    v, err := er.mem.GetString(key)
    if err != nil {
        status := er.client.Get(key)
        if status.Err() != nil {
            return "", status.Err()
        }
        return status.Val(), nil
    }
    return v, nil
}

func (er *Redis) Upgrade(key string, expire int) {
    er.mu.Lock()
    defer er.mu.Unlock()
    if v, err := er.Read(key); err != nil {
        er.Write(key, v, expire)
    }
    
}

func (er *Redis) Write(key string, value interface{}, expire int) {
    er.mu.Lock()
    defer er.mu.Unlock()
    if er.mem.Set(key, value, expire) == nil {
        err := er.flushToDB(key, value, expire)
        fmt.Println(key, err)
    }
}

func (er *Redis) WriteImmutable(key string, value interface{}, expire int) {
    er.mu.Lock()
    defer er.mu.Unlock()
    if er.mem.SetImmutable(key, value, expire) == nil {
        er.flushToDB(key, value, expire)
    }
}

func (er *Redis) flushToDB(key string, value interface{}, expire int) error {
    cmdStatus := er.client.Set(key, value, time.Duration(expire))
    return cmdStatus.Err()
}

func init() {
    storage.Register("redis", &Redis{})
}