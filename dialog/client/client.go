package client

import (
	"net/http"
	"context"
	"net"
	"io"
	"io/ioutil"
)

type (
	Cli struct {
		Client *http.Client
	}
)

func (c *Cli) Get(uri string) ([]byte, error) {
	if res, err := c.Client.Get("http://unix" + uri); err != nil {
		return nil, err
	} else {
		if body, err := ioutil.ReadAll(res.Body); err != nil {
			return nil, err
		} else {
			res.Body.Close()
			return body, nil
		}
	}
}

func (c *Cli) Post(uri string, bodyType string, body io.Reader) ([]byte, error) {
	if res, err := c.Client.Post("http://unix"+uri, bodyType, body); err != nil {
		return nil, err
	} else {
		if body, err := ioutil.ReadAll(res.Body); err != nil {
			return nil, err
		} else {
			res.Body.Close()
			return body, nil
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
