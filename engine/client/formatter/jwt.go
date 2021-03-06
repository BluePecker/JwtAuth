package formatter

import (
	"github.com/BluePecker/JwtAuth/engine/client/formatter/context"
	"github.com/BluePecker/JwtAuth/engine/server/parameter/jwt/response"
	"bytes"
	"strconv"
)

const (
	AddrHeader   = "CLIENT ADDR"
	TLLHeader    = "TOKEN TLL"
	DeviceHeader = "DEVICE"
	TokenHeader  = "TOKEN"

	QuietFormat    = "{{.Token}}"
	JwtTableFormat = "table {{.Addr}}\t{{.Tll}}\t{{.Device}}\t{{.Token}}"
)

type (
	JsonWebToken struct {
		context.BaseSubjectContext
		truncate bool
		jwt      response.JsonWebToken
	}

	JsonWebTokenContext struct {
		context.Context
		JsonWebTokens []response.JsonWebToken
	}
)

func (ctx JsonWebTokenContext) Write() {
	switch ctx.Template {
	case context.RawKey:
		if ctx.Quiet {
			ctx.Template = `Token: {{.Token}}`
		} else {
			ctx.Template = `Client Addr: {{.Addr}}\nToken TTL: {{.Tll}}\nDevice: {{.Device}}\nToken: {{.Token}}\n`
		}
	case context.TableKey:
		if ctx.Quiet {
			ctx.Template = QuietFormat
		} else {
			ctx.Template = JwtTableFormat
		}
	}

	ctx.Buffer = bytes.NewBufferString("")
	ctx.PreFormat()

	tpl, err := ctx.Parser()
	if err != nil {
		return
	}

	for _, jwt := range ctx.JsonWebTokens {
		jwtCtx := &JsonWebToken{
			truncate: ctx.Truncate,
			jwt:      jwt,
		}
		err = ctx.ContextFormat(tpl, jwtCtx)
		if err != nil {
			return
		}
	}

	ctx.FormFormat(tpl, &JsonWebToken{})
}

func (j *JsonWebToken) Addr() string {
	j.AddHeader(AddrHeader)
	return j.jwt.Addr
}

func (j *JsonWebToken) Tll() string {
	j.AddHeader(TLLHeader)
	return strconv.FormatFloat(j.jwt.TTL, 'f', -1, 64)
}

func (j *JsonWebToken) Device() string {
	j.AddHeader(DeviceHeader)
	return j.jwt.Device
}

func (j *JsonWebToken) Token() string {
	j.AddHeader(TokenHeader)
	return j.jwt.Singed
}
