package storage

import (
	"fmt"
)

type (
	Engine interface {
		Initializer(opts string) error

		HSet(key, field string, value interface{}, maxLen, expire int64) error

		HScan(key string, do func(token string, ttl float64)) error

		HGet(key, field string) (string, float64, error)

		HRem(key string, field ... string) error
	}
)

var provider = make(map[string]Engine)

func Register(name string, driver Engine) {
	if driver == nil {
		panic("storage: register driver is nil")
	}
	if _, find := provider[name]; find {
		panic("storage: register called twice for " + name)
	}

	provider[name] = driver
}

func New(name string, opts string) (*Engine, error) {
	if storage, find := provider[name]; !find {
		return nil, fmt.Errorf("storage: unknown driver %q (forgotten import?)", name)
	} else {
		if err := storage.Initializer(opts); err != nil {
			return nil, fmt.Errorf("storage: %q driver init failed", name);
		}
		return &storage, nil;
	}
}
