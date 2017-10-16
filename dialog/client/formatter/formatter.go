package formatter

import (
	"text/template"
	"io"
	"bytes"
	"strings"
	"fmt"
	"github.com/BluePecker/JwtAuth/pkg/templates"
	"text/tabwriter"
)

const (
	TableKey = "table"
)

type (
	Context struct {
		Output io.Writer
		Format string

		header string
		buffer *bytes.Buffer
		format string
		table  bool
	}
)

func (c *Context) pretreatment() {
	c.format = c.Format

	if strings.HasPrefix(c.Format, TableKey) {
		c.table = true
		c.format = c.format[len(TableKey):]
	}

	c.format = strings.Trim(c.format, " ")
	c.format = strings.NewReplacer(`\t`, "\t", `\n`, "\n").Replace(c.format)
}

func (c *Context) template() (*template.Template, error) {
	tpl, err := templates.Parse(c.format)
	if err != nil {
		c.buffer.WriteString(fmt.Sprintf("Template parsing error: %v\n", err))
		c.buffer.WriteTo(c.Output)
	}
	return tpl, err
}

func (c *Context) postFormat(tpl *template.Template, custom customContext) {
	if c.table {
		if len(c.header) == 0 {
			// if we still don't have a header, we didn't have any containers so we need to fake it to get the right headers from the template
			tpl.Execute(bytes.NewBufferString(""), custom)
			c.header = custom.fullHeader()
		}

		t := tabwriter.NewWriter(c.Output, 20, 1, 3, ' ', 0)
		t.Write(append([]byte(c.header + "\n")))
		c.buffer.WriteTo(t)
		t.Flush()
	} else {
		c.buffer.WriteTo(c.Output)
	}
}

func (c *Context) contextFormat(tpl *template.Template, custom customContext) error {
	if err := tpl.Execute(c.buffer, custom); err != nil {
		c.buffer = bytes.NewBufferString(fmt.Sprintf("Template parsing error: %v\n", err))
		c.buffer.WriteTo(c.Output)
		return err
	}
	if c.table && len(c.header) == 0 {
		c.header = custom.fullHeader()
	}
	c.buffer.WriteString("\n")
	return nil
}
