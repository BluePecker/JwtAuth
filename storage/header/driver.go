package header

import (
    "time"
    "reflect"
    "strconv"
    "fmt"
    "errors"
)

type (
    Driver interface {
        Read(key string) (v interface{}, err error)
        
        ReadInt(key string) (v int, err error)
        
        ReadString(key string) string
        
        Upgrade(string string, expired int)
        
        Initializer(options Options) error
        
        Write(key string, value interface{}, expired int)
        
        WriteImmutable(key string, value interface{}, expired int)
    }
    
    Entry struct {
        value     interface{}
        ttl       int64
        immutable bool
        version   int64
    }
    
    MemStore map[string]Entry
)

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

func (ms *MemStore) save(key string, value interface{}, expired int, immutable bool) {
    tm := time.Now().UnixNano()
    if entry, find := (*ms)[key]; find {
        if !entry.immutable {
            (*ms)[key] = Entry{
                value: value,
                version: tm,
                ttl: tm + int64(expired) * 1e9,
                immutable: immutable,
            }
            if expired <= 0 {
                return
            }
            timer := time.Duration(expired) * time.Second
            time.AfterFunc(timer, func() {
                if _, ok := (*ms)[key]; ok {
                    if (*ms)[key].version == tm {
                        ms.Remove(key)
                    }
                }
            })
        }
        return
    }
    
    (*ms)[key] = Entry{
        value: value,
        version: tm,
        ttl: tm + int64(expired) * 1e9,
        immutable: immutable,
    }
    if expired <= 0 {
        return
    }
    timer := time.Duration(expired) * time.Second
    time.AfterFunc(timer, func() {
        if _, ok := (*ms)[key]; ok {
            if (*ms)[key].version == tm {
                ms.Remove(key)
            }
        }
    })
}

func (ms *MemStore) Set(key string, value interface{}, expired int) {
    ms.save(key, value, expired, false)
}

func (ms *MemStore) SetImmutable(key string, value interface{}, expired int) {
    ms.save(key, value, expired, true)
}

func (ms *MemStore) Get(key string) interface{} {
    return (*ms)[key].Value()
}

func (ms *MemStore) GetString(key string) string {
    if value, ok := ms.Get(key).(string); !ok {
        return ""
    } else {
        return value
    }
}

func (ms *MemStore) GetInt(key string) (int, error) {
    v := ms.Get(key)
    if vInt, ok := v.(int); ok {
        return vInt, nil
    }
    if vString, ok := v.(string); ok {
        return strconv.Atoi(vString)
    }
    return -1, errors.New(fmt.Sprintf("unable to find or parse the integer, found: %#v", v))
}