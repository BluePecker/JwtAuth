package signal

type Backend interface {
	// 终止
	Stop() error
}
