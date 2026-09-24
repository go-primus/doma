package http

import (
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
)

// var client *resty.Client

// func init() {
// 	client = resty.New()
// }

type Client struct {
	client *resty.Client
	host   string
	port   string
}

func NewClient(host, port string) *Client {
	client := resty.New()
	client.SetLogger(&logWrap{})
	client.SetAllowGetMethodPayload(true)
	return &Client{
		client: client,
		host:   host,
		port:   port,
	}
}

func (c *Client) SetDebug(d bool) *Client {
	c.client.SetDebug(d)
	return c
}

func (c *Client) getDomain() string {
	return fmt.Sprintf("http://%s:%s", c.host, c.port)
}

func (c *Client) GET(path string, req interface{}, response interface{}, opts ...option) (err error) {

	request := c.client.R()
	c.handleOpt(request, opts...)
	if response != nil {
		request.SetResult(response)
	}

	resp, err := request.
		SetBody(req).
		Get(fmt.Sprintf("%s%s", c.getDomain(), path))
	if err != nil {
		return errors.New("err" + err.Error())
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("err:%v:%s", resp.StatusCode(), string(resp.Body()))
	}
	return
}

func (c *Client) POST(path string, req interface{}, response interface{}, opts ...option) (err error) {

	request := c.client.R()
	c.handleOpt(request, opts...)

	if response != nil {
		request.SetResult(response)
	}

	resp, err := request.
		SetBody(req).
		Post(fmt.Sprintf("%s%s", c.getDomain(), path))
	if err != nil {
		return errors.New("err" + err.Error())
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("err:%v:%s", resp.StatusCode(), string(resp.Body()))
	}
	return
}

func (c *Client) DELETE(path string, req interface{}, response interface{}, opts ...option) (err error) {

	request := c.client.R()
	c.handleOpt(request, opts...)
	if response != nil {
		request = request.SetResult(response)
	}

	resp, err := request.SetBody(req).Delete(fmt.Sprintf("%s%s", c.getDomain(), path))
	if err != nil {
		return errors.New("err" + err.Error())
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("err:%v:%s", resp.StatusCode(), string(resp.Body()))
	}
	return
}

func (c *Client) PUT(path string, req interface{}, response interface{}, opts ...option) (err error) {

	request := c.client.R()
	c.handleOpt(request, opts...)

	if response != nil {
		request = request.SetResult(response)
	}
	resp, err := request.SetBody(req).Put(fmt.Sprintf("%s%s", c.getDomain(), path))
	if err != nil {
		return errors.New("err" + err.Error())
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("err:%v:%s", resp.StatusCode(), string(resp.Body()))
	}
	return
}
