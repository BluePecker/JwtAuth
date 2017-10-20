package context

import (
	"io"
	"bytes"
	"strings"
	"fmt"
	"github.com/BluePecker/JwtAuth/pkg/templates"
	"text/template"
	"text/tabwriter"
)

const (
	TableKey = "table"
)

type (
	Context struct {
		Truncate bool
		Writer   io.Writer
		Template string
		Quiet    bool

		header string
		table  bool
		format string
		buffer *bytes.Buffer
	}
)

func (c *Context) PreFormat() {
	c.format = c.Template

	if strings.HasPrefix(c.Template, TableKey) {
		c.table = true
		c.format = c.format[len(TableKey):]
	}

	c.format = strings.Trim(c.format, " ")
	r := strings.NewReplacer(`\t`, "\t", `\n`, "\n")
	c.format = r.Replace(c.format)
}

func (c *Context) Parser() (*template.Template, error) {
	tpl, err := templates.Parse(c.format)
	if err != nil {
		c.buffer.WriteString(fmt.Sprintf("Template parsing error: %v\n", err))
		c.buffer.WriteTo(c.Writer)
	}
	return tpl, err
}

func (c *Context) FormFormat(tpl *template.Template, subject SubjectContext) {
	if c.table {
		if len(c.header) == 0 {
			// if we still don't have a header, we didn't have any containers so we need to fake it to get the right headers from the template
			tpl.Execute(bytes.NewBufferString(""), subject)
			c.header = subject.fullHeader()
		}

		t := tabwriter.NewWriter(c.Writer, 20, 1, 3, ' ', 0)
		t.Write([]byte(c.header))
		t.Write([]byte("\n"))
		c.buffer.WriteTo(t)
		t.Flush()
	} else {
		c.buffer.WriteTo(c.Writer)
	}
}

func (c *Context) ContextFormat(tpl *template.Template, subject SubjectContext) error {
	if err := tpl.Execute(c.buffer, subject); err != nil {
		c.buffer = bytes.NewBufferString(fmt.Sprintf("Template parsing error: %v\n", err))
		c.buffer.WriteTo(c.Writer)
		return err
	}
	if c.table && len(c.header) == 0 {
		c.header = subject.fullHeader()
	}
	c.buffer.WriteString("\n")
	return nil
}
