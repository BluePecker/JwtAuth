package storage

import (
    "fmt"
    "reflect"
    "time"
    "strconv"
    "errors"
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
    
    Entry struct {
        value     interface{}
        ttl       int64
        immutable bool
        version   int64
    }
    
    MemStore map[string]Entry
)

var provider = make(map[string]Driver)

func (e *Entry) Value() interface{} {
    if !e.immutable {
        return e.value
    }
    vv := reflect.Indirect(reflect.ValueOf(e.value))
    switch vv.Type().Kind() {
    case reflect.Map:
        newMap := reflect.MakeMap(vv.Type())
        for _, k := range vv.MapKeys() {
            newMap.SetMapIndex(k, vv.MapIndex(k))
        }
        return newMap
    case reflect.Slice:
        newSlice := reflect.MakeSlice(vv.Type(), vv.Len(), vv.Cap())
        reflect.Copy(newSlice, vv)
        return newSlice
    default:
        return vv.Interface()
    }
}

func (ms *MemStore) Exist(key string) bool {
    if _, ok := (*ms)[key]; !ok {
        return false;
    }
    return true;
}

func (ms *MemStore) Len() int {
    return len(*ms)
}

func (ms *MemStore) Reset() {
    *ms = make(map[string]Entry)
}

func (ms *MemStore) Remove(key string) bool {
    args := *ms
    if _, find := args[key]; !find {
        return false
    } else {
        delete(args, key)
        return true
    }
}

func (ms *MemStore) Visit(visitor func(key string, value interface{})) {
    for key, value := range (*ms) {
        visitor(key, value)
    }
}

func (ms *MemStore) clear(key string, expire int, timestamp int64) {
    timer := time.Duration(expire)
    time.AfterFunc(time.Second * timer, func() {
        if _, ok := (*ms)[key]; ok {
            if (*ms)[key].version == timestamp {
                ms.Remove(key)
            }
        }
    })
}

func (ms *MemStore) save(key string, value interface{}, expire int, immutable bool) error {
    if expire >= 0 {
        if len(*ms) == 0 {
            *ms = make(map[string]Entry)
        }
        tm := time.Now().UnixNano()
        if entry, find := (*ms)[key]; find {
            if entry.immutable {
                return fmt.Errorf("this key(%s) write protection", key)
            }
        }
        (*ms)[key] = Entry{
            value: value,
            version: tm,
            ttl: tm + int64(expire) * 1e9,
            immutable: immutable,
        }
        if expire > 0 {
            ms.clear(key, expire, tm)
        }
    }
    return nil
}

func (ms *MemStore) Set(key string, value interface{}, expire int) error {
    return ms.save(key, value, expire, false)
}

func (ms *MemStore) SetImmutable(key string, value interface{}, expire int) error {
    return ms.save(key, value, expire, true)
}

func (ms *MemStore) Get(key string) (interface{}, error) {
    args := *ms
    if entry, find := args[key]; find {
        return entry.Value(), nil
    } else {
        return nil, fmt.Errorf("can not find value for %s", key)
    }
}

func (ms *MemStore) GetString(key string) (string, error) {
    if v, err := ms.Get(key); err != nil {
        return "", err
    } else {
        if value, ok := v.(string); !ok {
            return "", fmt.Errorf("can not convert %#v to string", v)
        } else {
            return value, nil
        }
    }
}

func (ms *MemStore) GetInt(key string) (int, error) {
    v, _ := ms.Get(key)
    if vInt, ok := v.(int); ok {
        return vInt, nil
    }
    if vString, ok := v.(string); ok {
        return strconv.Atoi(vString)
    }
    return -1, errors.New(fmt.Sprintf("unable to find or parse the integer, found: %#v", v))
}

func (ms *MemStore) TTL(key string) int {
    if _, ok := (*ms)[key]; !ok {
        return -1
    }
    return int(((*ms)[key].ttl - time.Now().UnixNano()) / 1e9)
}

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