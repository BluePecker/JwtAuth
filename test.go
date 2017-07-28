package main

import (
    "reflect"
    "fmt"
    "unsafe"
    "github.com/BluePecker/JwtAuth/storage"
    _ "github.com/BluePecker/JwtAuth/storage/redis"
)

type Hand struct {
    
}

type Foot struct {
    K string
}

func (f *Foot) Hi() {
    fmt.Println("fh")
}

type User struct {
    Name string
    //Hand *Hand
    Foot *Foot
}

func (u *User) Ec() string {
    fmt.Println("ec")
    return "ec"
}

type G struct {
    Name string
}

func (g *G) E() {
    *g = G{
        Name: "shuc",
    }
}

func main() {
    user := &User{
        Name: "SC",
    }
    
    fmt.Printf("测试: %+v", (*struct {
        Ne  string
        Age int
        H   int
    })(unsafe.Pointer(user)))
    
    fmt.Println(reflect.ValueOf(*user).FieldByName("Foot"))
    
    fmt.Println(reflect.New(reflect.ValueOf(*user).Type()))
    
    redis, err := storage.NewManager("redis", storage.Options{
        Host: "127.0.0.1",
        Port: 6379,
        PoolSize: 20,
    });
    
    if err != nil {
        fmt.Println(err)
    }
    
    redis.Write("jwt", "13658009009", 0)
    redis.Write("auth", "23658009009", 30)
    
    //fmt.Println("redis ttl: ", redis.TTL("jwt"))
    //fmt.Println("redis ttl: ", redis.TTL("auth"))
    //v, err := redis.ReadString("jwt")
    //fmt.Println("redis jwt: ", v, err)
    //v, err = redis.ReadString("auth")
    //fmt.Println("redis auth: ", v, err)
    //
    //time.Sleep(time.Duration(15) * time.Second)
    //v, err = redis.ReadString("auth")
    //fmt.Println("redis auth: ", v, err)
    
    //store := &storage.MemStore{}
    //
    ////store.SetImmutable("name", "shuchao", 3)
    //store.Set("name", "hi", 2)
    //
    //fmt.Println(store)
    //time.Sleep(time.Duration(1 * time.Second))
    //fmt.Println(store)
    //
    //store.Set("name", "me", 0)
    //fmt.Println(store)
    //time.Sleep(time.Duration(5000 * time.Second))
    //fmt.Println(store)
}
  
