package context

type (
	SubjectContext interface {
		fullHeader() string
		addHeader(header string)
	}
)
