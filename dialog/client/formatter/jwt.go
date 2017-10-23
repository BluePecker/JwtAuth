package formatter

import (
	"github.com/BluePecker/JwtAuth/dialog/client/formatter/context"
	"github.com/BluePecker/JwtAuth/dialog/server/parameter/jwt/response"
	"bytes"
)

const (
	AddrHeader   = "ADDR"
	TLLHeader    = "TLL"
	DeviceHeader = "DEVICE"
	SingedHeader = "SINGED"

	QuietFormat    = "{{.Singed}}"
	JwtTableFormat = "table {{.Addr}}\t{{.TTL}}\t{{.Device}}\t{{.Singed}}"
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
			ctx.Template = `Singed: {{.Singed}}`
		} else {
			//ctx.Template = `Addr: {{.Addr}}\nTTL: {{.TTL}}\nDevice: {{.Device}}\nSinged: {{.Singed}}\n`
			ctx.Template = `Addr: {{.Addr}}\nDevice: {{.Device}}\nSinged: {{.Singed}}\n`
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

//func (j *JsonWebToken) TLL() string {
//	j.AddHeader(TLLHeader)
//	return strconv.FormatFloat(j.jwt.TTL, 'f', -1, 64)
//}

func (j *JsonWebToken) Device() string {
	j.AddHeader(DeviceHeader)
	return j.jwt.Device
}

func (j *JsonWebToken) Singed() string {
	j.AddHeader(SingedHeader)
	return j.jwt.Singed
}
