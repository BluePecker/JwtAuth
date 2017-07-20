package redis

import "github.com/BluePecker/JwtAuth/storage/driver"

type Redis struct {
    
}

func init() {
    driver.Register("redis", &Redis{})
}