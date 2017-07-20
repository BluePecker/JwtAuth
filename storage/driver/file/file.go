package file

import "sync"

type File struct {
    mu sync.RWMutex
}
