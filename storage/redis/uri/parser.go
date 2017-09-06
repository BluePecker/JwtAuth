package uri

import (
    "net/url"
    "github.com/go-redis/redis"
    "fmt"
    "reflect"
    "strings"
    "strconv"
    "time"
)

func Parser(uri string) (*redis.Options, *redis.ClusterOptions, error) {
    parsed, err := url.Parse(uri)
    if err != nil {
        return nil, nil, fmt.Errorf("wrong uri format: %s", uri)
    }
    if parsed.Scheme != "redis" {
        return nil, nil, fmt.Errorf("valid schema: %s", parsed.Scheme)
    }
    
    addr := strings.Split(parsed.Host, ",")
    switch len(addr) {
    case 0:
        return nil, nil, fmt.Errorf("%s", "valid host")
    //case 1:
    //    db, _ := strconv.Atoi(strings.Trim(parsed.Path, "/"))
    //    options := &redis.Options{
    //        Addr: addr[0],
    //        Password: pwd,
    //        DB: db,
    //    }
    //    value := reflect.ValueOf(options)
    //    for k, v := range parsed.Query() {
    //        field := reflect.Indirect(value).FieldByName(k)
    //        if field.CanSet() {
    //            switch field.Type().String() {
    //            case "int":
    //                i, _ := strconv.Atoi(v[0])
    //                field.SetInt(int64(i))
    //                break;
    //            case "string":
    //                field.SetString(v[0])
    //                break;
    //            case "bool":
    //                b, _ := strconv.ParseBool(v[0])
    //                field.SetBool(b)
    //                break;
    //            case "time.Duration":
    //                i, _ := strconv.Atoi(v[0])
    //                t := time.Duration(i * 1000)
    //                field.Set(reflect.ValueOf(t))
    //            }
    //        }
    //    }
    //    return options, nil, nil
    default:
        options := &redis.ClusterOptions{
            Addrs: addr,
        }
        if parsed.User != nil {
            options.Password, _ = parsed.User.Password()
        }
        
        value := reflect.ValueOf(options)
        for k, v := range parsed.Query() {
            field := reflect.Indirect(value).FieldByName(k)
            if field.CanSet() {
                switch field.Type().String() {
                case "int":
                    i, _ := strconv.Atoi(v[0])
                    field.SetInt(int64(i))
                    break;
                case "string":
                    field.SetString(v[0])
                    break;
                case "bool":
                    b, _ := strconv.ParseBool(v[0])
                    field.SetBool(b)
                    break;
                case "time.Duration":
                    i, _ := strconv.Atoi(v[0])
                    t := time.Duration(i * 1000)
                    field.Set(reflect.ValueOf(t))
                }
            }
        }
        return nil, options, nil
    }
}