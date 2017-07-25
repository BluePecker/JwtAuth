package main

import (
    "reflect"
    "fmt"
    "github.com/BluePecker/JwtAuth/storage"
    "github.com/BluePecker/JwtAuth/storage/header"
    "unsafe"
    "time"
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
    
    storage.New("redis", header.Options{});
    
    
    store := &header.MemStore{}
    
    //store.SetImmutable("name", "shuchao", 3)
    store.Set("name", "hi", 2)
    
    fmt.Println(store)
    time.Sleep(time.Duration(1 * time.Second))
    fmt.Println(store)
    
    store.Set("name", "me", 0)
    fmt.Println(store)
    time.Sleep(time.Duration(5 * time.Second))
    fmt.Println(store)
}
  
