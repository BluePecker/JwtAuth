package version

type Backend interface {
    // 版本号
    Version() (string, error)
}