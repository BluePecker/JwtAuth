package redis

import (
    "strconv"
    "sync"
    "time"
    "github.com/BluePecker/JwtAuth/storage"
    "github.com/go-redis/redis"
    "crypto/md5"
    "encoding/hex"
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
        Addr: fmt.Sprintf("%s:%d", opt.Host, opt.Port),
        PoolSize: opt.PoolSize,
        DB: opt.Database,
        MaxRetries: opt.MaxRetries,
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
    return R.client.TTL(R.md5Key(key)).Val().Seconds()
}

func (R *Redis) Read(key string) (interface{}, error) {
    R.mu.RLock()
    defer R.mu.RUnlock()
    status := R.client.Get(R.md5Key(key))
    return status.Val(), status.Err()
}

func (R *Redis) ReadInt(key string) (int, error) {
    R.mu.RLock()
    defer R.mu.RUnlock()
    status := R.client.Get(R.md5Key(key))
    if status.Err() != nil {
        return 0, status.Err()
    }
    return strconv.Atoi(status.Val())
}

func (R *Redis) ReadString(key string) (string, error) {
    R.mu.RLock()
    defer R.mu.RUnlock()
    status := R.client.Get(R.md5Key(key))
    if status.Err() != nil {
        return "", status.Err()
    }
    return status.Val(), nil
}

func (R *Redis) Upgrade(key string, expire int) {
    R.mu.Lock()
    defer R.mu.Unlock()
    key = R.md5Key(key)
    if v, err := R.Read(key); err != nil {
        R.Write(key, v, expire)
    }
}

func (R *Redis) Write(key string, value interface{}, expire int) {
    R.mu.Lock()
    defer R.mu.Unlock()
    R.save(R.md5Key(key), value, expire, false)
}

func (R *Redis) WriteImmutable(key string, value interface{}, expire int) {
    R.mu.Lock()
    defer R.mu.Unlock()
    R.save(R.md5Key(key), value, expire, true)
}

func (R *Redis) Remove(key string) {
    R.mu.Lock()
    defer R.mu.Unlock()
    R.remove(R.md5Key(key))
}

func (R *Redis) remove(key string) error {
    status := R.client.Del(key)
    return status.Err()
}

func (R *Redis) save(key string, value interface{}, expire int, immutable bool) error {
    key = R.md5Key(key)
    if cmd := R.client.HGet(key, "i"); strconv.ParseBool(cmd.Val()) {
        return fmt.Errorf("this key(%s) write protection", key)
    }
    R.client.Pipelined(func(pipe redis.Pipeliner) error {
        pipe.HSet(key, "v", value);
        pipe.HSet(key, "i", immutable);
        pipe.Expire(key, time.Duration(expire) * time.Second);
        return nil
    })
    return nil
}

func (R *Redis) md5Key(key string) string {
    hash := md5.New()
    hash.Write([]byte(key))
    return hex.EncodeToString(hash.Sum(nil))
}

func init() {
    storage.Register("redis", &Redis{})
}