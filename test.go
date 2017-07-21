package main

import (
    "reflect"
    "fmt"
    "github.com/BluePecker/JwtAuth/storage"
    "github.com/BluePecker/JwtAuth/storage/header"
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
    
    fmt.Println(reflect.ValueOf(*user).FieldByName("Foot"))
    
    fmt.Println(reflect.New(reflect.ValueOf(*user).Type()))
    
    storage.New("redis", header.Options{});
    
    
    g := &G{}
    fmt.Println(g)
    g.E()
    fmt.Println(g)
    
    
    fmt.Println("g start: ")
    start := time.Now()
    fmt.Println(start)
    s1 := make([]string, 100000000, 100000000)
    m1 := make(map[int]string)
    // 生成数据
    for i := 0; i < 100000000; i++ {
        s1 = append(s1, "shu chao")
        m1[i] = "shu chao"
    }
    fmt.Println("g end: ")
    end := time.Now()
    fmt.Println(end)
    
    time.Sleep(time.Duration(2 * time.Second))
    
    fmt.Println("s start: ")
    fmt.Println(time.Now())
    for k := range s1 {
        if k < 0 {}
    }
    fmt.Println("s end: ")
    fmt.Println(time.Now())
    
    fmt.Println("m start: ")
    fmt.Println(time.Now())
    for k := range m1 {
        if k < 0 {}
    }
    fmt.Println("m end: ")
    fmt.Println(time.Now())
    
}
  
