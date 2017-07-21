package header

import (
    "time"
    "reflect"
    "strconv"
    "errors"
    "fmt"
)

type (
    Entry struct {
        Key       string
        ttl       int
        immutable bool
        value     interface{}
    }
    
    Store []Entry
)

func (e Entry) Value() interface{} {
    if e.immutable {
        vv := reflect.Indirect(reflect.ValueOf(e.value))
        switch vv.Type().Kind() {
        case reflect.Map:
            newMap := reflect.MakeMap(vv.Type())
            for _, k := range vv.MapKeys() {
                newMap.SetMapIndex(k, vv.MapIndex(k))
            }
        case reflect.Slice:
            newSlice := reflect.MakeSlice(vv.Type(), vv.Len(), vv.Cap())
            reflect.Copy(newSlice, vv)
            return newSlice
        default:
            return vv.Interface()
        }
    }
    return e.value
}

func (s *Store) Len() int {
    return len(*s)
}

func (s *Store) Reset() {
    *s = (*s)[0:0]
}

func (s *Store) save(key string, value interface{}, expired int, immutable bool) {
    args := *s
    num := len(args)
    ttl := 0
    if expired > 0 {
        ttl = expired + int(time.Now().Unix())
    }
    
    for i := 0; i < num; i++ {
        kv := &args[i]
        if kv.Key == key {
            if immutable && kv.immutable {
                kv.value = value
                kv.ttl = ttl
                kv.immutable = immutable
            } else if kv.immutable == false {
                kv.value = value
                kv.ttl = ttl
                kv.immutable = immutable
            }
            return
        }
    }
    
    c := cap(args)
    if c > num {
        args = args[:num + 1]
        kv := &args[num]
        kv.Key = key
        kv.ttl = ttl
        kv.value = value
        kv.immutable = immutable
        *s = args
        return
    }
    
    kv := Entry{
        Key: key,
        value: value,
        ttl: ttl,
        immutable: immutable,
    }
    
    *s = append(args, kv)
}

func (s *Store) Set(key string, value interface{}, expired int) {
    s.save(key, value, expired, false)
}

func (s *Store) SetImmutable(key string, value interface{}, expired int) {
    s.save(key, value, expired, true)
}

func (s *Store) Get(key string) interface{} {
    args := *s
    num := len(args)
    for i := 0; i < num; i++ {
        if (&args[i]).Key == key {
            return (&args[i]).Value()
        }
    }
    
    return nil
}

func (s *Store) GetString(key string) string {
    if v, ok := s.Get(key).(string); ok {
        return v
    }
    return ""
}

func (s *Store) GetInt(key string) (int, error) {
    v := s.Get(key)
    if vInt, ok := v.(int); ok {
        return vInt, nil
    } else if vString, ok := v.(string); ok {
        return strconv.Atoi(vString)
    }
    return -1, errors.New(fmt.Sprintf("unable to find or parse the integer, found: %#v", v))
}

func (s *Store) Visit(visitor func(key string, value interface{})) {
    args := *s
    for i, num := 0, len(args); i < num; i++ {
        visitor(args[i].Key, args[i].Value())
    }
}

func (s *Store) Remove(key string) bool {
    args := *s
    for i, num := 0, len(args); i < num; i++ {
        if (&args[i]).Key == key {
            args = append(args[:i], args[i + 1:]...)
        }
        *s = args
        return true
    }
    return false
}