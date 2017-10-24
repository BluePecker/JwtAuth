package redis

import (
	"strconv"
	"sync"
	"time"
	"github.com/BluePecker/JwtAuth/pkg/storage"
	"github.com/go-redis/redis"
	"crypto/md5"
	"encoding/hex"
	"github.com/BluePecker/JwtAuth/pkg/storage/redis/uri"
	"reflect"
	"github.com/Sirupsen/logrus"
	"github.com/kataras/iris/core/errors"
)

type (
	Redis struct {
		create time.Time
		mu     sync.RWMutex
		engine Client
	}

	Client interface {
		Ping() *redis.StatusCmd

		Close() error

		Pipelined(fn func(redis.Pipeliner) error) ([]redis.Cmder, error)

		ZScore(key, field string) *redis.FloatCmd

		ZRem(key string, members ... interface{}) *redis.IntCmd

		ZRange(key string, start, stop int64) *redis.StringSliceCmd
	}
)

func inject(from, target reflect.Value) {
	indirect := reflect.Indirect(target.Elem())
	for index := 0; index < from.Elem().NumField(); index++ {
		name := from.Elem().Type().Field(index).Name
		f1 := from.Elem().FieldByName(name)
		f2 := indirect.FieldByName(name)
		if f2.IsValid() {
			if f1.Type() == f2.Type() && f2.CanSet() {
				f2.Set(f1)
			}
		}
	}
}

func (r *Redis) Initializer(opts string) error {
	generic, err := uri.Parser(opts)
	if err != nil {
		return err
	}
	switch reflect.ValueOf(generic).Interface().(type) {
	case *redis.ClusterOptions:
		options := &redis.ClusterOptions{}
		inject(reflect.ValueOf(generic), reflect.ValueOf(options))
		r.engine = redis.NewClusterClient(options)
		break
	case *redis.Options:
		options := &redis.Options{}
		inject(reflect.ValueOf(generic), reflect.ValueOf(options))
		r.engine = redis.NewClient(options)
		break
	}
	statusCmd := r.engine.Ping()
	if statusCmd.Err() != nil {
		logrus.Error(statusCmd.Err())
		defer r.engine.Close()
	}
	return statusCmd.Err()
}

func (r *Redis) HSet(key, field string, value interface{}, maxLen, expire int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.engine.Pipelined(func(p redis.Pipeliner) error {
		exp := time.Duration(expire) * time.Second
		tmp := jwtMd5(key)
		val := jwtMd5(tmp + field)
		score := expire + int64(time.Now().Second())
		if cmd := p.ZAdd(tmp, redis.Z{Score: float64(score), Member: val}); cmd.Err() != nil {
			return cmd.Err()
		} else {
			if cmd := p.Expire(tmp, exp); cmd.Err() != nil {
				return cmd.Err()
			}
			if cmd := p.Set(val, value, exp); cmd.Err() != nil {
				p.ZRem(tmp, redis.Z{Score: float64(score), Member: val})
				return cmd.Err()
			}
		}
		if cmd := p.ZCard(tmp); cmd.Err() != nil {
			return cmd.Err()
		} else {
			if cmd.Val() > maxLen {
				if cmd := p.ZRange(tmp, 0, cmd.Val()-maxLen); cmd.Err() != nil {
					return cmd.Err()
				} else {
					p.Del(cmd.Val()...)
				}
				cmd = p.ZRemRangeByRank(tmp, 0, cmd.Val()-maxLen)
				if cmd.Err() != nil {
					return cmd.Err()
				}
			}
		}
		return nil
	})
	return err
}

func (r *Redis) HGetString(key, field string) (string, float64, error) {
	r.mu.RLock()
	defer r.mu.RLock()
	tmp := jwtMd5(key)
	val := jwtMd5(tmp + field)
	if cmd := r.engine.ZScore(tmp, val); cmd.Err() != nil {
		return "", -1, cmd.Err()
	} else if cmd.Val() < float64(time.Now().Second()) {
		if cmd := r.engine.ZRem(tmp, val); cmd.Err() != nil {
			return "", -1, cmd.Err()
		} else {
			return "", -1, errors.New("key has been expired.")
		}
	} else {
		cmd, err := r.engine.Pipelined(func(p redis.Pipeliner) error {
			if cmd := p.Get(val); cmd.Err() != nil {
				return cmd.Err()
			}
			if cmd := p.TTL(val); cmd.Err() != nil {
				return cmd.Err()
			}
			return nil
		})
		if err == nil {
			if ttl, err := strconv.ParseFloat(cmd[1].String(), 64); err != nil {
				return "", -1, err
			} else {
				return cmd[0].String(), ttl, nil
			}
		}
		return "", -1, err
	}
}

func (r *Redis) HKeys(key string) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RLock()
	tmp := jwtMd5(key)
	if cmd := r.engine.ZRange(tmp, 0, -1); cmd.Err() != nil {
		return []string{}, cmd.Err()
	} else {
		return cmd.Val(), nil
	}
}

func (r *Redis) HRem(key string, field ... string) error {
	r.mu.Lock()
	defer r.mu.Lock()
	var val []string
	tmp := jwtMd5(key)
	for _, v := range field {
		val = append(val, jwtMd5(tmp+v))
	}
	_, err := r.engine.Pipelined(func(p redis.Pipeliner) error {
		p.ZRem(tmp, val...)
		p.Del(val...)
		return nil
	})
	return err
}

func jwtMd5(key string) string {
	hash := md5.New()
	hash.Write([]byte(key))
	return hex.EncodeToString(hash.Sum([]byte("jwt#")))
}

func init() {
	storage.Register("redis", &Redis{})
}
