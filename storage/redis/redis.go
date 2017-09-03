package redis

import (
    "strconv"
    "sync"
    "time"
    "github.com/BluePecker/JwtAuth/storage"
    "github.com/go-redis/redis"
    "fmt"
)

type Redis struct {
    create time.Time
    mu     sync.RWMutex
    client *redis.Client
}

func (R *Redis) Initializer(opt storage.Option) error {
    R.client = redis.NewClient(&redis.Options{
        Network: "tcp",
        Addr: opt.Host + ":" + strconv.Itoa(opt.Port),
        PoolSize: opt.PoolSize,
    })
    err := R.client.Ping().Err()
    if err != nil {
        defer R.client.Close()
    }
    return err
}

func (R *Redis) TTL(key string) float64 {
    R.mu.RLock()
    defer R.mu.RUnlock()
    return R.client.TTL(key).Val().Seconds()
}

func (R *Redis) Read(key string) (interface{}, error) {
    R.mu.RLock()
    defer R.mu.RUnlock()
    status := R.client.Get(key)
    return status.Val(), status.Err()
}

func (R *Redis) ReadInt(key string) (int, error) {
    R.mu.RLock()
    defer R.mu.RUnlock()
    status := R.client.Get(key)
    if status.Err() != nil {
        return 0, status.Err()
    }
    return strconv.Atoi(status.Val())
}

func (R *Redis) ReadString(key string) (string, error) {
    R.mu.RLock()
    defer R.mu.RUnlock()
    status := R.client.Get(key)
    if status.Err() != nil {
        return "", status.Err()
    }
    return status.Val(), nil
}

func (R *Redis) Upgrade(key string, expire int) {
    R.mu.Lock()
    defer R.mu.Unlock()
    if v, err := R.Read(key); err != nil {
        R.Write(key, v, expire)
    }
}

func (R *Redis) Write(key string, value interface{}, expire int) {
    R.mu.Lock()
    defer R.mu.Unlock()
    R.save(key, value, expire, false)
}

func (R *Redis) WriteImmutable(key string, value interface{}, expire int) {
    R.mu.Lock()
    defer R.mu.Unlock()
    R.save(key, value, expire, true)
}

func (R *Redis) Remove(key string) {
    R.mu.Lock()
    defer R.mu.Unlock()
    R.remove(key)
}

func (R *Redis) remove(key string) error {
    status := R.client.Del(key)
    return status.Err()
}

func (R *Redis) save(key string, value interface{}, expire int, immutable bool) error {
    if immutable {
        if cmd := R.client.HGet(key, "i"); strconv.ParseBool(cmd.Val()) {
            return fmt.Errorf("this key(%s) write protection", key)
        }
    }
    R.client.Pipelined(func(pipe redis.Pipeliner) error {
        pipe.HSet(key, "v", value);
        pipe.HSet(key, "i", immutable);
        pipe.Expire(key, time.Duration(expire) * time.Second);
        return nil
    })
    return nil
}

func init() {
    storage.Register("redis", &Redis{})
}