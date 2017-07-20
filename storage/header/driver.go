package header

type Driver interface {
    Read() (n int, err error)
    
    Write() (n int, err error)
    
    Upgrade(name string) error
    
    Flush() error
    
    Initializer(options Options) error
}
