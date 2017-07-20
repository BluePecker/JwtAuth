package main

import (
    "reflect"
    "fmt"
    "github.com/BluePecker/JwtAuth/storage"
    "github.com/BluePecker/JwtAuth/storage/header"
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
}
  
