package client

// Mb to pkg

import (
	"github.com/valyala/fasthttp"
)

type TClient struct {
	client *fasthttp.Client
}

func NewClient() *TClient {
	ccx := new(TClient)
	ccx.client = &fasthttp.Client{
		NoDefaultUserAgentHeader: true,
	}
	return ccx
}

func (ccx *TClient) SendGetRequest(uri string) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(uri)
	req.Header.SetMethod(fasthttp.MethodGet)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	defer fasthttp.ReleaseRequest(req)

	err := ccx.client.Do(req, resp)
	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}
