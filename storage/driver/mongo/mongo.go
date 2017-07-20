package mongo

import "github.com/BluePecker/JwtAuth/storage/driver"

type Mongo struct {
    
}

func init() {
    driver.Register("redis", &Mongo{})
}