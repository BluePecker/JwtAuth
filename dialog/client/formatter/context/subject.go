package context

import "strings"

type (
	SubjectContext interface {
		fullHeader() string
		addHeader(header string)
	}

	BaseSubjectContext struct {
		header []string
	}
)

func (c *BaseSubjectContext) FullHeader() string {
	if c.header == nil {
		return ""
	}
	return strings.Join(c.header, "\t")
}

func (c *BaseSubjectContext) AddHeader(header string) {
	if c.header == nil {
		c.header = []string{}
	}
	c.header = append(c.header, strings.ToUpper(header))
}
