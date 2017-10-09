package coder

type Backend interface {
	Decoder()

	Encoder()
}
