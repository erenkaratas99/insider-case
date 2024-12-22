package clients

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
)

type BaseClient struct {
	client  *fasthttp.Client
	baseURL string
}

func NewBaseClient(baseURL string) *BaseClient {
	return &BaseClient{
		client:  new(fasthttp.Client),
		baseURL: baseURL,
	}
}

func (b *BaseClient) GET(path string, opts ...map[string]string) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	res := fasthttp.AcquireResponse()

	req.SetRequestURI(b.baseURL + path)
	req.Header.SetMethod("GET")

	for _, opt := range opts {
		for k, v := range opt {
			req.Header.Set(k, v)
		}
	}

	if err := b.client.Do(req, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (b *BaseClient) PUT(path string, body interface{}, opts ...map[string]string) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	res := fasthttp.AcquireResponse()

	req.SetRequestURI(b.baseURL + path)
	req.Header.SetMethod("PUT")

	for _, opt := range opts {
		for k, v := range opt {
			req.Header.Set(k, v)
		}
	}
	contentType := ""
	if len(opts) > 0 {
		contentType = opts[0]["Content-Type"]
	}

	if contentType == "application/json" {
		body, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		req.SetBody(body)
	} else {
		buf, ok := body.([]byte)
		if !ok {
			return nil, fmt.Errorf("invalid body type")
		}
		req.SetBody(buf)
	}

	if err := b.client.Do(req, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (b *BaseClient) POST(path string, body interface{}, opts ...map[string]string) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	res := fasthttp.AcquireResponse()

	req.SetRequestURI(b.baseURL + path)
	req.Header.SetMethod("POST")

	for _, opt := range opts {
		for k, v := range opt {
			req.Header.Set(k, v)
		}
	}

	contentType := ""
	if len(opts) > 0 {
		contentType = opts[0]["Content-Type"]
	}

	if contentType == "application/json" {
		body, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		req.SetBody(body)
	} else {
		buf, ok := body.([]byte)
		if !ok {
			return nil, fmt.Errorf("invalid body type")
		}
		req.SetBody(buf)
	}

	if err := b.client.Do(req, res); err != nil {
		return nil, err
	}

	return res, nil

}
