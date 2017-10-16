package formatter

import "strings"

type (
	customContext interface {
		fullHeader() string
		addHeader(header string)
	}
	baseContext struct {
		header []string
	}
)

func (c *baseContext) fullHeader() string {
	if c.header == nil {
		return ""
	}
	return strings.Join(c.header, "\t")
}

func (c *baseContext) addHeader(header string) {
	if c.header == nil {
		c.header = []string{}
	}
	c.header = append(c.header, strings.ToUpper(header))
}
