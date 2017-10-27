package client

import (
	"net/http"
	"context"
	"net"
	"encoding/json"
	"bytes"
	"github.com/kataras/iris/core/errors"
	"io"
	"fmt"
)

type (
	Cli struct {
		Client *http.Client
	}
)

func (c *Cli) Get(uri string) (io.ReadCloser, error) {
	if res, err := c.Client.Get("http://unix" + uri); err != nil {
		return nil, err
	} else {
		if res.StatusCode != 200 {
			return nil, errors.New(fmt.Sprintf("%s %d", "response code is ", res.StatusCode))
		}
		return res.Body, nil
	}
}

func (c *Cli) Post(uri string, req interface{}) (io.ReadCloser, error) {
	if byteReq, err := json.Marshal(req); err != nil {
		return nil, err
	} else {
		if res, err := c.Client.Post("http://unix"+uri,
			"application/json;charset=utf-8",
			bytes.NewBuffer(byteReq)); err != nil {
			return nil, err
		} else {
			if res.StatusCode != 200 {
				return nil, errors.New(fmt.Sprintf("%s %d", "response code is ", res.StatusCode))
			}
			return res.Body, nil
		}
	}
}

func NewClient(sock string) *Cli {
	return &Cli{
		&http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return net.Dial("unix", sock)
				},
			},
		},
	}
}
