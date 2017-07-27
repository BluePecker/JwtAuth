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
    mu     sync.RWMutex
    create time.Time
    values storage.MemStore
    client *redis.Client
}

func (driver *Redis) Initializer(opt storage.Options) error {
    driver.client = redis.NewClient(&redis.Options{
        Network: "tcp",
        Addr: opt.Host + ":" + strconv.Itoa(opt.Port),
        PoolSize: opt.PoolSize,
        OnConnect: func(conn *redis.Conn) error {
            fmt.Println("conn: ", conn)
            return nil
        },
    })
    err := driver.client.Ping().Err()
    if err != nil {
        defer driver.client.Close()
    }
    fmt.Println(err)
    return err
}

func (driver *Redis) TTL(key string) int {
    driver.mu.RLock()
    defer driver.mu.RUnlock()
    return driver.values.TTL(key)
}

func (driver *Redis) Read(key string) (interface{}, error) {
    return nil, nil
}

func (driver *Redis) ReadInt(key string) (int, error) {
    return 0, nil
}

func (driver *Redis) ReadString(key string) string {
    return ""
}

func (driver *Redis) Upgrade(key string, expire int) {
    
}

func (driver *Redis) Write(key string, value interface{}, expire int) {
    
}

func (driver *Redis) WriteImmutable(key string, value interface{}, expire int) {
    
}

func init() {
    storage.Register("redis", &Redis{})
}